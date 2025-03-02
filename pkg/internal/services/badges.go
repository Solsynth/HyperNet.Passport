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

func ActiveBadge(badge models.Badge) error {
	accountId := badge.AccountID
	tx := database.C.Begin()

	if err := tx.Model(&badge).Where("account_id = ?", accountId).Update("is_active", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&badge).Where("id = ?", badge.ID).Update("is_active", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
