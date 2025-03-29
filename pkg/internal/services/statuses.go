package services

import (
	"context"
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"

	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/samber/lo"
)

func KgStatusCache(uid uint) string {
	return fmt.Sprintf("user-status#%d", uid)
}

func GetStatus(uid uint) (models.Status, error) {
	if val, err := cachekit.Get[models.Status](gap.Ca, KgStatusCache(uid)); err == nil {
		return val, nil
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
	cachekit.Set[models.Status](
		gap.Ca,
		KgStatusCache(uid),
		status,
		time.Minute*5,
		fmt.Sprintf("user#%d", uid),
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
		cachekit.Delete(gap.Ca, KgStatusCache(user.ID))
	}

	return nil
}
