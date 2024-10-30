package models

import "gorm.io/datatypes"

type ThirdClient struct {
	BaseModel

	Alias       string                      `json:"alias" gorm:"uniqueIndex"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Secret      string                      `json:"secret"`
	Urls        datatypes.JSONSlice[string] `json:"urls"`
	Callbacks   datatypes.JSONSlice[string] `json:"callbacks"`
	IsDraft     bool                        `json:"is_draft"`
	AccountID   *uint                       `json:"account_id"`
}
