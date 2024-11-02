package authkit

import (
	"context"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"github.com/samber/lo"
)

func GetRealm(nx *nex.Conn, id uint) (models.Realm, error) {
	var realm models.Realm
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return realm, err
	}
	resp, err := proto.NewRealmServiceClient(conn).GetRealm(context.Background(), &proto.LookupRealmRequest{
		Id: lo.ToPtr(uint64(id)),
	})
	if err != nil {
		return realm, err
	}
	return models.NewRealmFromProto(resp), nil
}

func GetRealmByAlias(nx *nex.Conn, alias string) (models.Realm, error) {
	var realm models.Realm
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return realm, err
	}
	resp, err := proto.NewRealmServiceClient(conn).GetRealm(context.Background(), &proto.LookupRealmRequest{
		Alias: &alias,
	})
	if err != nil {
		return realm, err
	}
	return models.NewRealmFromProto(resp), nil
}

func ListRealm(nx *nex.Conn, id []uint) ([]models.Realm, error) {
	var realms []models.Realm
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return realms, err
	}
	resp, err := proto.NewRealmServiceClient(conn).ListRealm(context.Background(), &proto.ListRealmRequest{
		Id: lo.Map(id, func(item uint, _ int) uint64 {
			return uint64(item)
		}),
	})
	if err != nil {
		return realms, err
	}
	for _, realm := range resp.GetData() {
		realms = append(realms, models.NewRealmFromProto(realm))
	}
	return realms, nil
}

func GetRealmMember(nx *nex.Conn, realmID, userID uint) (models.RealmMember, error) {
	var member models.RealmMember
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return member, err
	}
	resp, err := proto.NewRealmServiceClient(conn).GetRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: lo.ToPtr(uint64(realmID)),
		UserId:  lo.ToPtr(uint64(userID)),
	})
	if err != nil {
		return member, err
	}
	return models.NewRealmMemberFromProto(resp), nil
}

func ListRealmMember(nx *nex.Conn, realmID uint) ([]models.RealmMember, error) {
	var members []models.RealmMember
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return members, err
	}
	resp, err := proto.NewRealmServiceClient(conn).ListRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: lo.ToPtr(uint64(realmID)),
	})
	if err != nil {
		return members, err
	}
	for _, member := range resp.GetData() {
		members = append(members, models.NewRealmMemberFromProto(member))
	}
	return members, nil
}

func CheckRealmMemberPerm(nx *nex.Conn, realmID uint, userID, power int) bool {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return false
	}
	resp, err := proto.NewRealmServiceClient(conn).CheckRealmMemberPerm(context.Background(), &proto.CheckRealmPermRequest{
		RealmId:    uint64(realmID),
		UserId:     uint64(userID),
		PowerLevel: int32(power),
	})
	if err != nil {
		return false
	}
	return resp.GetIsSuccess()
}
