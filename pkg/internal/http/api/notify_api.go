package api

import (
	"fmt"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
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
		UserID     uint           `json:"user_id" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	client, err := services.GetThirdClientWithUser(data.ClientID, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get client: %v", err))
	}

	var target models.Account
	if target, err = services.GetAccount(data.UserID); err != nil {
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
