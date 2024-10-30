package services

import (
	"fmt"
	"gorm.io/datatypes"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"
)

const InternalTokenAudience = "solar-network"

// DetectRisk is used for detect user environment is suitable for no multifactorial authenticating or not.
// Return the remaining steps, value is from 1 to 2, may appear 3 if user enabled the third-authentication-factor.
func DetectRisk(user models.Account, ip, ua string) int {
	var clue int64
	if err := database.C.
		Where(models.AuthTicket{AccountID: user.ID, IpAddress: ip}).
		Where("available_at IS NOT NULL").
		Model(models.AuthTicket{}).
		Count(&clue).Error; err == nil {
		if clue >= 1 {
			return 1
		}
	}

	return 3
}

// PickTicketAttempt is trying to pick up the ticket that hasn't completed but created by a same client (identify by ip address).
// Then the client can continue their journey to get ticket activated.
func PickTicketAttempt(user models.Account, ip string) (models.AuthTicket, error) {
	var ticket models.AuthTicket
	if err := database.C.
		Where("account_id = ? AND ip_address = ? AND expired_at < ? AND available_at IS NULL", user.ID, ip, time.Now()).
		First(&ticket).Error; err != nil {
		return ticket, err
	}
	return ticket, nil
}

func NewTicket(user models.Account, ip, ua string) (models.AuthTicket, error) {
	var ticket models.AuthTicket
	if ticket, err := PickTicketAttempt(user, ip); err == nil {
		return ticket, nil
	}

	steps := DetectRisk(user, ip, ua)
	if count := CountUserFactor(user.ID); count <= 0 {
		return ticket, fmt.Errorf("specified user didn't enable sign in")
	} else {
		steps = min(steps, int(count))

		cfg, err := GetAuthPreference(user.ID)
		if err == nil && cfg.Config.Data().MaximumAuthSteps >= 1 {
			steps = min(steps, cfg.Config.Data().MaximumAuthSteps)
		}
	}

	ticket = models.AuthTicket{
		Claims:      []string{"*"},
		Audiences:   []string{InternalTokenAudience},
		IpAddress:   ip,
		UserAgent:   ua,
		StepRemain:  steps,
		ExpiredAt:   nil,
		AvailableAt: nil,
		AccountID:   user.ID,
	}

	err := database.C.Save(&ticket).Error

	return ticket, err
}

func NewOauthTicket(
	user models.Account,
	client models.ThirdClient,
	claims, audiences []string,
	ip, ua string, nonce *string,
) (models.AuthTicket, error) {
	if nonce != nil && len(*nonce) == 0 {
		nonce = nil
	}

	ticket := models.AuthTicket{
		Claims:       claims,
		Audiences:    audiences,
		IpAddress:    ip,
		UserAgent:    ua,
		GrantToken:   lo.ToPtr(uuid.NewString()),
		AccessToken:  lo.ToPtr(uuid.NewString()),
		RefreshToken: lo.ToPtr(uuid.NewString()),
		AvailableAt:  lo.ToPtr(time.Now()),
		ExpiredAt:    lo.ToPtr(time.Now().Add(7 * 24 * time.Hour)),
		Nonce:        nonce,
		ClientID:     &client.ID,
		AccountID:    user.ID,
	}

	if err := database.C.Save(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func ActiveTicket(ticket models.AuthTicket) (models.AuthTicket, error) {
	if ticket.AvailableAt != nil {
		return ticket, nil
	} else if err := ticket.IsCanBeAvailble(); err != nil {
		return ticket, err
	}

	ticket.AvailableAt = lo.ToPtr(time.Now())
	ticket.GrantToken = lo.ToPtr(uuid.NewString())
	ticket.AccessToken = lo.ToPtr(uuid.NewString())
	ticket.RefreshToken = lo.ToPtr(uuid.NewString())

	if err := database.C.Save(&ticket).Error; err != nil {
		return ticket, err
	} else {
		_ = NewNotification(models.Notification{
			Topic:    "passport.security.alert",
			Title:    "New sign in alert",
			Subtitle: fmt.Sprintf("New sign in from %s", ticket.IpAddress),
			Body:     fmt.Sprintf("Your account just got a new sign in from %s. Make sure you recongize this device, or sign out it immediately and reset password.", ticket.IpAddress),
			Metadata: datatypes.JSONMap{
				"ip_address":   ticket.IpAddress,
				"created_at":   ticket.CreatedAt,
				"available_at": ticket.AvailableAt,
			},
			AccountID: ticket.AccountID,
			Priority:  5,
		})
	}

	return ticket, nil
}

func ActiveTicketWithPassword(ticket models.AuthTicket, password string) (models.AuthTicket, error) {
	if ticket.AvailableAt != nil {
		return ticket, nil
	} else if ticket.StepRemain == 1 {
		return ticket, fmt.Errorf("multi-factor authentication required")
	}

	factor, err := GetPasswordTypeFactor(ticket.AccountID)
	if err != nil {
		return ticket, fmt.Errorf("unable to authenticate, password factor was not found: %v", err)
	} else if err := CheckFactor(factor, password); err != nil {
		return ticket, fmt.Errorf("invalid password: %v", err)
	}

	ticket.StepRemain--
	ticket.FactorTrail = append(ticket.FactorTrail, int(factor.ID))

	ticket.AvailableAt = lo.ToPtr(time.Now())
	ticket.GrantToken = lo.ToPtr(uuid.NewString())
	ticket.AccessToken = lo.ToPtr(uuid.NewString())
	ticket.RefreshToken = lo.ToPtr(uuid.NewString())

	if err := database.C.Save(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func PerformTicketCheck(ticket models.AuthTicket, factor models.AuthFactor, code string) (models.AuthTicket, error) {
	if ticket.AvailableAt != nil {
		return ticket, nil
	} else if ticket.StepRemain <= 0 {
		return ticket, nil
	}

	if lo.Contains(ticket.FactorTrail, int(factor.ID)) {
		return ticket, fmt.Errorf("already checked this ticket with factor %d", factor.ID)
	}

	if err := CheckFactor(factor, code); err != nil {
		return ticket, fmt.Errorf("invalid code: %v", err)
	}

	ticket.StepRemain--
	ticket.FactorTrail = append(ticket.FactorTrail, int(factor.ID))

	if ticket.IsCanBeAvailble() == nil {
		return ActiveTicket(ticket)
	} else {
		if err := database.C.Save(&ticket).Error; err != nil {
			return ticket, err
		}
	}

	return ticket, nil
}

func RotateTicket(ticket models.AuthTicket, fullyRestart ...bool) (models.AuthTicket, error) {
	ticket.GrantToken = lo.ToPtr(uuid.NewString())
	ticket.AccessToken = lo.ToPtr(uuid.NewString())
	ticket.RefreshToken = lo.ToPtr(uuid.NewString())
	if len(fullyRestart) > 0 && fullyRestart[0] {
		ticket.LastGrantAt = nil
	}
	err := database.C.Save(&ticket).Error
	return ticket, err
}

func DoAutoSignoff() {
	duration := 7 * 24 * time.Hour
	deadline := time.Now().Add(-duration)

	log.Debug().Time("before", deadline).Msg("Now signing off tickets...")

	if tx := database.C.
		Where("last_grant_at < ?", deadline).
		Delete(&models.AuthTicket{}); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred when running auto sign off...")
	} else {
		log.Debug().Int64("affected", tx.RowsAffected).Msg("Auto sign off accomplished.")
	}
}
