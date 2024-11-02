package authkit

import (
	"context"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"github.com/samber/lo"
)

func GetUser(nx *nex.Conn, userId uint) (models.Account, error) {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return models.Account{}, err
	}
	raw, _ := proto.NewUserServiceClient(conn).GetUser(context.Background(), &proto.GetUserRequest{
		UserId: lo.ToPtr(uint64(userId)),
	})
	return GetAccountFromUserInfo(&sec.UserInfo{
		ID:        uint(raw.GetId()),
		Name:      raw.GetName(),
		PermNodes: nex.DecodeMap(raw.GetPermNodes()),
		Metadata:  nex.DecodeMap(raw.GetMetadata()),
	}), nil
}

func GetUserByName(nx *nex.Conn, name string) (models.Account, error) {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return models.Account{}, err
	}
	raw, _ := proto.NewUserServiceClient(conn).GetUser(context.Background(), &proto.GetUserRequest{
		Name: &name,
	})
	return GetAccountFromUserInfo(&sec.UserInfo{
		ID:        uint(raw.GetId()),
		Name:      raw.GetName(),
		PermNodes: nex.DecodeMap(raw.GetPermNodes()),
		Metadata:  nex.DecodeMap(raw.GetMetadata()),
	}), nil
}

func ListUser(nx *nex.Conn, userId []uint) ([]models.Account, error) {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return nil, err
	}
	raw, _ := proto.NewUserServiceClient(conn).ListUser(context.Background(), &proto.ListUserRequest{
		UserId: lo.Map(userId, func(item uint, index int) uint64 {
			return uint64(item)
		}),
	})
	var out []models.Account
	for _, item := range raw.GetData() {
		out = append(out, GetAccountFromUserInfo(&sec.UserInfo{
			ID:        uint(item.GetId()),
			Name:      item.GetName(),
			PermNodes: nex.DecodeMap(item.GetPermNodes()),
			Metadata:  nex.DecodeMap(item.GetMetadata()),
		}))
	}
	return out, nil
}
