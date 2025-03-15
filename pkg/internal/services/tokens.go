package services

import (
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/localize"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

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

func NotifyMagicToken(token models.MagicToken, skipCheck ...bool) error {
	if token.AccountID == nil {
		return fmt.Errorf("could notify a non-assign magic token")
	}
	if token.LastNotifiedAt != nil && (len(skipCheck) == 0 || !skipCheck[0]) {
		if token.LastNotifiedAt.Unix() >= time.Now().Add(-60*time.Minute).Unix() {
			return fmt.Errorf("magic token has been notified in an hour")
		}
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
		subject = fmt.Sprintf("[%s] %s", viper.GetString("name"), localize.L.GetLocalizedString("subjectConfirmRegistration", user.Language))
		content = localize.L.RenderLocalizedTemplateHTML("register-confirm.tmpl", user.Language, map[string]any{
			"User": user,
			"Link": link,
		})
	case models.ResetPasswordMagicToken:
		link := fmt.Sprintf("%s/flow/accounts/password-reset?code=%s", viper.GetString("frontend_app"), token.Code)
		subject = fmt.Sprintf("[%s] %s", viper.GetString("name"), localize.L.GetLocalizedString("subjectResetPassword", user.Language))
		content = localize.L.RenderLocalizedTemplateHTML("reset-password.tmpl", user.Language, map[string]any{
			"User": user,
			"Link": link,
		})
	case models.DeleteAccountMagicToken:
		link := fmt.Sprintf("%s/flow/accounts/deletion?code=%s", viper.GetString("frontend_app"), token.Code)
		subject = fmt.Sprintf("[%s] %s", viper.GetString("name"), localize.L.GetLocalizedString("subjectDeleteAccount", user.Language))
		content = localize.L.RenderLocalizedTemplateHTML("confirm-deletion.tmpl", user.Language, map[string]any{
			"User": user,
			"Link": link,
		})
	default:
		return fmt.Errorf("unsupported magic token type to notify")
	}

	err := gap.Px.PushEmail(pushkit.EmailDeliverRequest{
		To: user.GetPrimaryEmail().Content,
		Email: pushkit.EmailData{
			Subject: subject,
			HTML:    &content,
		},
	})

	if err == nil {
		database.C.Model(&token).Update("last_notified_at", time.Now())
	}

	return err
}
