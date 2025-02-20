package models

import (
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"gorm.io/datatypes"
)

type Realm struct {
	BaseModel

	Alias        string            `json:"alias" gorm:"uniqueIndex"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Members      []RealmMember     `json:"members"`
	Avatar       *string           `json:"avatar"`
	Banner       *string           `json:"banner"`
	Popularity   int               `json:"popularity"`
	AccessPolicy datatypes.JSONMap `json:"access_policy"`
	IsPublic     bool              `json:"is_public"`
	IsCommunity  bool              `json:"is_community"`
	AccountID    uint              `json:"account_id"`
}

func NewRealmFromProto(proto *proto.RealmInfo) Realm {
	return Realm{
		BaseModel: BaseModel{
			ID: uint(proto.GetId()),
		},
		Alias:        proto.GetAlias(),
		Name:         proto.GetName(),
		Description:  proto.GetDescription(),
		Avatar:       &proto.Avatar,
		Banner:       &proto.Banner,
		IsPublic:     proto.GetIsPublic(),
		IsCommunity:  proto.GetIsCommunity(),
		AccessPolicy: nex.DecodeMap(proto.GetAccessPolicy()),
	}
}

type RealmMember struct {
	BaseModel

	RealmID    uint    `json:"realm_id"`
	AccountID  uint    `json:"account_id"`
	Realm      Realm   `json:"realm"`
	Account    Account `json:"account"`
	PowerLevel int     `json:"power_level"`
}

func NewRealmMemberFromProto(proto *proto.RealmMemberInfo) RealmMember {
	return RealmMember{
		BaseModel: BaseModel{
			ID: uint(proto.GetId()),
		},
		RealmID:    uint(proto.GetRealmId()),
		AccountID:  uint(proto.GetUserId()),
		PowerLevel: int(proto.GetPowerLevel()),
	}
}
