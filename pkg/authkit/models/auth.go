package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type AuthConfig struct {
	MaximumAuthSteps int `json:"maximum_auth_steps" validate:"required,min=1,max=99"`
}

type AuthFactorType = int8

const (
	PasswordAuthFactor = AuthFactorType(iota)
	EmailPasswordFactor
	TimeOtpFactor
	InAppNotifyFactor
)

type AuthFactor struct {
	BaseModel

	Type   int8              `json:"type"`
	Secret string            `json:"-"`
	Config datatypes.JSONMap `json:"config"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}

type AuthTicket struct {
	BaseModel

	Location     string                      `json:"location"`
	IpAddress    string                      `json:"ip_address"`
	UserAgent    string                      `json:"user_agent"`
	StepRemain   int                         `json:"step_remain"`
	Claims       datatypes.JSONSlice[string] `json:"claims"`
	Audiences    datatypes.JSONSlice[string] `json:"audiences"`
	FactorTrail  datatypes.JSONSlice[int]    `json:"factor_trail"`
	GrantToken   *string                     `json:"grant_token"`
	AccessToken  *string                     `json:"access_token"`
	RefreshToken *string                     `json:"refresh_token"`
	ExpiredAt    *time.Time                  `json:"expired_at"`
	AvailableAt  *time.Time                  `json:"available_at"`
	LastGrantAt  *time.Time                  `json:"last_grant_at"`
	Nonce        *string                     `json:"nonce"`
	ClientID     *uint                       `json:"client_id"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}

func (v AuthTicket) IsAvailable() error {
	if v.StepRemain > 0 {
		return fmt.Errorf("ticket isn't authenticated yet")
	}
	if v.AvailableAt != nil && time.Now().Unix() < v.AvailableAt.Unix() {
		return fmt.Errorf("ticket isn't available yet")
	}
	if v.ExpiredAt != nil && time.Now().Unix() > v.ExpiredAt.Unix() {
		return fmt.Errorf("ticket expired")
	}

	return nil
}

func (v AuthTicket) IsCanBeAvailble() error {
	if v.StepRemain > 0 {
		return fmt.Errorf("ticket isn't authenticated yet")
	}

	return nil
}

type AuthContext struct {
	Ticket  AuthTicket `json:"ticket"`
	Account Account    `json:"account"`
}
