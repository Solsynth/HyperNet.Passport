package grpc

import (
	"context"
	"fmt"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/passport/pkg/internal/services"
	"git.solsynth.dev/hypernet/passport/pkg/internal/web/exts"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func (v *App) BroadcastEvent(ctx context.Context, request *proto.EventInfo) (*proto.EventResponse, error) {
	log.Debug().Str("event", request.GetEvent()).
		Msg("Got a broadcasting event...")

	switch request.GetEvent() {
	// Last seen at
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

func (v *App) PushStream(_ context.Context, request *proto.PushStreamRequest) (*proto.PushStreamResponse, error) {
	sc := proto.NewStreamServiceClient(gap.Nx.GetNexusGrpcConn())

	var in nex.WebSocketPackage
	if err := jsoniter.Unmarshal(request.GetBody(), &in); err != nil {
		return nil, err
	}

	switch in.Action {
	// PaKex (Key Exchange)
	case "kex.ask":
		var data struct {
			UserID    uint   `json:"user_id" validate:"required"`
			KeypairID string `json:"keypair_id" validate:"required"`
			ClientID  string `json:"client_id" validate:"required"`
		}

		err := jsoniter.Unmarshal(in.RawPayload(), &data)
		if request.ClientId != nil {
			data.ClientID = *request.ClientId
		}
		if err == nil {
			err = exts.ValidateStruct(data)
		}
		if err != nil {
			_, _ = sc.PushStream(context.Background(), &proto.PushStreamRequest{
				ClientId: request.ClientId,
				Body: nex.WebSocketPackage{
					Action:  "error",
					Message: fmt.Sprintf("unable parse payload: %v", err),
				}.Marshal(),
			})
			break
		}

		// Forward ask request
		sc.PushStream(context.Background(), &proto.PushStreamRequest{
			UserId: lo.ToPtr(uint64(data.UserID)),
			Body: nex.WebSocketPackage{
				Action:  "kex.ask",
				Payload: data,
			}.Marshal(),
		})
	case "kex.ack":
		var data struct {
			UserID     uint   `json:"user_id" validate:"required"`
			KeypairID  string `json:"keypair_id" validate:"required"`
			PublicKey  string `json:"public_key"`
			PrivateKey string `json:"private_key"`
			ClientID   string `json:"client_id" validate:"required"`
		}

		err := jsoniter.Unmarshal(in.RawPayload(), &data)
		if err == nil {
			err = exts.ValidateStruct(data)
		}
		if err != nil {
			_, _ = sc.PushStream(context.Background(), &proto.PushStreamRequest{
				ClientId: request.ClientId,
				Body: nex.WebSocketPackage{
					Action:  "error",
					Message: fmt.Sprintf("unable parse payload: %v", err),
				}.Marshal(),
			})
			break
		}
		if len(data.PublicKey) == 0 && len(data.PrivateKey) == 0 {
			_, _ = sc.PushStream(context.Background(), &proto.PushStreamRequest{
				ClientId: request.ClientId,
				Body: nex.WebSocketPackage{
					Action:  "error",
					Message: "one of public key and private key is required",
				}.Marshal(),
			})
			break
		}

		// Forward ack request
		sc.PushStream(context.Background(), &proto.PushStreamRequest{
			ClientId: request.ClientId,
			Body: nex.WebSocketPackage{
				Action:  "kex.ack",
				Payload: data,
			}.Marshal(),
		})
	}

	return &proto.PushStreamResponse{}, nil
}
