package gap

import (
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit/pushcon"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"strings"

	"github.com/spf13/viper"
)

var Nx *nex.Conn
var Px *pushcon.Conn

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

	return err
}
