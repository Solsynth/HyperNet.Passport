package services

import (
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/localize"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func GetPasswordTypeFactor(userId uint) (models.AuthFactor, error) {
	var factor models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		Type:      models.PasswordAuthFactor,
		AccountID: userId,
	}).First(&factor).Error

	return factor, err
}

func GetFactor(id uint) (models.AuthFactor, error) {
	var factor models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		BaseModel: models.BaseModel{ID: id},
	}).First(&factor).Error

	return factor, err
}

func ListUserFactor(userId uint) ([]models.AuthFactor, error) {
	var factors []models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		AccountID: userId,
	}).Find(&factors).Error

	return factors, err
}

func CountUserFactor(userId uint) int64 {
	var count int64
	database.C.Where(models.AuthFactor{
		AccountID: userId,
	}).Model(&models.AuthFactor{}).Count(&count)

	return count
}

func GetFactorCode(factor models.AuthFactor, ip string) (bool, error) {
	switch factor.Type {
	case models.InAppNotifyFactor:
		var user models.Account
		if err := database.C.Where(&models.Account{
			BaseModel: models.BaseModel{ID: factor.AccountID},
		}).First(&user).Error; err != nil {
			return true, err
		}

		secret := uuid.NewString()[:6]

		identifier := fmt.Sprintf("%s%d", gap.FactorOtpPrefix, factor.ID)
		_, err := gap.Jt.Publish(identifier, []byte(secret))
		if err != nil {
			return true, fmt.Errorf("error during publish message: %v", err)
		} else {
			log.Info().Uint("factor", factor.ID).Str("secret", secret).Msg("Published one-time-password to JetStream...")
		}

		err = PushNotification(models.Notification{
			Topic:     "passport.security.otp",
			Title:     localize.L.GetLocalizedString("subjectLoginOneTimePassword", user.Language),
			Body:      fmt.Sprintf(localize.L.GetLocalizedString("shortBodyLoginOneTimePassword", user.Language), secret),
			Account:   user,
			AccountID: user.ID,
			Metadata:  map[string]any{"secret": secret},
		}, true)
		if err != nil {
			log.Warn().Err(err).Uint("factor", factor.ID).Msg("Failed to delivery one-time-password via notify...")
			return true, nil
		}
		return true, nil
	case models.EmailPasswordFactor:
		var user models.Account
		if err := database.C.Where(&models.Account{
			BaseModel: models.BaseModel{ID: factor.AccountID},
		}).Preload("Contacts").First(&user).Error; err != nil {
			return true, err
		}

		secret := uuid.NewString()[:6]

		identifier := fmt.Sprintf("%s%d", gap.FactorOtpPrefix, factor.ID)
		_, err := gap.Jt.Publish(identifier, []byte(secret))
		if err != nil {
			return true, fmt.Errorf("error during publish message: %v", err)
		} else {
			log.Info().Uint("factor", factor.ID).Str("secret", secret).Msg("Published one-time-password to JetStream...")
		}

		subject := fmt.Sprintf("[%s] %s", viper.GetString("name"), localize.L.GetLocalizedString("subjectLoginOneTimePassword", user.Language))

		content := localize.L.RenderLocalizedTemplateHTML("email-otp.tmpl", user.Language, map[string]any{
			"Code": secret,
			"User": user,
			"IP":   ip,
			"Date": time.Now().Format(time.DateTime),
		})

		err = gap.Px.PushEmail(pushkit.EmailDeliverRequest{
			To: user.GetPrimaryEmail().Content,
			Email: pushkit.EmailData{
				Subject: subject,
				HTML:    &content,
			},
		})
		if err != nil {
			log.Warn().Err(err).Uint("factor", factor.ID).Msg("Failed to delivery one-time-password via mail...")
			return true, nil
		}
		return true, nil
	default:
		return false, nil
	}
}

func CheckFactor(factor models.AuthFactor, code string) error {
	switch factor.Type {
	case models.PasswordAuthFactor:
		return lo.Ternary(
			VerifyPassword(code, factor.Secret),
			nil,
			fmt.Errorf("invalid password"),
		)
	case models.TimeOtpFactor:
		return lo.Ternary(
			totp.Validate(code, factor.Secret),
			nil,
			fmt.Errorf("invalid verification code"),
		)
	case models.InAppNotifyFactor:
	case models.EmailPasswordFactor:
		identifier := fmt.Sprintf("%s%d", gap.FactorOtpPrefix, factor.ID)
		sub, err := gap.Jt.PullSubscribe(identifier, "otp_validator", nats.BindStream("OTPs"))
		if err != nil {
			log.Error().Err(err).Msg("Error subscribing to subject when validating factor code...")
			return fmt.Errorf("error subscribing to subject: %v", err)
		}
		defer sub.Unsubscribe()

		msgs, err := sub.Fetch(1, nats.MaxWait(3*time.Second))
		if err != nil {
			log.Error().Err(err).Msg("Error fetching message when validating factor code...")
			return fmt.Errorf("error fetching message: %v", err)
		}

		if len(msgs) > 0 {
			msg := msgs[0]
			if !strings.EqualFold(code, string(msg.Data)) {
				return fmt.Errorf("invalid verification code")
			}
			log.Info().Uint("factor", factor.ID).Str("secret", code).Msg("Verified one-time-password...")
			if err := msg.AckSync(); err != nil {
				log.Warn().Err(err).Uint("factor", factor.ID).Msg("Failed to acknowledge message when validating factor code...")
			}
			return nil
		}

		return fmt.Errorf("one-time-password not found or expired")
	}

	return nil
}
