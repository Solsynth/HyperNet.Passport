package services

import (
	"fmt"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const ConfirmRegistrationTemplate = `Dear %s,

Thank you for choosing to register with %s. We are excited to welcome you to our community and appreciate your trust in us.

Your registration details have been successfully received, and you are now a valued member of %s. Here are the confirm link of your registration:

	%s

As a confirmed registered member, you will have access to all our services.
We encourage you to explore our services and take full advantage of the resources available to you.

Once again, thank you for choosing us. We look forward to serving you and hope you have a positive experience with us.

Best regards,
%s`

const ResetPasswordTemplate = `Dear %s,

We received a request to reset the password for your account at %s. If you did not request a password reset, please ignore this email.

To confirm your password reset request and create a new password, please use the link below:

%s

This link will expire in 24 hours. If you do not reset your password within this time frame, you will need to submit a new password reset request.

If you have any questions or need further assistance, please do not hesitate to contact our support team.

Best regards,
%s`

const DeleteAccountTemplate = `Dear %s,

We received a request to delete your account at %s. If you did not request a account deletion, please change your account password right now.
If you changed your mind, please ignore this email.

To confirm your account deletion request, please use the link below:

%s

This link will expire in 24 hours. If you do not use that link within this time frame, you will need to submit an account deletion request.

If you have any questions or need further assistance, please do not hesitate to contact our support team.
Also, if you want to let us know why you decided to delete your account, send email us (lily@solsynth.dev) and tell us how could we improve our user experience.

Best regards,
%s`

func ValidateMagicToken(code string, mode models.MagicTokenType) (models.MagicToken, error) {
	var tk models.MagicToken
	if err := database.C.Where(models.MagicToken{Code: code, Type: mode}).First(&tk).Error; err != nil {
		return tk, err
	} else if tk.ExpiredAt != nil && time.Now().Unix() >= tk.ExpiredAt.Unix() {
		return tk, fmt.Errorf("token has been expired")
	}

	return tk, nil
}

func NewMagicToken(mode models.MagicTokenType, assignTo *models.Account, expiredAt *time.Time) (models.MagicToken, error) {
	var uid uint
	if assignTo != nil {
		uid = assignTo.ID
	}

	token := models.MagicToken{
		Code:      strings.Replace(uuid.NewString(), "-", "", -1),
		Type:      mode,
		AccountID: &uid,
		ExpiredAt: expiredAt,
	}

	if err := database.C.Save(&token).Error; err != nil {
		return token, err
	} else {
		return token, nil
	}
}

func NotifyMagicToken(token models.MagicToken) error {
	if token.AccountID == nil {
		return fmt.Errorf("could notify a non-assign magic token")
	}

	var user models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: *token.AccountID},
	}).Preload("Contacts").First(&user).Error; err != nil {
		return err
	}

	var subject string
	var content string
	switch token.Type {
	case models.ConfirmMagicToken:
		link := fmt.Sprintf("%s/flow/accounts/confirm?code=%s", viper.GetString("frontend_app"), token.Code)
		subject = fmt.Sprintf("[%s] Confirm your registration", viper.GetString("name"))
		content = fmt.Sprintf(
			ConfirmRegistrationTemplate,
			user.Name,
			viper.GetString("name"),
			viper.GetString("maintainer"),
			link,
			viper.GetString("maintainer"),
		)
	case models.ResetPasswordMagicToken:
		link := fmt.Sprintf("%s/flow/accounts/password-reset?code=%s", viper.GetString("frontend_app"), token.Code)
		subject = fmt.Sprintf("[%s] Reset your password", viper.GetString("name"))
		content = fmt.Sprintf(
			ResetPasswordTemplate,
			user.Name,
			viper.GetString("name"),
			link,
			viper.GetString("maintainer"),
		)
	case models.DeleteAccountMagicToken:
		link := fmt.Sprintf("%s/flow/accounts/account-delete?code=%s", viper.GetString("frontend_app"), token.Code)
		subject = fmt.Sprintf("[%s] Confirm your account deletion", viper.GetString("name"))
		content = fmt.Sprintf(
			DeleteAccountTemplate,
			user.Name,
			viper.GetString("name"),
			link,
			viper.GetString("maintainer"),
		)
	default:
		return fmt.Errorf("unsupported magic token type to notify")
	}

	err := gap.Px.PushEmail(pushkit.EmailDeliverRequest{
		To: user.GetPrimaryEmail().Content,
		Email: pushkit.EmailData{
			Subject: subject,
			Text:    &content,
		},
	})
	return err
}
