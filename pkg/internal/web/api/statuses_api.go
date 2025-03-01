package api

import (
	"fmt"
	"strconv"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getStatus(c *fiber.Ctx) error {
	alias := c.Params("alias")

	var user models.Account
	if err := database.C.Where(models.Account{
		Name: alias,
	}).Preload("Profile").First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("account not found: %s", err))
	}

	status, err := services.GetStatus(user.ID)
	disturbable := services.GetStatusDisturbable(user.ID) == nil
	online := services.GetStatusOnline(user.ID) == nil

	return c.JSON(fiber.Map{
		"status":         lo.Ternary(err == nil, &status, nil),
		"last_seen_at":   user.Profile.LastSeenAt,
		"is_disturbable": disturbable,
		"is_online":      online,
	})
}

func getMyselfStatus(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	status, err := services.GetStatus(user.ID)
	disturbable := services.GetStatusDisturbable(user.ID) == nil
	online := services.GetStatusOnline(user.ID) == nil

	return c.JSON(fiber.Map{
		"status":         lo.Ternary(err == nil, &status, nil),
		"last_seen_at":   user.Profile.LastSeenAt,
		"is_disturbable": disturbable,
		"is_online":      online,
	})
}

func setStatus(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var req struct {
		Type        string     `json:"type" validate:"required"`
		Label       string     `json:"label" validate:"required"`
		Attitude    uint       `json:"attitude"`
		IsNoDisturb bool       `json:"is_no_disturb"`
		IsInvisible bool       `json:"is_invisible"`
		ClearAt     *time.Time `json:"clear_at"`
	}

	if err := exts.BindAndValidate(c, &req); err != nil {
		return err
	}

	// End the status already exists
	if status, err := services.GetStatus(user.ID); err == nil {
		status.ClearAt = lo.ToPtr(time.Now())
		database.C.Save(&status)
	}

	status := models.Status{
		Type:        req.Type,
		Label:       req.Label,
		Attitude:    models.StatusAttitude(req.Attitude),
		IsNoDisturb: req.IsNoDisturb,
		IsInvisible: req.IsInvisible,
		ClearAt:     req.ClearAt,
		AccountID:   user.ID,
	}

	if status, err := services.NewStatus(user, status); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "statuses.set", strconv.Itoa(int(status.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(status)
	}
}

func editStatus(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var req struct {
		Type        string     `json:"type" validate:"required"`
		Label       string     `json:"label" validate:"required"`
		Attitude    uint       `json:"attitude"`
		IsNoDisturb bool       `json:"is_no_disturb"`
		IsInvisible bool       `json:"is_invisible"`
		ClearAt     *time.Time `json:"clear_at"`
	}

	if err := exts.BindAndValidate(c, &req); err != nil {
		return err
	}

	status, err := services.GetStatus(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "you must set a status first and then can edit it")
	}

	status.Type = req.Type
	status.Label = req.Label
	status.Attitude = models.StatusAttitude(req.Attitude)
	status.IsNoDisturb = req.IsNoDisturb
	status.IsInvisible = req.IsInvisible
	status.ClearAt = req.ClearAt

	if status, err := services.EditStatus(user, status); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		services.AddEvent(user.ID, "statuses.edit", strconv.Itoa(int(status.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.JSON(status)
	}
}

func clearStatus(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if err := services.ClearStatus(user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "statuses.clear", strconv.Itoa(int(user.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
	}

	return c.SendStatus(fiber.StatusOK)
}
