package admin

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/authkit/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
)

func listUser(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if err := exts.EnsureGrantedPerm(c, "AdminUser", true); err != nil {
		return err
	}

	var count int64
	if err := database.C.Model(&models.Account{}).Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var items []models.Account
	if err := database.C.Offset(offset).Limit(take).Order("id ASC").Find(&items).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func getUser(c *fiber.Ctx) error {
	userId, _ := c.ParamsInt("user")

	if err := exts.EnsureGrantedPerm(c, "AdminUser", true); err != nil {
		return err
	}

	var user models.Account
	if err := database.C.Where("id = ?", userId).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err))
	}

	return c.JSON(user)
}

func forceConfirmAccount(c *fiber.Ctx) error {
	userId, _ := c.ParamsInt("user")

	if err := exts.EnsureGrantedPerm(c, "AdminUserConfirmation", true); err != nil {
		return err
	}
	operator := c.Locals("user").(*sec.UserInfo)

	var user models.Account
	if err := database.C.Where("id = ?", userId).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err))
	}

	if err := services.ForceConfirmAccount(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddAuditRecord(operator.ID, "user.confirm", c.IP(), c.Get(fiber.HeaderUserAgent), map[string]any{
			"user_id": user.ID,
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
