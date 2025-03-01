package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	localCache "git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"gorm.io/gorm"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getOtherUserinfo(c *fiber.Ctx) error {
	alias := c.Params("alias")

	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	ctx := context.Background()

	if val, err := marshal.Get(ctx, services.GetAccountCacheKey(alias), new(models.Account)); err == nil {
		return c.JSON(*val.(*models.Account))
	}

	tx := database.C.Where("name = ?", alias)

	numericId, err := strconv.Atoi(alias)
	if err == nil {
		if val, err := marshal.Get(ctx, services.GetAccountCacheKey(numericId), new(models.Account)); err == nil {
			return c.JSON(*val.(*models.Account))
		}
		tx = tx.Or("id = ?", numericId)
	}

	var account models.Account
	if err := tx.
		Preload("Profile").
		Preload("Badges", func(db *gorm.DB) *gorm.DB {
			return db.Order("badges.type DESC")
		}).
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	groups, err := services.GetUserAccountGroup(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("unable to get account groups: %v", err))
	}
	for _, group := range groups {
		for k, v := range group.PermNodes {
			if _, ok := account.PermNodes[k]; !ok {
				account.PermNodes[k] = v
			}
		}
	}

	services.CacheAccount(account)

	return c.JSON(account)
}

func getOtherUserinfoBatch(c *fiber.Ctx) error {
	idFilter := c.Query("id")
	nameFilter := c.Query("name")
	idSet := strings.Split(idFilter, ",")
	nameSet := strings.Split(nameFilter, ",")
	if len(idSet) == 0 && len(nameSet) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "query filter is required")
	}

	if len(idSet)+len(nameSet) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "only support 100 users in a single batch")
	}

	tx := database.C.Model(&models.Account{}).Limit(100)
	if len(idFilter) > 0 {
		tx = tx.Where("id IN ?", idSet)
	}
	if len(nameFilter) > 0 {
		tx = tx.Where("name IN ?", nameSet)
	}

	var accounts []models.Account
	if err := tx.Find(&accounts).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(accounts)
}

func getUserinfoForOidc(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		Preload("Contacts").
		Preload("Badges").
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		data.PermNodes = c.Locals("nex_user").(*sec.UserInfo).PermNodes
	}

	return c.JSON(fiber.Map{
		"sub":                fmt.Sprintf("%d", data.ID),
		"family_name":        data.Profile.FirstName,
		"given_name":         data.Profile.LastName,
		"name":               data.Name,
		"email":              data.GetPrimaryEmail().Content,
		"email_verified":     data.GetPrimaryEmail().VerifiedAt != nil,
		"preferred_username": data.Nick,
		"picture":            data.GetAvatar(),
		"birthdate":          data.Profile.Birthday,
		"updated_at":         data.UpdatedAt,
	})
}
