package services

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
)

func GetTicket(id uint) (models.AuthTicket, error) {
	var ticket models.AuthTicket
	if err := database.C.
		Where(&models.AuthTicket{BaseModel: models.BaseModel{ID: id}}).
		First(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func GetTicketWithToken(tokenId string) (models.AuthTicket, error) {
	var ticket models.AuthTicket
	if err := database.C.
		Where(models.AuthTicket{AccessToken: &tokenId}).
		Or(models.AuthTicket{RefreshToken: &tokenId}).
		First(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}
