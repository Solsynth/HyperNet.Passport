package models

import (
	"time"

	"gorm.io/datatypes"
)

type AccountProfile struct {
	BaseModel

	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	Description string            `json:"description"`
	TimeZone    string            `json:"time_zone"`
	Location    string            `json:"location"`
	Pronouns    string            `json:"pronouns"`
	Gender      string            `json:"gender"`
	Links       datatypes.JSONMap `json:"links"`
	Experience  uint64            `json:"experience"`
	LastSeenAt  *time.Time        `json:"last_seen_at"`
	Birthday    *time.Time        `json:"birthday"`
	AccountID   uint              `json:"account_id"`
}

type AccountPage struct {
	BaseModel

	Content   string `json:"content"`
	AccountID uint   `json:"account_id"`
}
