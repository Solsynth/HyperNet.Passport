package grpc

import (
	"context"
	"git.solsynth.dev/hydrogen/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"github.com/samber/lo"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hypernet/nexus/pkg/proto"
)

type authenticateServer struct {
	proto.UnimplementedAuthServiceServer
}

func (v *App) Authenticate(_ context.Context, in *proto.AuthRequest) (*proto.AuthReply, error) {
	ticket, perms, err := services.Authenticate(uint(in.GetSessionId()))
	if err != nil {
		return &proto.AuthReply{
			IsValid: false,
		}, nil
	} else {
		user := ticket.Account
		userinfo := &proto.UserInfo{
			Id:        uint64(user.ID),
			Name:      user.Name,
			PermNodes: nex.EncodeMap(perms),
			Metadata:  nex.EncodeMap(user),
		}

		return &proto.AuthReply{
			IsValid: true,
			Info: &proto.AuthInfo{
				SessionId: uint64(ticket.ID),
				Info:      userinfo,
			},
		}, nil
	}
}

func (v *App) EnsurePermGranted(_ context.Context, in *proto.CheckPermRequest) (*proto.CheckPermResponse, error) {
	ctx, err := services.GetAuthContext(uint(in.GetSessionId()))
	if err != nil {
		return nil, err
	}

	var heldPerms map[string]any
	rawHeldPerms, _ := jsoniter.Marshal(ctx.Account.PermNodes)
	_ = jsoniter.Unmarshal(rawHeldPerms, &heldPerms)

	var value any
	_ = jsoniter.Unmarshal(in.GetValue(), &value)
	perms := services.FilterPermNodes(heldPerms, ctx.Claims)
	valid := services.HasPermNode(perms, in.GetKey(), value)

	return &proto.CheckPermResponse{
		IsValid: valid,
	}, nil
}

func (v *App) EnsureUserPermGranted(_ context.Context, in *proto.CheckUserPermRequest) (*proto.CheckUserPermResponse, error) {
	relation, err := services.GetRelationWithTwoNode(uint(in.GetUserId()), uint(in.GetOtherId()))
	if err != nil {
		return &proto.CheckUserPermResponse{
			IsValid: false,
		}, nil
	}

	defaultPerm := relation.Status == models.RelationshipFriend

	var value any
	_ = jsoniter.Unmarshal(in.GetValue(), &value)
	valid := services.HasPermNodeWithDefault(relation.PermNodes, in.GetKey(), value, defaultPerm)

	return &proto.CheckUserPermResponse{
		IsValid: valid,
	}, nil
}

func (v *App) ListUserFriends(_ context.Context, in *proto.ListUserRelativeRequest) (*proto.ListUserRelativeResponse, error) {
	tx := database.C.Preload("Account").Where("status = ?", models.RelationshipFriend)

	if in.GetIsRelated() {
		tx = tx.Where("related_id = ?", in.GetUserId())
	} else {
		tx = tx.Where("account_id = ?", in.GetUserId())
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

func (v *App) ListUserBlocklist(_ context.Context, in *proto.ListUserRelativeRequest) (*proto.ListUserRelativeResponse, error) {
	tx := database.C.Preload("Account").Where("status = ?", models.RelationshipBlocked)

	if in.GetIsRelated() {
		tx = tx.Where("related_id = ?", in.GetUserId())
	} else {
		tx = tx.Where("account_id = ?", in.GetUserId())
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
