package models

import "gorm.io/datatypes"

type Realm struct {
	BaseModel

	Alias        string            `json:"alias" gorm:"uniqueIndex"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Members      []RealmMember     `json:"members"`
	Avatar       *string           `json:"avatar"`
	Banner       *string           `json:"banner"`
	AccessPolicy datatypes.JSONMap `json:"access_policy"`
	IsPublic     bool              `json:"is_public"`
	IsCommunity  bool              `json:"is_community"`
	AccountID    uint              `json:"account_id"`
}

type RealmMember struct {
	BaseModel

	RealmID    uint    `json:"realm_id"`
	AccountID  uint    `json:"account_id"`
	Realm      Realm   `json:"realm"`
	Account    Account `json:"account"`
	PowerLevel int     `json:"power_level"`
}
