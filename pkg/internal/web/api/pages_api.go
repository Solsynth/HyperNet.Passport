package api

import (
	"strconv"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
)

func getAccountPage(c *fiber.Ctx) error {
	alias := c.Params("alias")

	tx := database.C.Where("name = ?", alias)

	numericId, err := strconv.Atoi(alias)
	if err == nil {
		tx = tx.Or("id = ?", numericId)
	}

	var account models.Account
	if err := tx.First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var page models.AccountPage
	if err := database.C.Where("account_id = ?", account.ID).First(&page).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(page)
}

func getOwnAccountPage(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var page models.AccountPage
	if err := database.C.Where("account_id = ?", user.ID).First(&page).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(page)
}

func updateAccountPage(c *fiber.Ctx) error {
	if err := sec.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Content string `json:"content" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var page models.AccountPage
	if err := database.C.Where("account_id = ?", user.ID).First(&page).Error; err != nil {
		page = models.AccountPage{
			AccountID: user.ID,
		}
	}

	page.Content = data.Content

	if err := database.C.Save(&page).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(page)
}
