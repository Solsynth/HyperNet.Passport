package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
)

func getTickets(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var tickets []models.AuthTicket
	if err := database.C.
		Where(&models.AuthTicket{AccountID: user.ID}).
		Model(&models.AuthTicket{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Order("created_at desc").
		Where(&models.AuthTicket{AccountID: user.ID}).
		Limit(take).
		Offset(offset).
		Find(&tickets).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  tickets,
	})
}

func killTicket(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)
	id, _ := c.ParamsInt("ticketId", 0)

	if err := database.C.Delete(&models.AuthTicket{}, &models.AuthTicket{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
