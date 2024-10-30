package models

import "gorm.io/datatypes"

type RelationshipStatus = int8

const (
	RelationshipPending = RelationshipStatus(iota)
	RelationshipFriend
	RelationshipBlocked
	RelationshipWaiting
)

type AccountRelationship struct {
	BaseModel

	AccountID uint               `json:"account_id"`
	RelatedID uint               `json:"related_id"`
	Account   Account            `json:"account"`
	Related   Account            `json:"related"`
	Status    RelationshipStatus `json:"status"`
	PermNodes datatypes.JSONMap  `json:"perm_nodes"`
}
