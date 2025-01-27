package gap

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/rx"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit/pushcon"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/spf13/viper"
)

var (
	Nx *nex.Conn
	Px *pushcon.Conn
	Rx *rx.MqConn
	Jt nats.JetStreamContext
)

const (
	FactorOtpPrefix = "passport.otp."
)

func InitializeToNexus() error {
	grpcBind := strings.SplitN(viper.GetString("grpc_bind"), ":", 2)
	httpBind := strings.SplitN(viper.GetString("bind"), ":", 2)

	outboundIp, _ := nex.GetOutboundIP()

	grpcOutbound := fmt.Sprintf("%s:%s", outboundIp, grpcBind[1])
	httpOutbound := fmt.Sprintf("%s:%s", outboundIp, httpBind[1])

	var err error
	Nx, err = nex.NewNexusConn(viper.GetString("nexus_addr"), &proto.ServiceInfo{
		Id:       viper.GetString("id"),
		Type:     nex.ServiceTypeAuth,
		Label:    "Passport",
		GrpcAddr: grpcOutbound,
		HttpAddr: lo.ToPtr("http://" + httpOutbound + "/api"),
	})
	if err == nil {
		go func() {
			err := Nx.RunRegistering()
			if err != nil {
				log.Error().Err(err).Msg("An error occurred while registering service...")
			}
		}()
	}

	Px, err = pushcon.NewConn(Nx)
	if err != nil {
		return fmt.Errorf("error during initialize pushcon: %v", err)
	}

	Rx, err = rx.NewMqConn(Nx)
	if err != nil {
		return fmt.Errorf("error during initialize nexus rx module: %v", err)
	}
	Jt, err = Rx.Nt.JetStream()
	if err != nil {
		return fmt.Errorf("error during initialize nats jetstream: %v", err)
	}

	jetstreamCfg := &nats.StreamConfig{
		Name:     "Passport OTPs",
		Subjects: []string{FactorOtpPrefix + ">"},
		Storage:  nats.MemoryStorage,
		MaxAge:   5 * time.Minute,
	}
	_, err = Jt.AddStream(jetstreamCfg)
	if err != nil && !errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		return fmt.Errorf("error during initialize jetstream stream: %v", err)
	}

	return err
}
