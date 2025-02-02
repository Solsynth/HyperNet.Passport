package grpc

import (
	"context"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v *App) GetThirdClient(ctx context.Context, request *proto.GetThirdClientRequest) (*proto.GetThirdClientResponse, error) {
	tx := database.C
	if request.Id == nil && request.Alias == nil {
		return nil, status.Error(codes.InvalidArgument, "either id or alias must be specified")
	}
	if request.Id != nil {
		tx = tx.Where("id = ?", request.Id)
	} else if request.Alias != nil {
		tx = tx.Where("alias = ?", request.Alias)
	}

	var client models.ThirdClient
	if err := tx.First(&client).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "requested client was not found")
	}

	if request.Secret != nil {
		if client.Secret != request.GetSecret() {
			return nil, status.Errorf(codes.PermissionDenied, "invalid secret")
		}
	}

	return &proto.GetThirdClientResponse{
		Info: &proto.ThirdClientInfo{
			Id:          uint64(client.ID),
			Name:        client.Name,
			Description: client.Description,
		},
	}, nil
}
