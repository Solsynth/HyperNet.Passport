package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/authkit/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"github.com/samber/lo"
)

func GetUserAccountGroup(user models.Account) ([]models.AccountGroup, error) {
	var members []models.AccountGroupMember
	if err := database.C.Where(&models.AccountGroupMember{
		AccountID: user.ID,
	}).Find(&members).Error; err != nil {
		return nil, err
	}

	var groups []models.AccountGroup
	if err := database.C.Where("id IN ?", lo.Map(members, func(item models.AccountGroupMember, index int) uint {
		return item.GroupID
	})).Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}
