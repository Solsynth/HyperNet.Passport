package admin

import (
	"fmt"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
)

func editUserPermission(c *fiber.Ctx) error {
	userId, _ := c.ParamsInt("user")

	if err := exts.EnsureGrantedPerm(c, "AdminUserPermission", true); err != nil {
		return err
	}
	operator := c.Locals("user").(models.Account)

	var data struct {
		PermNodes map[string]any `json:"perm_nodes" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var user models.Account
	if err := database.C.Where("id = ?", userId).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err))
	}

	prev := user.PermNodes
	user.PermNodes = data.PermNodes

	services.InvalidAuthCacheWithUser(user.ID)

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddAuditRecord(operator, "user.permissions.edit", c.IP(), c.Get(fiber.HeaderUserAgent), map[string]any{
			"user_id":              user.ID,
			"previous_permissions": prev,
			"new_permissions":      data.PermNodes,
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
