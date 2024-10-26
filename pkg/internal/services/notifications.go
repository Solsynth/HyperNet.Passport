package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"
	"reflect"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/gap"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
)

func AddNotifySubscriber(user models.Account, provider, id, tk, ua string) (models.NotificationSubscriber, error) {
	var prev models.NotificationSubscriber
	var subscriber models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		DeviceID:  id,
		AccountID: user.ID,
	}); err != nil {
		subscriber = models.NotificationSubscriber{
			UserAgent:   ua,
			Provider:    provider,
			DeviceID:    id,
			DeviceToken: tk,
			AccountID:   user.ID,
		}
	} else {
		prev = subscriber
	}

	subscriber.UserAgent = ua
	subscriber.Provider = provider
	subscriber.DeviceToken = tk

	var err error
	if !reflect.DeepEqual(subscriber, prev) {
		err = database.C.Save(&subscriber).Error
	}

	return subscriber, err
}

// NewNotification will create a notification and push via the push method it
// Pleases provide the notification with the account field is not empty
func NewNotification(notification models.Notification) error {
	if ok := CheckNotificationNotifiable(notification.Account, notification.Topic); !ok {
		log.Info().Str("topic", notification.Topic).Uint("uid", notification.AccountID).Msg("Notification dismissed by user...")
		return nil
	}

	if err := database.C.Save(&notification).Error; err != nil {
		return err
	}
	if err := PushNotification(notification, true); err != nil {
		return err
	}

	return nil
}

func NewNotificationBatch(notifications []models.Notification) error {
	if len(notifications) == 0 {
		return nil
	}

	notifiable := CheckNotificationNotifiableBatch(lo.Map(notifications, func(item models.Notification, index int) models.Account {
		return item.Account
	}), notifications[0].Topic)

	notifications = lo.Filter(notifications, func(item models.Notification, index int) bool {
		return notifiable[index]
	})

	if err := database.C.CreateInBatches(notifications, 1000).Error; err != nil {
		return err
	}

	PushNotificationBatch(notifications, true)
	return nil
}

// PushNotification will push a notification to the user, via websocket, firebase, or APNs
// Please provide the notification with the account field is not empty
func PushNotification(notification models.Notification, skipNotifiableCheck ...bool) error {
	if len(skipNotifiableCheck) == 0 || !skipNotifiableCheck[0] {
		if ok := CheckNotificationNotifiable(notification.Account, notification.Topic); !ok {
			log.Info().Str("topic", notification.Topic).Uint("uid", notification.AccountID).Msg("Notification dismissed by user...")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := proto.NewStreamServiceClient(gap.Nx.GetNexusGrpcConn()).PushStream(ctx, &proto.PushStreamRequest{
		UserId: lo.ToPtr(uint64(notification.AccountID)),
		Body: nex.WebSocketPackage{
			Action:  "notifications.new",
			Payload: notification,
		}.Marshal(),
	})
	if err != nil {
		return fmt.Errorf("failed to push via websocket: %v", err)
	}

	// Skip push notification
	if GetStatusDisturbable(notification.AccountID) != nil {
		return nil
	}

	var subscribers []models.NotificationSubscriber
	if err := database.C.Where(&models.NotificationSubscriber{
		AccountID: notification.AccountID,
	}).Find(&subscribers).Error; err != nil {
		return err
	}

	var providers []string
	var tokens []string
	for _, subscriber := range subscribers {
		providers = append(providers, subscriber.Provider)
		tokens = append(tokens, subscriber.DeviceToken)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = gap.Px.PushNotifyBatch(pushkit.NotificationPushBatchRequest{
		Providers: providers,
		Tokens:    tokens,
		Notification: pushkit.Notification{
			Topic:    notification.Topic,
			Title:    notification.Title,
			Subtitle: notification.Subtitle,
			Body:     notification.Body,
			Metadata: notification.Metadata,
		},
	})

	return err
}

func PushNotificationBatch(notifications []models.Notification, skipNotifiableCheck ...bool) {
	if len(notifications) == 0 {
		return
	}

	var accountIdx []uint
	if len(skipNotifiableCheck) == 0 || !skipNotifiableCheck[0] {
		notifiable := CheckNotificationNotifiableBatch(lo.Map(notifications, func(item models.Notification, index int) models.Account {
			return item.Account
		}), notifications[0].Topic)
		accountIdx = lo.Map(
			lo.Filter(notifications, func(item models.Notification, index int) bool {
				return notifiable[index]
			}),
			func(item models.Notification, index int) uint {
				return item.AccountID
			},
		)
	} else {
		accountIdx = lo.Map(
			notifications,
			func(item models.Notification, index int) uint {
				return item.AccountID
			},
		)
	}

	if len(accountIdx) == 0 {
		return
	}

	var subscribers []models.NotificationSubscriber
	database.C.Where("account_id IN ?", accountIdx).Find(&subscribers)

	stream := proto.NewStreamServiceClient(gap.Nx.GetNexusGrpcConn())
	for _, notification := range notifications {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, _ = stream.PushStream(ctx, &proto.PushStreamRequest{
			UserId: lo.ToPtr(uint64(notification.AccountID)),
			Body: nex.WebSocketPackage{
				Action:  "notifications.new",
				Payload: notification,
			}.Marshal(),
		})
		cancel()

		// Skip push notification
		if GetStatusDisturbable(notification.AccountID) != nil {
			continue
		}

		var providers []string
		var tokens []string
		for _, subscriber := range lo.Filter(subscribers, func(item models.NotificationSubscriber, index int) bool {
			return item.AccountID == notification.AccountID
		}) {
			providers = append(providers, subscriber.Provider)
			tokens = append(tokens, subscriber.DeviceToken)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		_ = gap.Px.PushNotifyBatch(pushkit.NotificationPushBatchRequest{
			Providers: providers,
			Tokens:    tokens,
			Notification: pushkit.Notification{
				Topic:    notification.Topic,
				Title:    notification.Title,
				Subtitle: notification.Subtitle,
				Body:     notification.Body,
				Metadata: notification.Metadata,
			},
		})
		cancel()
	}
}
