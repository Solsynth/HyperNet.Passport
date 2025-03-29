package grpc

import (
	"context"
	"fmt"

	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v *App) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserInfo, error) {
	var account models.Account
	var err error
	if request.UserId != nil {
		account, err = services.GetAccountForEnd(uint(request.GetUserId()))
	} else if request.Name != nil {
		account, err = services.GetAccountForEnd(request.GetName())
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("unable to get account punishments: %v", err))
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
	tx := database.C.Preload("Account").Preload("Related").Where("status = ?", request.GetStatus())

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
			account := lo.Ternary(request.GetIsRelated(), item.Account, item.Related)
			val := &proto.UserInfo{
				Id:   uint64(account.ID),
				Name: account.Name,
			}

			return val
		}),
	}, nil
}
