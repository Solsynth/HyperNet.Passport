package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/http/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"strconv"
	"time"
)

func getNotifications(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	tx := database.C.Where(&models.Notification{AccountID: user.ID}).Model(&models.Notification{})

	var count int64
	if err := tx.
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var notifications []models.Notification
	if err := tx.
		Limit(take).
		Offset(offset).
		Order("read_at DESC, created_at DESC").
		Find(&notifications).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  notifications,
	})
}

func markNotificationRead(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)
	id, _ := c.ParamsInt("notificationId", 0)

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}

	var notify models.Notification
	if err := database.C.Where(&models.Notification{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).First(&notify).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	notify.ReadAt = lo.ToPtr(time.Now())

	if err := database.C.Save(&notify).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "notifications.mark.read", strconv.Itoa(int(notify.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.SendStatus(fiber.StatusOK)
	}
}

func markNotificationReadBatch(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	var data struct {
		MessageIDs []uint `json:"messages"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := database.C.Model(&models.Notification{}).
		Where("account_id = ? AND id IN ?", user.ID, data.MessageIDs).
		Updates(&models.Notification{ReadAt: lo.ToPtr(time.Now())}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddEvent(user.ID, "notifications.markAll.read", strconv.Itoa(int(user.ID)), c.IP(), c.Get(fiber.HeaderUserAgent))
		return c.SendStatus(fiber.StatusOK)
	}
}

func getNotifySubscriber(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	var subscribers []models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		AccountID: user.ID,
	}).Find(&subscribers).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(subscribers)
}

func addNotifySubscriber(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	var data struct {
		Provider    string `json:"provider" validate:"required"`
		DeviceToken string `json:"device_token" validate:"required"`
		DeviceID    string `json:"device_id" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var count int64
	if err := database.C.Where(&models.NotificationSubscriber{
		DeviceID:    data.DeviceID,
		DeviceToken: data.DeviceToken,
		AccountID:   user.ID,
	}).Model(&models.NotificationSubscriber{}).Count(&count).Error; err != nil || count > 0 {
		return c.SendStatus(fiber.StatusOK)
	}

	var account models.Account
	if err := database.C.Where("id = ?", user.ID).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err))
	}

	subscriber, err := services.AddNotifySubscriber(
		account,
		data.Provider,
		data.DeviceID,
		data.DeviceToken,
		c.Get(fiber.HeaderUserAgent),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	services.AddEvent(user.ID, "notifications.subscribe.push", data.DeviceID, c.IP(), c.Get(fiber.HeaderUserAgent))
	return c.JSON(subscriber)
}

func removeNotifySubscriber(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(*sec.UserInfo)

	device := c.Params("deviceId")

	if err := database.C.Where(&models.NotificationSubscriber{
		DeviceID:  device,
		AccountID: user.ID,
	}).Delete(&models.NotificationSubscriber{}).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	services.AddEvent(user.ID, "notifications.unsubscribe.push", device, c.IP(), c.Get(fiber.HeaderUserAgent))
	return c.SendStatus(fiber.StatusOK)
}
