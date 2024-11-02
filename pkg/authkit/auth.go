package authkit

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"github.com/samber/lo"
)

func EnsureUserPermGranted(nx *nex.Conn, userId, otherId uint, key string, val any) error {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return fmt.Errorf("failed to get auth service client: %v", err)
	}
	resp, err := proto.NewAuthServiceClient(conn).EnsureUserPermGranted(context.Background(), &proto.CheckUserPermRequest{
		UserId:  uint64(userId),
		OtherId: uint64(otherId),
		Key:     key,
		Value:   nex.EncodeMap(val),
	})
	if err != nil {
		return err
	}
	return lo.Ternary(resp.GetIsValid(), nil, fmt.Errorf("missing permission: %v", key))
}
