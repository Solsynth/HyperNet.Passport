package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func setAvatar(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	var data struct {
		AttachmentID string `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := database.C.Where("id = ?", user.ID).Updates(&models.Account{Avatar: &data.AttachmentID}).Error; err != nil {
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
	user := c.Locals("user").(*sec.UserInfo)

	var data struct {
		AttachmentID string `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := database.C.Where("id = ?", user.ID).Updates(&models.Account{Banner: &data.AttachmentID}).Error; err != nil {
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
	user := c.Locals("user").(*sec.UserInfo)

	var account models.Account
	if err := database.C.Where("id = ?", user.ID).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err))
	}

	if content := account.GetAvatar(); content == nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.Redirect(*content, fiber.StatusFound)
	}
}

func getBanner(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	var account models.Account
	if err := database.C.Where("id = ?", user.ID).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err))
	}

	if content := account.GetBanner(); content == nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.Redirect(*content, fiber.StatusFound)
	}
}
