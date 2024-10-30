package models

import "gorm.io/datatypes"

type PreferenceAuth struct {
	BaseModel

	Config    datatypes.JSONType[AuthConfig] `json:"config"`
	AccountID uint                           `json:"account_id"`
	Account   Account                        `json:"account"`
}

type PreferenceNotification struct {
	BaseModel

	Config    datatypes.JSONMap `json:"config"`
	AccountID uint              `json:"account_id"`
	Account   Account           `json:"account"`
}
