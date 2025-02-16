package models

import (
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"gorm.io/datatypes"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Account struct {
	BaseModel

	Name        string            `json:"name" gorm:"uniqueIndex"`
	Nick        string            `json:"nick"`
	Description string            `json:"description"`
	Avatar      *string           `json:"avatar"`
	Banner      *string           `json:"banner"`
	ConfirmedAt *time.Time        `json:"confirmed_at"`
	SuspendedAt *time.Time        `json:"suspended_at"`
	PermNodes   datatypes.JSONMap `json:"perm_nodes"`
	Language    string            `json:"language"`

	AutomatedBy *Account `json:"automated_by" gorm:"foreignKey:AutomatedID"`
	AutomatedID *uint    `json:"automated_id"`

	AffiliatedTo *Realm `json:"affiliated_to" gorm:"foreignKey:AffiliatedID"`
	AffiliatedID *uint  `json:"affiliated_id"`

	Profile  AccountProfile   `json:"profile,omitempty"`
	Contacts []AccountContact `json:"contacts,omitempty"`
	Badges   []Badge          `json:"badges,omitempty"`

	Tickets []AuthTicket `json:"tickets,omitempty"`
	Factors []AuthFactor `json:"factors,omitempty"`

	Relations []AccountRelationship `json:"relations,omitempty" gorm:"foreignKey:AccountID"`
}

func (v Account) GetAvatar() *string {
	if v.Avatar != nil {
		return lo.ToPtr(fmt.Sprintf("%s/%s", viper.GetString("content_endpoint"), *v.Avatar))
	}
	return nil
}

func (v Account) GetBanner() *string {
	if v.Banner != nil {
		return lo.ToPtr(fmt.Sprintf("%s/%s", viper.GetString("content_endpoint"), *v.Banner))
	}
	return nil
}

func (v Account) GetPrimaryEmail() AccountContact {
	val, _ := lo.Find(v.Contacts, func(item AccountContact) bool {
		return item.Type == EmailAccountContact && item.IsPrimary
	})
	return val
}

func (v Account) EncodeToUserInfo() *proto.UserInfo {
	return &proto.UserInfo{
		Id:        uint64(v.ID),
		Name:      v.Name,
		Language:  v.Language,
		PermNodes: nex.EncodeMap(v.PermNodes),
		Metadata:  nex.EncodeMap(v),
	}
}

type AccountContactType = int8

const (
	EmailAccountContact = AccountContactType(iota)
)

type AccountContact struct {
	BaseModel

	Type       int8       `json:"type"`
	Content    string     `json:"content" gorm:"uniqueIndex"`
	IsPublic   bool       `json:"is_public"`
	IsPrimary  bool       `json:"is_primary"`
	VerifiedAt *time.Time `json:"verified_at"`
	AccountID  uint       `json:"account_id"`
}
