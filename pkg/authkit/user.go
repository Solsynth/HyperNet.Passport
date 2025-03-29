package authkit

import (
	"context"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"github.com/samber/lo"
)

func GetUser(nx *nex.Conn, userId uint) (models.Account, error) {
	cacheConn, err := cachekit.NewConn(nx, 3*time.Second)
	if err == nil {
		key := cachekit.FKey(cachekit.DAAttachment, userId)
		if user, err := cachekit.Get[models.Account](cacheConn, key); err == nil {
			return user, nil
		}
	}

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
	cacheConn, err := cachekit.NewConn(nx, 3*time.Second)
	if err == nil {
		key := cachekit.FKey(cachekit.DAAttachment, name)
		if user, err := cachekit.Get[models.Account](cacheConn, key); err == nil {
			return user, nil
		}
	}

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

func ListUser(nx *nex.Conn, userIds []uint) ([]models.Account, error) {
	var accounts []models.Account
	var missingId []uint
	cachedUsers := make(map[uint]models.Account)

	// Try to get users from cache
	cacheConn, err := cachekit.NewConn(nx, 3*time.Second)
	if err == nil {
		for _, userId := range userIds {
			key := cachekit.FKey(cachekit.DAAttachment, userId)
			if user, err := cachekit.Get[models.Account](cacheConn, key); err == nil {
				cachedUsers[userId] = user
			} else {
				missingId = append(missingId, userId)
			}
		}
	}

	// If all users are found in cache, return them
	if len(missingId) == 0 {
		for _, account := range cachedUsers {
			accounts = append(accounts, account)
		}
		return accounts, nil
	}

	// Fetch missing users from the gRPC service
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return nil, err
	}

	raw, _ := proto.NewUserServiceClient(conn).ListUser(context.Background(), &proto.ListUserRequest{
		UserId: lo.Map(missingId, func(item uint, index int) uint64 {
			return uint64(item)
		}),
	})

	// Convert fetched users and add to the result
	for _, item := range raw.GetData() {
		account := GetAccountFromUserInfo(&sec.UserInfo{
			ID:        uint(item.GetId()),
			Name:      item.GetName(),
			PermNodes: nex.DecodeMap(item.GetPermNodes()),
			Metadata:  nex.DecodeMap(item.GetMetadata()),
		})
		accounts = append(accounts, account)
	}

	// Merge cached and fetched results
	for _, account := range cachedUsers {
		accounts = append(accounts, account)
	}

	return accounts, nil
}
