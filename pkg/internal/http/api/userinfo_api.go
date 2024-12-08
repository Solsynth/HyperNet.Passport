package api

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	localCache "git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"gorm.io/gorm"
	"strconv"
	"strings"

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
