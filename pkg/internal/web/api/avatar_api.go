package api

import (
	"strconv"

	"git.solsynth.dev/hypernet/paperclip/pkg/filekit"
	"git.solsynth.dev/hypernet/paperclip/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
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

	og := user.Avatar
	user.Avatar = &data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "profile.edit.avatar", strconv.Itoa(int(user.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		services.InvalidAuthCacheWithUser(user.ID)
	}

	if og != nil && len(*og) > 0 {
		filekit.CountAttachmentUsage(gap.Nx, &proto.UpdateUsageRequest{
			Rid:   []string{*og},
			Delta: -1,
		})
	}
	filekit.CountAttachmentUsage(gap.Nx, &proto.UpdateUsageRequest{
		Rid:   []string{*user.Avatar},
		Delta: 1,
	})

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

	og := user.Banner
	user.Banner = &data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "profile.edit.banner", strconv.Itoa(int(user.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		services.InvalidAuthCacheWithUser(user.ID)
	}

	if og != nil && len(*og) > 0 {
		filekit.CountAttachmentUsage(gap.Nx, &proto.UpdateUsageRequest{
			Rid:   []string{*og},
			Delta: -1,
		})
	}
	filekit.CountAttachmentUsage(gap.Nx, &proto.UpdateUsageRequest{
		Rid:   []string{*user.Banner},
		Delta: 1,
	})

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
