package services

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/authkit/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"github.com/samber/lo"
)

func ListCommunityRealm() ([]models.Realm, error) {
	var realms []models.Realm
	if err := database.C.Where(&models.Realm{
		IsCommunity: true,
	}).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func ListOwnedRealm(user models.Account) ([]models.Realm, error) {
	var realms []models.Realm
	if err := database.C.Where(&models.Realm{AccountID: user.ID}).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func ListAvailableRealm(user models.Account) ([]models.Realm, error) {
	var realms []models.Realm
	var members []models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		AccountID: user.ID,
	}).Find(&members).Error; err != nil {
		return realms, err
	}

	idx := lo.Map(members, func(item models.RealmMember, index int) uint {
		return item.RealmID
	})

	if err := database.C.Where("id IN ?", idx).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func GetRealmWithAlias(alias string) (models.Realm, error) {
	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		Alias: alias,
	}).First(&realm).Error; err != nil {
		return realm, err
	}
	return realm, nil
}

func NewRealm(realm models.Realm, user models.Account) (models.Realm, error) {
	realm.Members = []models.RealmMember{
		{AccountID: user.ID, PowerLevel: 100},
	}

	err := database.C.Save(&realm).Error
	return realm, err
}

func ListRealmMember(realmId uint) ([]models.RealmMember, error) {
	var members []models.RealmMember

	if err := database.C.
		Where(&models.RealmMember{RealmID: realmId}).
		Preload("Account").
		Find(&members).Error; err != nil {
		return members, err
	}

	return members, nil
}

func GetRealmMember(userId uint, realmId uint) (models.RealmMember, error) {
	var member models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		AccountID: userId,
		RealmID:   realmId,
	}).Find(&member).Error; err != nil {
		return member, err
	}
	return member, nil
}

func AddRealmMember(user models.Account, affected models.Account, target models.Realm) error {
	if !target.IsPublic && !target.IsCommunity {
		if member, err := GetRealmMember(user.ID, target.ID); err != nil {
			return fmt.Errorf("only realm member can add people: %v", err)
		} else if member.PowerLevel < 50 {
			return fmt.Errorf("only realm moderator can add people")
		}
		rel, err := GetRelationWithTwoNode(affected.ID, user.ID)
		if err != nil || HasPermNodeWithDefault(
			rel.PermNodes,
			"RealmAdd",
			true,
			rel.Status == models.RelationshipFriend,
		) {
			return fmt.Errorf("you unable to add this user to your realm")
		}
	}

	member := models.RealmMember{
		RealmID:   target.ID,
		AccountID: affected.ID,
	}
	err := database.C.Save(&member).Error
	return err
}

func RemoveRealmMember(user models.Account, affected models.Account, target models.Realm) error {
	if user.ID != affected.ID {
		if member, err := GetRealmMember(user.ID, target.ID); err != nil {
			return fmt.Errorf("only realm member can remove other member: %v", err)
		} else if member.PowerLevel < 50 {
			return fmt.Errorf("only realm moderator can invite people")
		}
	}

	var member models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		RealmID:   target.ID,
		AccountID: affected.ID,
	}).First(&member).Error; err != nil {
		return err
	}

	return database.C.Delete(&member).Error
}

func EditRealm(realm models.Realm) (models.Realm, error) {
	err := database.C.Save(&realm).Error
	return realm, err
}

func DeleteRealm(realm models.Realm) error {
	return database.C.Delete(&realm).Error
}
