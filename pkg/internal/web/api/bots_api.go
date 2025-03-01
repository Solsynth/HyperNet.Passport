package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"strings"
	"time"
)

func listBots(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	tx := database.C.Where("automated_id = ?", user.ID)

	countTx := tx
	var count int64
	if err := countTx.Model(&models.Account{}).Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var bots []models.Account
	if err := tx.Find(&bots).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  bots,
	})
}

func createBot(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	cnt, _ := services.GetBotCount(user)
	if err := exts.EnsureGrantedPerm(c, "CreateBots", cnt+1); err != nil {
		return err
	}

	var data struct {
		Name        string `json:"name" validate:"required,lowercase,alphanum,min=4,max=16"`
		Nick        string `json:"nick" validate:"required"`
		Description string `json:"description"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	} else {
		data.Name = strings.TrimSpace(data.Name)
		data.Nick = strings.TrimSpace(data.Nick)
	}

	if !services.ValidateAccountName(data.Nick, 4, 24) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid bot nick, length requires 4 to 24")
	}

	bot, err := services.NewBot(user, models.Account{
		Name:        data.Name,
		Nick:        data.Nick,
		Description: data.Description,
		ConfirmedAt: lo.ToPtr(time.Now()),
		PermNodes:   datatypes.JSONMap{},
	})

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(bot)
	}
}

func deleteBot(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	id, _ := c.ParamsInt("botId", 0)

	var bot models.Account
	if err := database.C.Where("id = ? AND automated_id = ?", id, user.ID).First(&bot).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteAccount(bot.ID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(bot)
}
