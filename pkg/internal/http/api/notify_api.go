package api

import (
	"fmt"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func notifyUser(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "DevNotifyUser", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		ClientID   string         `json:"client_id" validate:"required"`
		Topic      string         `json:"type" validate:"required"`
		Title      string         `json:"subject" validate:"required,max=1024"`
		Subtitle   string         `json:"subtitle" validate:"max=1024"`
		Body       string         `json:"content" validate:"required,max=4096"`
		Metadata   map[string]any `json:"metadata"`
		Priority   int            `json:"priority"`
		IsRealtime bool           `json:"is_realtime"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	client, err := services.GetThirdClientWithUser(data.ClientID, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get client: %v", err))
	}

	userId, _ := c.ParamsInt("user")

	var target models.Account
	if target, err = services.GetAccount(uint(userId)); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	notification := models.Notification{
		Topic:     data.Topic,
		Subtitle:  data.Subtitle,
		Title:     data.Title,
		Body:      data.Body,
		Metadata:  data.Metadata,
		Priority:  data.Priority,
		Account:   target,
		AccountID: target.ID,
		SenderID:  &client.ID,
	}

	if data.IsRealtime {
		if err := services.PushNotification(notification); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else {
		if err := services.NewNotification(notification); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func notifyAllUser(c *fiber.Ctx) error {
	if err := exts.EnsureGrantedPerm(c, "DevNotifyAllUser", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		ClientID   string         `json:"client_id" validate:"required"`
		Topic      string         `json:"type" validate:"required"`
		Title      string         `json:"subject" validate:"required,max=1024"`
		Subtitle   string         `json:"subtitle" validate:"max=1024"`
		Body       string         `json:"content" validate:"required,max=4096"`
		Metadata   map[string]any `json:"metadata"`
		Priority   int            `json:"priority"`
		IsRealtime bool           `json:"is_realtime"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	client, err := services.GetThirdClientWithUser(data.ClientID, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get client: %v", err))
	}

	var accounts []models.Account
	if err := database.C.Find(&accounts).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var notifications []models.Notification
	for _, account := range accounts {
		notification := models.Notification{
			Topic:     data.Topic,
			Subtitle:  data.Subtitle,
			Title:     data.Title,
			Body:      data.Body,
			Metadata:  data.Metadata,
			Priority:  data.Priority,
			Account:   account,
			AccountID: account.ID,
			SenderID:  &client.ID,
		}
		notifications = append(notifications, notification)
	}

	if data.IsRealtime {
		go services.PushNotificationBatch(notifications)
	} else {
		go services.NewNotificationBatch(notifications)
	}

	return c.SendStatus(fiber.StatusOK)
}
