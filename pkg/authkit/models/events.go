package models

import "gorm.io/datatypes"

type ActionEvent struct {
	BaseModel

	Type      string            `json:"type"`
	Metadata  datatypes.JSONMap `json:"metadata"`
	Location  string            `json:"location"`
	IpAddress string            `json:"ip_address"`
	UserAgent string            `json:"user_agent"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}
