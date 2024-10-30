package admin

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func notifyAllUser(c *fiber.Ctx) error {
	var data struct {
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

	if err := exts.EnsureGrantedPerm(c, "AdminNotifyAll", true); err != nil {
		return err
	}
	operator := c.Locals("user").(*sec.UserInfo)

	var users []models.Account
	if err := database.C.Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddAuditRecord(operator.ID, "notify.all", c.IP(), c.Get(fiber.HeaderUserAgent), map[string]any{
			"payload": data,
		})
	}

	go func() {
		for _, user := range users {
			notification := models.Notification{
				Topic:     data.Topic,
				Subtitle:  data.Subtitle,
				Title:     data.Title,
				Body:      data.Body,
				Metadata:  data.Metadata,
				Priority:  data.Priority,
				Account:   user,
				AccountID: user.ID,
			}

			if data.IsRealtime {
				if err := services.PushNotification(notification); err != nil {
					log.Error().Err(err).Uint("user", user.ID).Msg("Failed to push notification...")
				}
			} else {
				if err := services.NewNotification(notification); err != nil {
					log.Error().Err(err).Uint("user", user.ID).Msg("Failed to create notification...")
				}
			}
		}
	}()

	return c.SendStatus(fiber.StatusOK)
}

func notifyOneUser(c *fiber.Ctx) error {
	var data struct {
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

	if err := exts.EnsureGrantedPerm(c, "AdminNotifyAll", true); err != nil {
		return err
	}
	operator := c.Locals("user").(*sec.UserInfo)

	var user models.Account
	if err := database.C.Where("id = ?", data.UserID).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddAuditRecord(operator.ID, "notify.one", c.IP(), c.Get(fiber.HeaderUserAgent), map[string]any{
			"user_id": user.ID,
			"payload": data,
		})
	}

	notification := models.Notification{
		Topic:     data.Topic,
		Subtitle:  data.Subtitle,
		Title:     data.Title,
		Body:      data.Body,
		Priority:  data.Priority,
		AccountID: user.ID,
	}

	if data.IsRealtime {
		if err := services.PushNotification(notification); err != nil {
			log.Error().Err(err).Uint("user", user.ID).Msg("Failed to push notification...")
		}
	} else {
		if err := services.NewNotification(notification); err != nil {
			log.Error().Err(err).Uint("user", user.ID).Msg("Failed to create notification...")
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
