package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/authkit/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
)

func GetBotCount(user models.Account) (int64, error) {
	var count int64
	if err := database.C.Where("automated_id = ?", user.ID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func NewBot(user models.Account, bot models.Account) (models.Account, error) {
	bot.AutomatedBy = &user
	bot.AutomatedID = &user.ID

	if err := database.C.Save(&bot).Error; err != nil {
		return bot, err
	}
	return bot, nil
}
