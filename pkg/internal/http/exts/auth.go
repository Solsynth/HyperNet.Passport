package exts

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
)

func EnsureAuthenticated(c *fiber.Ctx) error {
	if _, ok := c.Locals("user").(*sec.UserInfo); !ok {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	return nil
}

func EnsureGrantedPerm(c *fiber.Ctx, key string, val any) error {
	if err := EnsureAuthenticated(c); err != nil {
		return err
	}
	perms := c.Locals("permissions").(map[string]any)
	if !services.HasPermNode(perms, key, val) {
		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("missing permission: %s", key))
	}
	return nil
}
