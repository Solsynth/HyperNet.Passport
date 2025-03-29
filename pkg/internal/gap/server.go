package gap

import (
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/localize"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/rx"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit/pushcon"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/spf13/viper"
)

var (
	Nx *nex.Conn
	Px *pushcon.Conn
	Rx *rx.MqConn
	Ca *cachekit.Conn
)

const (
	FactorOtpPrefix = "auth-otp"
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
	Ca, err = cachekit.NewConn(Nx, time.Second*3)
	if err != nil {
		return fmt.Errorf("error during initialize nexus cache module: %v", err)
	}

	return err
}

func LoadLocalization() error {
	return localize.LoadLocalization(viper.GetString("locales_dir"), viper.GetString("templates_dir"))
}
