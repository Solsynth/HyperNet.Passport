package api

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"github.com/gofiber/fiber/v2"
)

func listContact(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var contacts []models.AccountContact
	if err := database.C.Where("account_id = ?", user.ID).Find(&contacts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(contacts)
}

func getContact(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	contactId, _ := c.ParamsInt("contactId")

	var contact models.AccountContact
	if err := database.C.Where("account_id = ? AND id = ?", user.ID, contactId).First(&contact).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(contact)
}

func createContact(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Type int8 `json:"type"`
		Content string `json:"content" validate:"required"`
		IsPublic bool `json:"is_public"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	contact := models.AccountContact{
        Type:       data.Type,
        Content:    data.Content,
        IsPublic:   data.IsPublic,
        IsPrimary:  false,
        VerifiedAt: nil,
        AccountID:  user.ID,
    }
	if err := database.C.Create(&contact).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(contact)
}

func updateContact(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	contactId, _ := c.ParamsInt("contactId")

	var data struct {
		Type int8 `json:"type"`
		Content string `json:"content" validate:"required"`
		IsPublic bool `json:"is_public"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var contact models.AccountContact
	if err := database.C.Where("account_id = ? AND id = ?", user.ID, contactId).First(&contact).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	contact.Type = data.Type
	contact.IsPublic = data.IsPublic
	if contact.Content != data.Content {
		contact.Content = data.Content
		contact.VerifiedAt = nil
	}

	if err := database.C.Save(&contact).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(contact)
}

func deleteContact(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	contactId, _ := c.ParamsInt("contactId")

	var contact models.AccountContact
	if err := database.C.Where("account_id = ? AND id = ?", user.ID, contactId).First(&contact).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := database.C.Delete(&contact).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}