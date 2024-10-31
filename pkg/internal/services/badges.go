package services

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
)

func GrantBadge(user models.Account, badge models.Badge) error {
	badge.AccountID = user.ID
	return database.C.Save(badge).Error
}

func RevokeBadge(badge models.Badge) error {
	return database.C.Delete(&badge).Error
}
