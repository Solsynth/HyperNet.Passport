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
	var client models.ThirdClient
	if err := database.C.Where("id = ?", request.ClientId).First(&client).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "requested client with id %d was not found", request.ClientId)
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
