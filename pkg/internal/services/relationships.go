package services

import (
	"errors"
	"fmt"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"gorm.io/gorm"
)

func ListAllRelationship(user uint) ([]models.AccountRelationship, error) {
	var relationships []models.AccountRelationship
	if err := database.C.
		Where("account_id = ?", user).
		Preload("Account").
		Preload("Related").
		Find(&relationships).Error; err != nil {
		return relationships, err
	}

	return relationships, nil
}

func ListRelationshipWithFilter(user uint, status models.RelationshipStatus) ([]models.AccountRelationship, error) {
	var relationships []models.AccountRelationship
	if err := database.C.
		Where("account_id = ? AND status = ?", user, status).
		Preload("Account").
		Preload("Related").
		Find(&relationships).Error; err != nil {
		return relationships, err
	}

	return relationships, nil
}

func GetRelationship(otherId uint) (models.AccountRelationship, error) {
	var relationship models.AccountRelationship
	if err := database.C.
		Where(&models.AccountRelationship{AccountID: otherId}).
		Preload("Account").
		Preload("Related").
		First(&relationship).Error; err != nil {
		return relationship, err
	}

	return relationship, nil
}

func GetRelationWithTwoNode(userId, relatedId uint, noPreload ...bool) (models.AccountRelationship, error) {
	var tx *gorm.DB
	if len(noPreload) > 0 && noPreload[0] {
		tx = database.C
	} else {
		tx = database.C.Preload("Account").Preload("Related")
	}

	var relationship models.AccountRelationship
	if err := tx.
		Where(&models.AccountRelationship{AccountID: userId, RelatedID: relatedId}).
		First(&relationship).Error; err != nil {
		return relationship, err
	}

	return relationship, nil
}

func EditRelationship(relationship models.AccountRelationship) (models.AccountRelationship, error) {
	if err := database.C.Save(&relationship).Error; err != nil {
		return relationship, err
	}
	return relationship, nil
}

func DeleteRelationship(relationship models.AccountRelationship) error {
	if err := database.C.Delete(&relationship).Error; err != nil {
		return err
	}
	return nil
}

func NewBlockship(userA models.Account, userB models.Account) (models.AccountRelationship, error) {
	var err error
	var rel models.AccountRelationship
	if rel, err = GetRelationWithTwoNode(userA.ID, userB.ID, true); err == nil {
		rel.Status = models.RelationshipBlocked
	} else {
		rel = models.AccountRelationship{
			AccountID: userA.ID,
			RelatedID: userB.ID,
			Status:    models.RelationshipBlocked,
		}
	}

	if err := database.C.Save(&rel).Error; err != nil {
		return rel, err
	}

	return rel, nil
}

func NewFriend(userA models.Account, userB models.Account, skipPending ...bool) (models.AccountRelationship, error) {
	relA := models.AccountRelationship{
		AccountID: userA.ID,
		RelatedID: userB.ID,
		Status:    models.RelationshipWaiting,
	}
	relB := models.AccountRelationship{
		AccountID: userB.ID,
		RelatedID: userA.ID,
		Status:    models.RelationshipPending,
	}

	if len(skipPending) > 0 && skipPending[0] {
		relA.Status = models.RelationshipFriend
		relB.Status = models.RelationshipFriend
	}

	if userA.ID == userB.ID {
		return relA, fmt.Errorf("unable to make relationship with yourself")
	} else if _, err := GetRelationWithTwoNode(userA.ID, userB.ID, true); err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return relA, fmt.Errorf("unable to recreate a relationship with that user")
	}

	if err := database.C.Save(&relA).Error; err != nil {
		return relA, err
	} else if err = database.C.Save(&relB).Error; err != nil {
		return relA, err
	} else {
		_ = NewNotification(models.Notification{
			Title:     "New Friend Request",
			Subtitle:  fmt.Sprintf("New friend request from %s", userA.Name),
			Body:      fmt.Sprintf("You got a new friend request from %s. Go to your account page and decide how to deal it.", userA.Nick),
			Account:   userB,
			AccountID: userB.ID,
		})
	}

	return relA, nil
}

func HandleFriend(userA models.Account, userB models.Account, isAccept bool) error {
	relA, err := GetRelationWithTwoNode(userA.ID, userB.ID, true)
	if err != nil {
		return fmt.Errorf("relationship was not found: %v", err)
	} else if relA.Status != models.RelationshipPending {
		return fmt.Errorf("relationship already handled")
	}

	if isAccept {
		relA.Status = models.RelationshipFriend
	} else {
		relA.Status = models.RelationshipBlocked
	}

	if err := database.C.Save(&relA).Error; err != nil {
		return err
	}

	relB, err := GetRelationWithTwoNode(userB.ID, userA.ID, true)
	if err == nil && relB.Status == models.RelationshipWaiting {
		relB.Status = models.RelationshipFriend
		if err := database.C.Save(&relB).Error; err != nil {
			return err
		}

		_ = NewNotification(models.Notification{
			Title:     "Friend Request Processed",
			Subtitle:  fmt.Sprintf("Your friend request to %s has been processsed.", userA.Name),
			Body:      fmt.Sprintf("Your relationship status with %s has been updated, go check it out!", userA.Nick),
			Account:   userB,
			AccountID: userB.ID,
		})
	}

	return nil
}
