package models

import (
	"time"

	"gorm.io/datatypes"
)

const (
	PunishmentTypeStrike = iota
	PunishmentTypeLimited
	PunishmentTypeDisabled
)

type Punishment struct {
	BaseModel

	Reason      string            `json:"reason"`
	Type        int               `json:"type"`
	PermNodes   datatypes.JSONMap `json:"perm_nodes"`
	ExpiredAt   *time.Time        `json:"expired_at"`
	Account     Account           `json:"account"`
	AccountID   uint              `json:"account_id"`
	Moderator   *Account          `json:"moderator"`
	ModeratorID *uint             `json:"moderator_id"`
}
