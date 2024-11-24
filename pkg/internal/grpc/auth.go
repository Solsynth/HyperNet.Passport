package grpc

import (
	"context"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hypernet/nexus/pkg/proto"
)

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
