package services

import (
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"

	jsoniter "github.com/json-iterator/go"

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

func KgAuthContextCache(sessionId uint) string {
	return cachekit.FKey("auth-context", sessionId)
}

func GetAuthContext(sessionId uint) (models.AuthTicket, error) {
	var err error
	var ctx models.AuthTicket

	key := KgAuthContextCache(sessionId)
	if val, err := cachekit.Get[models.AuthTicket](gap.Ca, key); err == nil {
		ctx = val
	} else {
		log.Error().Err(err).Msg("Unable to get auth context cache")
		ctx, err = CacheAuthContext(sessionId)
		if err != nil {
			log.Error().Err(err).Msg("Unable to cache auth context")
		} else {
			log.Debug().Uint("session", sessionId).Msg("Created a new auth context cache")
		}
	}

	return ctx, err
}

func CacheAuthContext(sessionId uint) (models.AuthTicket, error) {
	// Query data from primary database
	var ticket models.AuthTicket
	if err := database.C.
		Where("id = ?", sessionId).
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
	key := KgAuthContextCache(sessionId)
	err = cachekit.Set[models.AuthTicket](
		gap.Ca,
		key,
		ticket,
		time.Minute*10,
		"auth-context",
		fmt.Sprintf("user#%d", user.ID),
	)
	if err != nil {
		log.Error().Err(err).Msg("Unable to cache auth context...")
	}

	return ticket, err
}

func InvalidUserAuthCache(uid uint) {
	cachekit.DeleteByTags(gap.Ca, "auth-context", fmt.Sprintf("user#%d", uid))
}
