package grpc

import (
	"context"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
)

func (v *App) RecordEvent(ctx context.Context, request *proto.RecordEventRequest) (*proto.RecordEventResponse, error) {
	services.AddEvent(
		uint(request.GetUserId()),
		request.GetAction(),
		nex.DecodeMap(request.GetMetadata()),
		request.GetIp(),
		request.GetUserAgent(),
	)

	return &proto.RecordEventResponse{IsSuccess: true}, nil
}
