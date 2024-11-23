package grpc

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v *App) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserInfo, error) {
	tx := database.C
	if request.UserId != nil {
		tx = tx.Where("id = ?", uint(request.GetUserId()))
	}
	if request.Name != nil {
		tx = tx.Where("name = ?", request.GetName())
	}

	var account models.Account
	if err := tx.First(&account).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("requested user with id %d was not found", request.GetUserId()))
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
