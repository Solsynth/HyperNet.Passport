package exts

import (
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func EnsureAuthenticated(c *fiber.Ctx) error {
	if _, ok := c.Locals("nex_user").(*sec.UserInfo); !ok {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	return nil
}

func EnsureGrantedPerm(c *fiber.Ctx, key string, val any) error {
	if err := EnsureAuthenticated(c); err != nil {
		return err
	}
	perms := c.Locals("nex_user").(*sec.UserInfo).PermNodes
	if !services.HasPermNode(perms, key, val) {
		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("missing permission: %s", key))
	}
	return nil
}
