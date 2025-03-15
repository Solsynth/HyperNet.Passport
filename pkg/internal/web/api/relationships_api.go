package api

import (
	"strconv"
	"strings"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func listRelationship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	statusQuery := c.Query("status")

	status := lo.Map(strings.Split(statusQuery, ","), func(s string, _ int) models.RelationshipStatus {
		num, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return models.RelationshipStatus(num)
	})

	var err error
	var friends []models.AccountRelationship
	if len(status) == 0 {
		if friends, err = services.ListAllRelationship(user); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		if friends, err = services.ListRelationshipWithFilter(user, status...); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(friends)
}

func getRelationship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if friend, err := services.GetRelationWithTwoNode(user.ID, related.ID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(friend)
	}
}

func editRelationship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	var data struct {
		Status    uint8          `json:"status"`
		PermNodes map[string]any `json:"perm_nodes"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	relationship, err := services.GetRelationWithTwoNode(user.ID, uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	relationship.Status = models.RelationshipStatus(data.Status)
	relationship.PermNodes = data.PermNodes

	if friendship, err := services.EditRelationship(relationship); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "relationships.edit", map[string]any{
			"relationship": relationship,
		}, c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(friendship)
	}
}

func deleteRelationship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	relationship, err := services.GetRelationWithTwoNode(user.ID, related.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteRelationship(relationship); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "relationships.delete", map[string]any{
			"relationship": relationship,
		}, c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(relationship)
	}
}

// Friends stuff

func makeFriendship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Related string `json:"related" validate:"required"`
	}
	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var related models.Account
	if numericId, err := strconv.Atoi(data.Related); err == nil {
		related, err = services.GetAccount(uint(numericId))
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	} else {
		related, err = services.LookupAccount(data.Related)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}

	friend, err := services.NewFriend(user, related)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "relationships.friends.new", map[string]any{
			"related": related,
		}, c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(friend)
	}
}

func makeBlockship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Related string `json:"related" validate:"required"`
	}
	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var related models.Account
	if numericId, err := strconv.Atoi(data.Related); err == nil {
		related, err = services.GetAccount(uint(numericId))
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	} else {
		related, err = services.LookupAccount(data.Related)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}

	friend, err := services.NewBlockship(user, related)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "relationships.blocks.new", map[string]any{
			"related": related,
		}, c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(friend)
	}
}

func acceptFriend(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.HandleFriend(user, related, true); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "relationships.friends.accept", map[string]any{
			"related": related,
		}, c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.SendStatus(fiber.StatusOK)
	}
}

func declineFriend(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.HandleFriend(user, related, false); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "relationships.friends.decline", map[string]any{
			"related": related,
		}, c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.SendStatus(fiber.StatusOK)
	}
}
