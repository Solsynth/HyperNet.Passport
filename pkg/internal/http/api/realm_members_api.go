package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func listRealmMembers(c *fiber.Ctx) error {
	alias := c.Params("realm")
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if realm, err := services.GetRealmWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if count, err := services.CountRealmMember(realm.ID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if members, err := services.ListRealmMember(realm.ID, take, offset); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(fiber.Map{
			"count": count,
			"data":  members,
		})
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
		Related string `json:"related" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	var numericId int
	if numericId, err = strconv.Atoi(data.Related); err == nil {
		err = database.C.Where(&models.Account{
			BaseModel: models.BaseModel{ID: uint(numericId)},
		}).First(&account).Error
	} else {
		err = database.C.Where(&models.Account{
			Name: data.Related,
		}).First(&account).Error
	}
	if err != nil {
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
	memberId, _ := c.ParamsInt("memberId", 0)

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var member models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		BaseModel: models.BaseModel{ID: uint(memberId)},
		RealmID:   realm.ID,
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

	if err := database.C.Delete(&member).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
