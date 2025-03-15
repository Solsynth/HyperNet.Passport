package models

import "time"

type MagicTokenType = int8

const (
	ConfirmMagicToken = MagicTokenType(iota)
	RegistrationMagicToken
	ResetPasswordMagicToken
	DeleteAccountMagicToken
)

type MagicToken struct {
	BaseModel

	Code           string     `json:"code"`
	Type           int8       `json:"type"`
	AccountID      *uint      `json:"account_id"`
	ExpiredAt      *time.Time `json:"expired_at"`
	LastNotifiedAt *time.Time `json:"last_notified_at"`
}
