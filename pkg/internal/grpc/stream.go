package grpc

import (
	"context"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"github.com/rs/zerolog/log"
)

func (v *App) BroadcastEvent(ctx context.Context, request *proto.EventInfo) (*proto.EventResponse, error) {
	log.Debug().Str("event", request.GetEvent()).
		Msg("Got a broadcasting event...")
	switch request.GetEvent() {
	case "ws.client.register":
		// No longer need update user online status
		// Based on realtime sever connection status
		break
	case "ws.client.unregister":
		// Update user last seen at
		data := nex.DecodeMap(request.GetData())
		err := services.SetAccountLastSeen(uint(data["user"].(float64)))
		log.Debug().Err(err).Any("event", data).Msg("Setting account last seen...")
	}

	return &proto.EventResponse{}, nil
}
