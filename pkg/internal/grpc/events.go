package grpc

import (
	"context"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
)

func (v *App) RecordEvent(ctx context.Context, request *proto.RecordEventRequest) (*proto.RecordEventResponse, error) {
	services.AddEvent(
		uint(request.GetUserId()),
		request.GetAction(),
		request.GetTarget(),
		request.GetIp(),
		request.GetUserAgent(),
	)

	return &proto.RecordEventResponse{IsSuccess: true}, nil
}
