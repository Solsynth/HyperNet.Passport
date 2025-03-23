package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func checkPermission(c *fiber.Ctx) error {
	var data struct {
		PermNode string `json:"perm_node" validate:"required"`
		Value    any    `json:"value" validate:"required"`
	}
	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	var heldPerms map[string]any
	rawHeldPerms, _ := jsoniter.Marshal(user.PermNodes)
	_ = jsoniter.Unmarshal(rawHeldPerms, &heldPerms)
	valid := services.HasPermNode(heldPerms, data.PermNode, data.Value)
	if !valid {
		return c.SendStatus(fiber.StatusForbidden)
	}
	return c.SendStatus(fiber.StatusOK)
}

func checkUserPermission(c *fiber.Ctx) error {
	var data struct {
		PermNode string `json:"perm_node" validate:"required"`
		Value    any    `json:"value" validate:"required"`
	}
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("userId")
	relation, err := services.GetRelationWithTwoNode(user.ID, uint(relatedId))
	if err != nil {
		return err
	}
	defaultPerm := relation.Status == models.RelationshipFriend
	valid := services.HasPermNodeWithDefault(relation.PermNodes, data.PermNode, data.Value, defaultPerm)
	if !valid {
		return c.SendStatus(fiber.StatusForbidden)
	}
	return c.SendStatus(fiber.StatusOK)
}
