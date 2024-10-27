package grpc

import (
	"net"

	"google.golang.org/grpc/reflection"

	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	nroto "git.solsynth.dev/hypernet/nexus/pkg/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	nroto.UnimplementedAuthServiceServer
	nroto.UnimplementedDirectoryServiceServer
	proto.UnimplementedNotifierServer
	proto.UnimplementedRealmServer
	proto.UnimplementedEventRecorderServer
	health.UnimplementedHealthServer

	srv *grpc.Server
}

func NewServer() *App {
	server := &App{
		srv: grpc.NewServer(),
	}

	nroto.RegisterAuthServiceServer(server.srv, server)
	proto.RegisterNotifierServer(server.srv, server)
	proto.RegisterRealmServer(server.srv, server)
	proto.RegisterEventRecorderServer(server.srv, server)
	health.RegisterHealthServer(server.srv, server)

	reflection.Register(server.srv)

	return server
}

func (v *App) Listen() error {
	listener, err := net.Listen("tcp", viper.GetString("grpc_bind"))
	if err != nil {
		return err
	}

	return v.srv.Serve(listener)
}
