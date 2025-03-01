package admin

import (
	"fmt"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func grantBadge(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "AdminGrantBadges", true); err != nil {
		return err
	}

	var data struct {
		Type      string         `json:"type" validate:"required"`
		Metadata  map[string]any `json:"metadata"`
		AccountID uint           `json:"account_id"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var account models.Account
	if account, err = services.GetAccount(data.AccountID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("target account was not found: %v", err))
	}

	badge := models.Badge{
		Type:     data.Type,
		Metadata: data.Metadata,
	}

	if err := services.GrantBadge(account, badge); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(badge)
	}
}

func revokeBadge(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "AdminRevokeBadges", true); err != nil {
		return err
	}

	id, _ := c.ParamsInt("badgeId", 0)

	var badge models.Badge
	if err := database.C.Where("id = ?", id).First(&badge).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("target badge was not found: %v", err))
	}

	if err := services.RevokeBadge(badge); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(badge)
	}
}
