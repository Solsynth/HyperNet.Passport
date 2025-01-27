package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

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
