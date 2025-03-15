package models

import "gorm.io/datatypes"

type AuditRecord struct {
	BaseModel

	Action      string            `json:"action"`
	Metadata    datatypes.JSONMap `json:"metadata"`
	Location    *string           `json:"location"`
	CoordinateX *float64          `json:"coordinate_x"`
	CoordinateY *float64          `json:"coordinate_y"`
	UserAgent   string            `json:"user_agent"`
	IpAddress   string            `json:"ip_address"`
	AccountID   uint              `json:"account_id"`
}
