package models

import (
	"gorm.io/datatypes"
	"time"
)

type Notification struct {
	BaseModel

	Topic    string            `json:"topic"`
	Title    string            `json:"title"`
	Subtitle string            `json:"subtitle"`
	Body     string            `json:"body"`
	Metadata datatypes.JSONMap `json:"metadata"`
	Priority int               `json:"priority"`
	SenderID *uint             `json:"sender_id"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`

	ReadAt *time.Time `json:"read_at"`

	IsRealtime bool `json:"is_realtime" gorm:"-"`
}

const (
	NotifySubscriberFirebase = "firebase"
	NotifySubscriberAPNs     = "apple"
)

type NotificationSubscriber struct {
	BaseModel

	UserAgent   string `json:"user_agent"`
	Provider    string `json:"provider"`
	DeviceID    string `json:"device_id" gorm:"uniqueIndex"`
	DeviceToken string `json:"device_token"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}
