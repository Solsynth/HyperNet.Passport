package authkit

import (
	"context"
	"fmt"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"github.com/samber/lo"
)

func GetThirdClient(nx *nex.Conn, id uint, secret *string) (*models.ThirdClient, error) {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth service client: %v", err)
	}
	resp, err := proto.NewThirdClientServiceClient(conn).
		GetThirdClient(context.Background(), &proto.GetThirdClientRequest{
			Id:     lo.ToPtr(uint64(id)),
			Secret: secret,
		})
	if err != nil {
		return nil, err
	}

	return &models.ThirdClient{
		Alias:       resp.GetInfo().GetAlias(),
		Name:        resp.GetInfo().GetName(),
		Description: resp.GetInfo().GetDescription(),
		IsDraft:     resp.GetInfo().GetIsDraft(),
		AccountID: lo.TernaryF(resp.GetInfo().AccountId != nil, func() *uint {
			return lo.ToPtr(uint(resp.GetInfo().GetAccountId()))
		}, func() *uint {
			return nil
		}),
	}, nil
}

func GetThirdClientByAlias(nx *nex.Conn, alias string, secret *string) (*models.ThirdClient, error) {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth service client: %v", err)
	}
	resp, err := proto.NewThirdClientServiceClient(conn).
		GetThirdClient(context.Background(), &proto.GetThirdClientRequest{
			Alias:  &alias,
			Secret: secret,
		})
	if err != nil {
		return nil, err
	}

	return &models.ThirdClient{
		Alias:       resp.GetInfo().GetAlias(),
		Name:        resp.GetInfo().GetName(),
		Description: resp.GetInfo().GetDescription(),
		IsDraft:     resp.GetInfo().GetIsDraft(),
		AccountID: lo.TernaryF(resp.GetInfo().AccountId != nil, func() *uint {
			return lo.ToPtr(uint(resp.GetInfo().GetAccountId()))
		}, func() *uint {
			return nil
		}),
	}, nil
}
