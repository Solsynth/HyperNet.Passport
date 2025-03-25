package database

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"gorm.io/gorm"
)

var AutoMaintainRange = []any{
	&models.Account{},
	&models.AccountGroup{},
	&models.AccountGroupMember{},
	&models.AuthFactor{},
	&models.AccountProfile{},
	&models.AccountPage{},
	&models.AccountContact{},
	&models.AccountRelationship{},
	&models.Status{},
	&models.Badge{},
	&models.Realm{},
	&models.RealmMember{},
	&models.AuthTicket{},
	&models.MagicToken{},
	&models.ThirdClient{},
	&models.ActionEvent{},
	&models.Notification{},
	&models.NotificationSubscriber{},
	&models.AuditRecord{},
	&models.ApiKey{},
	&models.CheckInRecord{},
	&models.PreferenceNotification{},
	&models.PreferenceAuth{},
	&models.AbuseReport{},
	&models.Program{},
	&models.ProgramMember{},
	&models.Punishment{},
}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(AutoMaintainRange...); err != nil {
		return err
	}

	return nil
}
