package api

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func listCheckInRecord(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var count int64
	if err := database.C.
		Model(&models.CheckInRecord{}).
		Where("account_id = ?", user.ID).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var records []models.CheckInRecord
	if err := database.C.
		Where("account_id = ?", user.ID).
		Limit(take).Offset(offset).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  records,
	})
}

func listOtherUserCheckInRecord(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	alias := c.Params("alias")

	var account models.Account
	if err := database.C.
		Where(&models.Account{Name: alias}).
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var count int64
	if err := database.C.
		Model(&models.CheckInRecord{}).
		Where("account_id = ?", account.ID).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var records []models.CheckInRecord
	if err := database.C.
		Where("account_id = ?", account.ID).
		Limit(take).Offset(offset).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  records,
	})
}

func getTodayCheckIn(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if record, err := services.GetTodayCheckIn(user); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(record)
	}
}

func doCheckIn(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if record, err := services.CheckIn(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "checkIn", strconv.Itoa(int(record.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(record)
	}
}
