package models

import (
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"
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
}

func (v Notification) EncodeToPushkit() pushkit.Notification {
	return pushkit.Notification{
		Topic:    v.Topic,
		Title:    v.Title,
		Subtitle: v.Subtitle,
		Body:     v.Body,
		Metadata: v.Metadata,
		Priority: v.Priority,
	}
}

func NewNotificationFromPushkit(pk pushkit.Notification) Notification {
	return Notification{
		Topic:    pk.Topic,
		Title:    pk.Title,
		Subtitle: pk.Subtitle,
		Body:     pk.Body,
		Metadata: pk.Metadata,
		Priority: pk.Priority,
		SenderID: nil,
	}
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
