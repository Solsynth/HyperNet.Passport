package api

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func listUserBadge(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var badges []models.Badge
	if err := database.C.Where("account_id = ?", user.ID).Find(&badges).Error; err != nil {
		return err
	}

	return c.JSON(badges)
}

func activeUserBadge(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	badgeId, _ := c.ParamsInt("badgeId", 0)

	var badge models.Badge
	if err := database.C.Where("id = ? AND account_id = ?", badgeId, user.ID).First(&badge).Error; err != nil {
		return err
	}

	if err := services.ActiveBadge(badge); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
