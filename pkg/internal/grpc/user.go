package grpc

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	localCache "git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (v *App) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserInfo, error) {
	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	contx := context.Background()

	var account models.Account

	tx := database.C
	hitCache := false
	if request.UserId != nil {
		if val, err := marshal.Get(contx, services.GetAccountCacheKey(request.GetUserId()), new(models.Account)); err == nil {
			account = *val.(*models.Account)
			hitCache = true
		} else {
			tx = tx.Where("id = ?", uint(request.GetUserId()))
		}
	}
	if request.Name != nil {
		if val, err := marshal.Get(contx, services.GetAccountCacheKey(request.GetName()), new(models.Account)); err == nil {
			account = *val.(*models.Account)
			hitCache = true
		} else {
			tx = tx.Where("name = ?", request.GetName())
		}
	}

	if !hitCache {
		if err := tx.
			Preload("Profile").
			Preload("Badges", func(db *gorm.DB) *gorm.DB {
				return db.Order("badges.type DESC")
			}).
			First(&account).Error; err != nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("requested user with id %d was not found", request.GetUserId()))
		}

		groups, err := services.GetUserAccountGroup(account)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("unable to get account groups: %v", err))
		}
		for _, group := range groups {
			for k, v := range group.PermNodes {
				if _, ok := account.PermNodes[k]; !ok {
					account.PermNodes[k] = v
				}
			}
		}

		services.CacheAccount(account)
	}

	return account.EncodeToUserInfo(), nil
}

func (v *App) ListUser(ctx context.Context, request *proto.ListUserRequest) (*proto.MultipleUserInfo, error) {
	var accounts []models.Account
	if err := database.C.
		Where("id IN ?", lo.Map(request.GetUserId(), func(id uint64, _ int) interface{} { return id })).
		Find(&accounts).Error; err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to list users: %v", err))
	}
	return &proto.MultipleUserInfo{
		Data: lo.Map(request.GetUserId(), func(item uint64, index int) *proto.UserInfo {
			val, ok := lo.Find(accounts, func(x models.Account) bool {
				return uint(item) == x.ID
			})
			if !ok {
				return nil
			}
			return val.EncodeToUserInfo()
		}),
	}, nil
}

func (v *App) ListUserRelative(ctx context.Context, request *proto.ListUserRelativeRequest) (*proto.ListUserRelativeResponse, error) {
	tx := database.C.Preload("Account").Where("status = ?", request.GetStatus())

	if request.GetIsRelated() {
		tx = tx.Where("related_id = ?", request.GetUserId())
	} else {
		tx = tx.Where("account_id = ?", request.GetUserId())
	}

	var data []models.AccountRelationship
	if err := tx.Find(&data).Error; err != nil {
		return nil, err
	}

	return &proto.ListUserRelativeResponse{
		Data: lo.Map(data, func(item models.AccountRelationship, index int) *proto.UserInfo {
			val := &proto.UserInfo{
				Id:   uint64(item.AccountID),
				Name: item.Account.Name,
			}

			return val
		}),
	}, nil
}
