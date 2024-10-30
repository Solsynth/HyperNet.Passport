package models

import "gorm.io/datatypes"

type Badge struct {
	BaseModel

	Type      string            `json:"type"`
	Metadata  datatypes.JSONMap `json:"metadata"`
	AccountID uint              `json:"account_id"`
}
