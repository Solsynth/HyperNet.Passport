package services

import (
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const EmailPasswordTemplate = `Dear %s,

We hope this message finds you well.
As part of our ongoing commitment to ensuring the security of your account, we require you to complete the login process by entering the verification code below:

Your Login Verification Code: %s

Please use the provided code within the next 2 hours to complete your login. 
If you did not request this code, please update your information, maybe your username or email has been leak.

Thank you for your cooperation in helping us maintain the security of your account.

Best regards,
%s`

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

func GetFactorCode(factor models.AuthFactor) (bool, error) {
	switch factor.Type {
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

		subject := fmt.Sprintf("[%s] Login verification code", viper.GetString("name"))
		content := fmt.Sprintf(EmailPasswordTemplate, user.Name, secret, viper.GetString("maintainer"))

		err = gap.Px.PushEmail(pushkit.EmailDeliverRequest{
			To: user.GetPrimaryEmail().Content,
			Email: pushkit.EmailData{
				Subject: subject,
				Text:    &content,
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
	case models.EmailPasswordFactor:
		identifier := fmt.Sprintf("%s%d", gap.FactorOtpPrefix, factor.ID)
		sub, err := gap.Jt.PullSubscribe(identifier, "otp_validator", nats.BindStream("OTPs"))
		if err != nil {
			log.Error().Err(err).Msg("Error subscribing to subject when validating factor code...")
			return fmt.Errorf("error subscribing to subject: %v", err)
		}
		// ChatGPT said the subscription will be reused, so we don't need to unsubscribe
		// defer sub.Unsubscribe()

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
			msg.Ack()
			return nil
		}

		return fmt.Errorf("one-time-password not found or expired")
	}

	return nil
}
