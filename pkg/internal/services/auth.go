package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	jsoniter "github.com/json-iterator/go"

	localCache "git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Authenticate(sessionId uint) (ctx models.AuthTicket, perms map[string]any, err error) {
	if ctx, err = GetAuthContext(sessionId); err == nil {
		var heldPerms map[string]any
		rawHeldPerms, _ := jsoniter.Marshal(ctx.Account.PermNodes)
		_ = jsoniter.Unmarshal(rawHeldPerms, &heldPerms)

		perms = FilterPermNodes(heldPerms, ctx.Claims)
		ctx.Account.PermNodes = perms
		return
	}

	err = fiber.NewError(fiber.StatusUnauthorized, err.Error())
	return
}

func GetAuthContextCacheKey(sessionId uint) string {
	return fmt.Sprintf("auth-context#%d", sessionId)
}

func GetAuthContext(sessionId uint) (models.AuthTicket, error) {
	var err error
	var ctx models.AuthTicket

	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)

	key := GetAuthContextCacheKey(sessionId)
	if val, err := marshal.Get(context.Background(), key, new(models.AuthTicket)); err == nil {
		ctx = *val.(*models.AuthTicket)
	} else {
		ctx, err = CacheAuthContext(sessionId)
		log.Debug().Uint("session", sessionId).Msg("Created a new auth context cache")
	}

	return ctx, err
}

func CacheAuthContext(sessionId uint) (models.AuthTicket, error) {
	// Query data from primary database
	var ticket models.AuthTicket
	if err := database.C.
		Where("id = ?", sessionId).
		Preload("Account").
		First(&ticket).Error; err != nil {
		return ticket, fmt.Errorf("invalid auth ticket: %v", err)
	} else if err := ticket.IsAvailable(); err != nil {
		return ticket, fmt.Errorf("unavailable auth ticket: %v", err)
	}

	user, err := GetAccount(ticket.AccountID)
	if err != nil {
		return ticket, fmt.Errorf("invalid account: %v", err)
	}
	groups, err := GetUserAccountGroup(user)
	if err != nil {
		return ticket, fmt.Errorf("unable to get account groups: %v", err)
	}

	for _, group := range groups {
		for k, v := range group.PermNodes {
			if _, ok := user.PermNodes[k]; !ok {
				user.PermNodes[k] = v
			}
		}
	}
	ticket.Account = user

	// Put the data into the cache
	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)

	key := GetAuthContextCacheKey(sessionId)
	err = marshal.Set(
		context.Background(),
		key,
		ticket,
		store.WithExpiration(10*time.Minute),
		store.WithTags([]string{"auth-context", fmt.Sprintf("user#%d", user.ID)}),
	)

	return ticket, nil
}

func InvalidAuthCacheWithUser(userId uint) {
	cacheManager := cache.New[any](localCache.S)
	ctx := context.Background()

	cacheManager.Invalidate(
		ctx,
		store.WithInvalidateTags([]string{"auth-context", fmt.Sprintf("user#%d", userId)}),
	)
}
