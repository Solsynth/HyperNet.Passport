package models

import "gorm.io/datatypes"

type Badge struct {
	BaseModel

	Type      string            `json:"type"`
	Metadata  datatypes.JSONMap `json:"metadata"`
	IsActive  bool              `json:"is_active"`
	AccountID uint              `json:"account_id"`
}
