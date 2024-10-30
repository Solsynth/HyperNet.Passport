package models

import "time"

type StatusAttitude = uint8

const (
	AttitudeNeutral = StatusAttitude(iota)
	AttitudePositive
	AttitudeNegative
)

type Status struct {
	BaseModel

	Type        string         `json:"type"`
	Label       string         `json:"label"`
	Attitude    StatusAttitude `json:"attitude"`
	IsNoDisturb bool           `json:"is_no_disturb"`
	IsInvisible bool           `json:"is_invisible"`
	ClearAt     *time.Time     `json:"clear_at"`
	AccountID   uint           `json:"account_id"`
}
