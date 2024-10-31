package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func getRealm(c *fiber.Ctx) error {
	alias := c.Params("realm")
	if realm, err := services.GetRealmWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(realm)
	}
}

func listCommunityRealm(c *fiber.Ctx) error {
	realms, err := services.ListCommunityRealm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(realms)
}

func listOwnedRealm(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	if realms, err := services.ListOwnedRealm(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(realms)
	}
}

func listAvailableRealm(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	if realms, err := services.ListAvailableRealm(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(realms)
	}
}

func createRealm(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "CreateRealms", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Alias        string         `json:"alias" validate:"required,lowercase,min=4,max=32"`
		Name         string         `json:"name" validate:"required"`
		Description  string         `json:"description"`
		Avatar       *string        `json:"avatar"`
		Banner       *string        `json:"banner"`
		AccessPolicy map[string]any `json:"access_policy"`
		IsPublic     bool           `json:"is_public"`
		IsCommunity  bool           `json:"is_community"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	realm, err := services.NewRealm(models.Realm{
		Alias:        data.Alias,
		Name:         data.Name,
		Description:  data.Description,
		Avatar:       data.Avatar,
		Banner:       data.Banner,
		AccessPolicy: data.AccessPolicy,
		IsPublic:     data.IsPublic,
		IsCommunity:  data.IsCommunity,
		AccountID:    user.ID,
	}, user)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "realms.new", strconv.Itoa(int(realm.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
	}

	return c.JSON(realm)
}

func editRealm(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	id, _ := c.ParamsInt("realmId", 0)

	var data struct {
		Alias        string         `json:"alias" validate:"required,lowercase,min=4,max=32"`
		Name         string         `json:"name" validate:"required"`
		Description  string         `json:"description"`
		Avatar       *string        `json:"avatar"`
		Banner       *string        `json:"banner"`
		AccessPolicy map[string]any `json:"access_policy"`
		IsPublic     bool           `json:"is_public"`
		IsCommunity  bool           `json:"is_community"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	realm.Alias = data.Alias
	realm.Name = data.Name
	realm.Description = data.Description
	realm.Avatar = data.Avatar
	realm.Banner = data.Banner
	realm.AccessPolicy = data.AccessPolicy
	realm.IsPublic = data.IsPublic
	realm.IsCommunity = data.IsCommunity

	realm, err := services.EditRealm(realm)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "realms.edit", strconv.Itoa(int(realm.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
	}

	return c.JSON(realm)
}

func deleteRealm(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	id, _ := c.ParamsInt("realmId", 0)

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteRealm(realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "realms.delete", strconv.Itoa(int(realm.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
	}

	return c.SendStatus(fiber.StatusOK)
}
