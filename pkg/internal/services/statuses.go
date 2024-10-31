package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"time"

	localCache "git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/samber/lo"
)

func GetStatusCacheKey(uid uint) string {
	return fmt.Sprintf("user-status#%d", uid)
}

func GetStatus(uid uint) (models.Status, error) {
	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	contx := context.Background()

	if val, err := marshal.Get(contx, GetStatusCacheKey(uid), new(models.Status)); err == nil {
		status := val.(models.Status)
		if status.ClearAt != nil && status.ClearAt.Before(time.Now()) {
			marshal.Delete(contx, GetStatusCacheKey(uid))
		} else {
			return status, nil
		}
	}
	var status models.Status
	if err := database.C.
		Where("account_id = ?", uid).
		Where("clear_at > ?", time.Now()).
		First(&status).Error; err != nil {
		return status, err
	} else {
		CacheUserStatus(uid, status)
	}
	return status, nil
}

func CacheUserStatus(uid uint, status models.Status) {
	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	contx := context.Background()

	marshal.Set(
		contx,
		GetStatusCacheKey(uid),
		status,
		store.WithTags([]string{"user-status", fmt.Sprintf("user#%d", uid)}),
	)
}

func GetUserOnline(uid uint) bool {
	pc := proto.NewStreamServiceClient(gap.Nx.GetNexusGrpcConn())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := pc.CountStreamConnection(ctx, &proto.CountConnectionRequest{
		UserId: uint64(uid),
	})
	if err != nil {
		return false
	}
	return resp.Count > 0
}

func GetStatusDisturbable(uid uint) error {
	status, err := GetStatus(uid)
	isOnline := GetUserOnline(uid)
	if isOnline && err != nil {
		return nil
	} else if err == nil && status.IsNoDisturb {
		return fmt.Errorf("do not disturb")
	} else {
		return nil
	}
}

func GetStatusOnline(uid uint) error {
	status, err := GetStatus(uid)
	isOnline := GetUserOnline(uid)
	if isOnline && err != nil {
		return nil
	} else if err == nil && status.IsInvisible {
		return fmt.Errorf("invisible")
	} else if !isOnline {
		return fmt.Errorf("offline")
	} else {
		return nil
	}
}

func NewStatus(user models.Account, status models.Status) (models.Status, error) {
	if err := database.C.Save(&status).Error; err != nil {
		return status, err
	} else {
		CacheUserStatus(user.ID, status)
	}
	return status, nil
}

func EditStatus(user models.Account, status models.Status) (models.Status, error) {
	if err := database.C.Save(&status).Error; err != nil {
		return status, err
	} else {
		CacheUserStatus(user.ID, status)
	}
	return status, nil
}

func ClearStatus(user models.Account) error {
	if err := database.C.
		Where("account_id = ?", user.ID).
		Where("clear_at > ?", time.Now()).
		Updates(models.Status{ClearAt: lo.ToPtr(time.Now())}).Error; err != nil {
		return err
	} else {
		cacheManager := cache.New[any](localCache.S)
		marshal := marshaler.New(cacheManager)
		contx := context.Background()

		marshal.Delete(contx, GetStatusCacheKey(user.ID))
	}

	return nil
}
