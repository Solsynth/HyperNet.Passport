package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getAuthPreference(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	cfg, err := services.GetAuthPreference(user)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(cfg.Config.Data())
}

func updateAuthPreference(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data models.AuthConfig
	if err := exts.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	cfg, err := services.UpdateAuthPreference(user, data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "preferences.edit", "auth", c.IP(), c.Get(fiber.HeaderUserAgent))
	}

	return c.JSON(cfg.Config.Data())
}

func getNotificationPreference(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	notification, err := services.GetNotificationPreference(user)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(notification)
}

func updateNotificationPreference(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Config map[string]bool `json:"config"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	notification, err := services.UpdateNotificationPreference(user, data.Config)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "preferences.edit", "notifications", c.IP(), c.Get(fiber.HeaderUserAgent))
	}

	return c.JSON(notification)
}
