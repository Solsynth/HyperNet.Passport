package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func setAvatar(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		AttachmentID string `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	user.Avatar = &data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "profile.edit.avatar", strconv.Itoa(int(user.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		services.InvalidAuthCacheWithUser(user.ID)
	}

	return c.SendStatus(fiber.StatusOK)
}

func setBanner(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		AttachmentID string `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	user.Banner = &data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "profile.edit.banner", strconv.Itoa(int(user.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		services.InvalidAuthCacheWithUser(user.ID)
	}

	return c.SendStatus(fiber.StatusOK)
}

func getAvatar(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if content := user.GetAvatar(); content == nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.Redirect(*content, fiber.StatusFound)
	}
}

func getBanner(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if content := user.GetBanner(); content == nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.Redirect(*content, fiber.StatusFound)
	}
}
