package services

import (
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/cachekit"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"github.com/rs/zerolog/log"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

func GetAuthPreference(account models.Account) (models.PreferenceAuth, error) {
	var auth models.PreferenceAuth
	if err := database.C.Where("account_id = ?", account.ID).First(&auth).Error; err != nil {
		return auth, err
	}

	return auth, nil
}

func UpdateAuthPreference(account models.Account, config models.AuthConfig) (models.PreferenceAuth, error) {
	var auth models.PreferenceAuth
	var err error
	if auth, err = GetAuthPreference(account); err != nil {
		auth = models.PreferenceAuth{
			AccountID: account.ID,
			Config:    datatypes.NewJSONType(config),
		}
	} else {
		auth.Config = datatypes.NewJSONType(config)
	}

	err = database.C.Save(&auth).Error
	return auth, err
}

func KgNotifyPreferenceCache(accountId uint) string {
	return fmt.Sprintf("notification-preference#%d", accountId)
}

func GetNotifyPreference(account models.Account) (models.PreferenceNotification, error) {
	var notification models.PreferenceNotification
	if val, err := cachekit.Get[models.PreferenceNotification](
		gap.Ca,
		KgNotifyPreferenceCache(account.ID),
	); err == nil {
		return val, nil
	}
	if err := database.C.Where("account_id = ?", account.ID).First(&notification).Error; err != nil {
		return notification, err
	}
	CacheNotifyPreference(notification)
	return notification, nil
}

func CacheNotifyPreference(prefs models.PreferenceNotification) {
	cachekit.Set[models.PreferenceNotification](
		gap.Ca,
		KgNotifyPreferenceCache(prefs.AccountID),
		prefs,
		time.Minute*60,
		fmt.Sprintf("user#%d", prefs.AccountID),
	)
}

func UpdateNotifyPreference(account models.Account, config map[string]bool) (models.PreferenceNotification, error) {
	var notification models.PreferenceNotification
	var err error
	if notification, err = GetNotifyPreference(account); err != nil {
		notification = models.PreferenceNotification{
			AccountID: account.ID,
			Config:    lo.MapValues(config, func(v bool, k string) any { return v }),
		}
	} else {
		notification.Config = lo.MapValues(config, func(v bool, k string) any { return v })
	}

	err = database.C.Save(&notification).Error
	if err == nil {
		CacheNotifyPreference(notification)
	}

	return notification, err
}

func CheckNotificationNotifiable(account models.Account, topic string) bool {
	var notification models.PreferenceNotification
	if val, err := cachekit.Get[models.PreferenceNotification](
		gap.Ca,
		KgNotifyPreferenceCache(account.ID),
	); err == nil {
		notification = val
	} else {
		if err := database.C.Where("account_id = ?", account.ID).First(&notification).Error; err != nil {
			return true
		}
		CacheNotifyPreference(notification)
	}

	if val, ok := notification.Config[topic]; ok {
		if status, ok := val.(bool); ok {
			return status
		}
	}
	return true
}

func CheckNotificationNotifiableBatch(accounts []models.Account, topic string) []bool {
	notifiable := make([]bool, len(accounts))
	var queryNeededIdx []uint
	notificationMap := make(map[uint]models.PreferenceNotification)

	// Check cache for each account
	for _, account := range accounts {
		cacheKey := KgNotifyPreferenceCache(account.ID)
		if val, err := cachekit.Get[models.PreferenceNotification](gap.Ca, cacheKey); err == nil {
			notificationMap[account.ID] = val
		} else {
			// Add to the list of accounts that need to be queried
			queryNeededIdx = append(queryNeededIdx, account.ID)
		}
	}

	// Query the database for missing account IDs
	if len(queryNeededIdx) > 0 {
		var dbNotifications []models.PreferenceNotification
		if err := database.C.Where("account_id IN ?", queryNeededIdx).Find(&dbNotifications).Error; err != nil {
			// Handle error by returning false for accounts without cached notifications
			return lo.Map(accounts, func(item models.Account, index int) bool {
				return true
			})
		}

		// Cache the newly fetched notifications and add to the notificationMap
		for _, notification := range dbNotifications {
			notificationMap[notification.AccountID] = notification
			CacheNotifyPreference(notification) // Cache the result
		}
	}

	log.Debug().Any("notifiable", notificationMap).Msg("Fetched notifiable status...")

	// Process the notifiable status for the fetched notifications
	for idx, account := range accounts {
		if notification, exists := notificationMap[account.ID]; exists {
			if val, ok := notification.Config[topic]; ok {
				if status, ok := val.(bool); ok {
					notifiable[idx] = status
					continue
				}
			}
			notifiable[idx] = true
		} else {
			notifiable[idx] = true
		}
	}

	return notifiable
}
