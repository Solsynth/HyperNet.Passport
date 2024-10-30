package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JSONMap = datatypes.JSONType[map[string]any]

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
