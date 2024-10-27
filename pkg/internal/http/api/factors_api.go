package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
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

func requestResetPassword(c *fiber.Ctx) error {
	var data struct {
		UserID uint `json:"user_id" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	user, err := services.GetAccount(data.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = services.CheckAbleToResetPassword(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if err = services.RequestResetPassword(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func confirmResetPassword(c *fiber.Ctx) error {
	var data struct {
		Code        string `json:"code" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := services.ConfirmResetPassword(data.Code, data.NewPassword); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
