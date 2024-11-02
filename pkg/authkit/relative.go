package authkit

import (
	"context"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
)

func ListRelative(nx *nex.Conn, userId uint, status int32, isRelated bool) ([]*proto.UserInfo, error) {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return nil, err
	}
	resp, err := proto.NewUserServiceClient(conn).ListUserRelative(context.Background(), &proto.ListUserRelativeRequest{
		UserId:    uint64(userId),
		Status:    status,
		IsRelated: isRelated,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetData(), err
}
