package models

import "gorm.io/datatypes"

type AccountGroup struct {
	BaseModel

	Name      string            `json:"name"`
	PermNodes datatypes.JSONMap `json:"perm_nodes"`
}

type AccountGroupMember struct {
	BaseModel

	Account   Account      `json:"account"`
	Group     AccountGroup `json:"group"`
	AccountID uint         `json:"account_id"`
	GroupID   uint         `json:"group_id"`
}
