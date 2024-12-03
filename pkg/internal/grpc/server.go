package grpc

import (
	"net"

	"google.golang.org/grpc/reflection"

	nroto "git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	nroto.UnimplementedAuthServiceServer
	nroto.UnimplementedDirectoryServiceServer
	nroto.UnimplementedUserServiceServer
	proto.UnimplementedRealmServiceServer
	proto.UnimplementedAuditServiceServer
	proto.UnimplementedNotifyServiceServer
	health.UnimplementedHealthServer

	srv *grpc.Server
}

func NewServer() *App {
	server := &App{
		srv: grpc.NewServer(),
	}

	nroto.RegisterAuthServiceServer(server.srv, server)
	nroto.RegisterUserServiceServer(server.srv, server)
	nroto.RegisterDirectoryServiceServer(server.srv, server)
	proto.RegisterNotifyServiceServer(server.srv, server)
	proto.RegisterRealmServiceServer(server.srv, server)
	proto.RegisterAuditServiceServer(server.srv, server)
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
