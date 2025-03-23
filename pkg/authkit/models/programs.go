package models

import (
	"time"

	"gorm.io/datatypes"
)

type ProgramPrice struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type ProgramBadge struct {
	Type     string         `json:"type"`
	Metadata map[string]any `json:"metadata"`
}

type ProgramGroup struct {
	ID uint `json:"id"`
}

type Program struct {
	BaseModel

	Name           string                           `json:"name"`
	Description    string                           `json:"description"`
	Alias          string                           `json:"alias" gorm:"uniqueIndex"`
	ExpRequirement int64                            `json:"exp_requirement"`
	Price          datatypes.JSONType[ProgramPrice] `json:"price"`
	Badge          datatypes.JSONType[ProgramBadge] `json:"badge"`
	Group          datatypes.JSONType[ProgramGroup] `json:"group"`
	Appearance     datatypes.JSONMap                `json:"appearance"`
}

type ProgramMember struct {
	BaseModel

	LastPaid  *time.Time `json:"last_paid"`
	Account   Account    `json:"account"`
	AccountID uint       `json:"account_id"`
	Program   Program    `json:"program"`
	ProgramID uint       `json:"program_id"`
}
