package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/paperclip/pkg/filekit"
	pproto "git.solsynth.dev/hypernet/paperclip/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func ListCommunityRealm() ([]models.Realm, error) {
	var realms []models.Realm
	if err := database.C.Where(&models.Realm{
		IsCommunity: true,
	}).Order("popularity DESC").Find(&realms).Error; err != nil {
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

func ListAvailableRealm(user models.Account, includePublic ...bool) ([]models.Realm, error) {
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

	tx := database.C
	if len(includePublic) > 0 && includePublic[0] {
		tx = tx.Where("is_public = ?", true)
	}

	if err := tx.Where("id IN ?", idx).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func GetRealmWithAlias(alias string) (models.Realm, error) {
	tx := database.C.Where("alias = ?", alias)

	numericId, err := strconv.Atoi(alias)
	if err == nil {
		tx.Or("id = ?", numericId)
	}

	var realm models.Realm
	if err := tx.First(&realm).Error; err != nil {
		return realm, err
	}
	return realm, nil
}

func NewRealm(realm models.Realm, user models.Account) (models.Realm, error) {
	realm.Members = []models.RealmMember{
		{AccountID: user.ID, PowerLevel: 100},
	}

	var attachments []string
	if realm.Avatar != nil && len(*realm.Avatar) > 0 {
		attachments = append(attachments, *realm.Avatar)
	}
	if realm.Banner != nil && len(*realm.Banner) > 0 {
		attachments = append(attachments, *realm.Banner)
	}
	if len(attachments) > 0 {
		filekit.CountAttachmentUsage(gap.Nx, &pproto.UpdateUsageRequest{
			Rid:   attachments,
			Delta: 1,
		})
	}

	err := database.C.Save(&realm).Error
	return realm, err
}

func CountRealmMember(realmId uint) (int64, error) {
	var count int64
	if err := database.C.Where(&models.RealmMember{
		RealmID: realmId,
	}).Model(&models.RealmMember{}).Count(&count).Error; err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

func ListRealmMember(realmId uint, take int, offset int) ([]models.RealmMember, error) {
	var members []models.RealmMember

	if err := database.C.
		Limit(take).Offset(offset).
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
	var member models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		AccountID: affected.ID,
		RealmID:   target.ID,
	}).First(&member).Error; err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if !target.IsCommunity {
		if member, err := GetRealmMember(user.ID, target.ID); err != nil {
			return fmt.Errorf("only realm member can add people: %v", err)
		} else if member.PowerLevel < 50 {
			return fmt.Errorf("only realm moderator can add member")
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

	member = models.RealmMember{
		RealmID:   target.ID,
		AccountID: affected.ID,
	}

	err := database.C.Save(&member).Error
	if err == nil {
		database.C.Model(&models.Realm{}).
			Where("id = ?", target.ID).
			Update("popularity", gorm.Expr("popularity + ?", models.RealmPopularityMemberFactor))
	}

	return err
}

func RemoveRealmMember(user models.Account, affected models.RealmMember, target models.Realm) error {
	if user.ID != affected.AccountID {
		if member, err := GetRealmMember(user.ID, target.ID); err != nil {
			return fmt.Errorf("only realm member can remove other member: %v", err)
		} else if member.PowerLevel < 50 {
			return fmt.Errorf("only realm moderator can kick member")
		}
	}

	if err := database.C.Delete(&affected).Error; err != nil {
		return err
	}

	database.C.Model(&models.Realm{}).
		Where("id = ?", target.ID).
		Update("popularity", gorm.Expr("popularity - ?", models.RealmPopularityMemberFactor))

	return nil
}

func EditRealm(realm, og models.Realm) (models.Realm, error) {
	err := database.C.Save(&realm).Error
	if err == nil {
		var minusAttachments, plusAttachments []string
		if realm.Avatar != og.Avatar && realm.Avatar != nil {
			minusAttachments = append(minusAttachments, *og.Avatar)
			plusAttachments = append(plusAttachments, *realm.Avatar)
		}
		if realm.Banner != og.Banner && realm.Banner != nil {
			minusAttachments = append(minusAttachments, *og.Banner)
			plusAttachments = append(plusAttachments, *realm.Banner)
		}
		if len(minusAttachments) > 0 {
			filekit.CountAttachmentUsage(gap.Nx, &pproto.UpdateUsageRequest{
				Rid:   minusAttachments,
				Delta: -1,
			})
		}
		if len(plusAttachments) > 0 {
			filekit.CountAttachmentUsage(gap.Nx, &pproto.UpdateUsageRequest{
				Rid:   plusAttachments,
				Delta: 1,
			})
		}
	}
	return realm, err
}

func DeleteRealm(realm models.Realm) error {
	tx := database.C.Begin()
	if err := tx.Where("realm_id = ?", realm.ID).Delete(&models.RealmMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&realm).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	} else {
		var attachments []string
		if realm.Avatar != nil && len(*realm.Avatar) > 0 {
			attachments = append(attachments, *realm.Avatar)
		}
		if realm.Banner != nil && len(*realm.Banner) > 0 {
			attachments = append(attachments, *realm.Banner)
		}
		if len(attachments) > 0 {
			filekit.CountAttachmentUsage(gap.Nx, &pproto.UpdateUsageRequest{
				Rid:   attachments,
				Delta: -1,
			})
		}

		conn := gap.Nx.GetNexusGrpcConn()
		_, _ = proto.NewDirectoryServiceClient(conn).BroadcastEvent(context.Background(), &proto.EventInfo{
			Event: "deletion",
			Data: nex.EncodeMap(map[string]any{
				"type": "realm",
				"id":   realm.ID,
			}),
		})
	}
	return nil
}
