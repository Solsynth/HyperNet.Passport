package api

import (
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
)

func getPunishment(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	data, err := services.GetPunishment(uint(id), true)
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func listUserPunishment(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	data, err := services.ListPunishments(user)
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func listMadePunishment(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	moderator := c.Locals("user").(models.Account)

	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if c.QueryBool("all", false) {
		if err := exts.EnsureGrantedPerm(c, "OverridePunishments", true); err != nil {
			return err
		}
		count, err := services.CountAllPunishments()
		data, err := services.ListAllPunishments(take, offset)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"count": count,
			"data":  data,
		})
	}

	count, err := services.CountMadePunishments(moderator)
	data, err := services.ListMadePunishments(moderator, take, offset)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"count": count,
		"data":  data,
	})
}

func createPunishment(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "CreatePunishments", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Reason    string         `json:"reason" validate:"required"`
		Type      int            `json:"type"`
		ExpiredAt *time.Time     `json:"expired_at"`
		PermNodes map[string]any `json:"perm_nodes"`
		AccountID uint           `json:"account_id"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var account models.Account
	if err := database.C.Where("id = ?", data.AccountID).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	punishment := models.Punishment{
		Reason:    data.Reason,
		Type:      data.Type,
		PermNodes: data.PermNodes,
		ExpiredAt: data.ExpiredAt,
		Account:   account,
		AccountID: account.ID,
	}

	if punishment, err := services.NewPunishment(punishment, user); err != nil {
		return err
	} else {
		return c.JSON(punishment)
	}
}

func editPunishment(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	id, _ := c.ParamsInt("id", 0)

	var data struct {
		Reason    string         `json:"reason" validate:"required"`
		Type      int            `json:"type"`
		ExpiredAt *time.Time     `json:"expired_at"`
		PermNodes map[string]any `json:"perm_nodes"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var punishment models.Punishment
	if c.QueryBool("override", false) {
		if err = exts.EnsureGrantedPerm(c, "OverridePunishments", true); err != nil {
			return err
		}
		punishment, err = services.GetPunishment(uint(id))
	} else {
		punishment, err = services.GetMadePunishment(uint(id), user)
	}
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	punishment.Reason = data.Reason
	punishment.Type = data.Type
	punishment.ExpiredAt = data.ExpiredAt
	punishment.PermNodes = data.PermNodes

	if punishment, err := services.EditPunishment(punishment); err != nil {
		return err
	} else {
		return c.JSON(punishment)
	}
}

func deletePunishment(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	id := c.QueryInt("id")

	var err error
	var punishment models.Punishment
	if c.QueryBool("override", false) {
		if err = exts.EnsureGrantedPerm(c, "OverridePunishments", true); err != nil {
			return err
		}
		punishment, err = services.GetPunishment(uint(id))
	} else {
		punishment, err = services.GetMadePunishment(uint(id), user)
	}
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeletePunishment(punishment); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
