package models

import "gorm.io/datatypes"

type CheckInRecord struct {
	BaseModel

	ResultTier       int     `json:"result_tier"`
	ResultExperience int     `json:"result_experience"`
	ResultCoin       float64 `json:"result_coin"`

	// The result modifiers are some random tips that will show up in the client;
	// This field is to use to make sure the tips will be the same when the client is reloaded.
	// For now, this modifier slice will contain four random numbers from 0 to 1024.
	// Client should mod this modifier by the length of total available tips.
	ResultModifiers datatypes.JSONSlice[int] `json:"result_modifiers"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}
