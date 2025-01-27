package api

import (
	"fmt"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getAvailableFactors(c *fiber.Ctx) error {
	ticketId := c.QueryInt("ticketId", 0)
	if ticketId <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "must provide ticket id as a query parameter")
	}

	ticket, err := services.GetTicket(uint(ticketId))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("ticket was not found: %v", err))
	}
	factors, err := services.ListUserFactor(ticket.AccountID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(factors)
}

func requestFactorToken(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("factorId", 0)

	factor, err := services.GetFactor(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if sent, err := services.GetFactorCode(factor); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if !sent {
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func listFactor(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var factors []models.AuthFactor
	if err := database.C.Where("account_id = ?", user.ID).Find(&factors).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(factors)
}

func createFactor(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Type   models.AuthFactorType `json:"type"`
		Secret string                `json:"secret"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	typeWhitelist := []models.AuthFactorType{
		models.EmailPasswordFactor,
		models.InAppNotifyFactor,
		models.TimeOtpFactor,
	}
	if !lo.Contains(typeWhitelist, data.Type) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid factor type")
	}

	// Currently, each type of factor can only be created once
	var currentCount int64
	if err := database.C.Model(&models.AuthFactor{}).
		Where("account_id = ? AND type = ?", user.ID, data.Type).
		Count(&currentCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("unable to check current factor count: %v", err))
	} else if currentCount > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "this type of factor already exists")
	}

	factor := models.AuthFactor{
		Type:      data.Type,
		Secret:    data.Secret,
		Account:   user,
		AccountID: user.ID,
	}
	if err := database.C.Create(&factor).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(factor)
}

func deleteFactor(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("factorId", 0)

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var factor models.AuthFactor
	if err := database.C.Where("id = ? AND account_id = ?", id, user.ID).First(&factor).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if factor.Type == models.PasswordAuthFactor {
		return fiber.NewError(fiber.StatusBadRequest, "unable to delete password factor")
	}

	if err := database.C.Delete(&factor).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
