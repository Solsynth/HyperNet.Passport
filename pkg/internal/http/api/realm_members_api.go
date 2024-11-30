package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func listRealmMembers(c *fiber.Ctx) error {
	alias := c.Params("realm")

	if realm, err := services.GetRealmWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if members, err := services.ListRealmMember(realm.ID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(members)
	}
}

func getMyRealmMember(c *fiber.Ctx) error {
	alias := c.Params("realm")
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if realm, err := services.GetRealmWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if member, err := services.GetRealmMember(user.ID, realm.ID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(member)
	}
}

func addRealmMember(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	alias := c.Params("realm")

	var data struct {
		Target string `json:"target" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		Name: data.Target,
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.AddRealmMember(user, account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func removeRealmMember(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	alias := c.Params("realm")

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		Name: data.Target,
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var member models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		RealmID:   realm.ID,
		AccountID: account.ID,
	}).First(&member).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.RemoveRealmMember(user, member, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func leaveRealm(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	alias := c.Params("realm")

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if user.ID == realm.AccountID {
		return fiber.NewError(fiber.StatusBadRequest, "you cannot leave your own realm")
	}

	var member models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		RealmID:   realm.ID,
		AccountID: user.ID,
	}).First(&member).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.RemoveRealmMember(user, member, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
