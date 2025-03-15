package models

import "gorm.io/datatypes"

type ActionEvent struct {
	BaseModel

	Type        string            `json:"type"`
	Metadata    datatypes.JSONMap `json:"metadata"`
	Location    *string           `json:"location"`
	CoordinateX *float64          `json:"coordinate_x"`
	CoordinateY *float64          `json:"coordinate_y"`
	IpAddress   string            `json:"ip_address"`
	UserAgent   string            `json:"user_agent"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}
