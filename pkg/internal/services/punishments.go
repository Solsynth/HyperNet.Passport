package services

import (
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/nexus/pkg/nex/localize"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/rs/zerolog/log"
)

func NewPunishment(in models.Punishment, moderator ...models.Account) (models.Punishment, error) {
	if len(moderator) > 0 {
		in.Moderator = &moderator[0]
		in.ModeratorID = &moderator[0].ID
	}

	// If user got more than 2 strikes, it will upgrade to limited
	if in.Type == models.PunishmentTypeStrike {
		var count int64
		if err := database.C.Model(&models.Punishment{}).
			Where("account_id = ? AND type = ?", in.AccountID, models.PunishmentTypeStrike).
			Count(&count).Error; err != nil {
			return in, err
		}
		if count > 2 {
			in.Type = models.PunishmentTypeLimited
		}
	}

	if err := database.C.Create(&in).Error; err != nil {
		return in, err
	} else {
		var moderator = "System"
		if in.Moderator != nil {
			moderator = fmt.Sprintf("@%s", in.Moderator.Name)
		}
		err = NewNotification(models.Notification{
			Topic:     "passport.punishments",
			Title:     localize.L.GetLocalizedString("subjectPunishmentCreated", in.Account.Language),
			Subtitle:  fmt.Sprintf(localize.L.GetLocalizedString("subtitlePunishment", in.Account.Language), in.ID, moderator),
			Body:      fmt.Sprintf(localize.L.GetLocalizedString("shortBodyPunishmentCreated", in.Account.Language), in.Reason),
			Account:   in.Account,
			AccountID: in.Account.ID,
			Metadata:  map[string]any{"punishment": in},
		})
		if err != nil {
			log.Warn().Err(err).Uint("case", in.ID).Msg("Failed to delivery punishment via notify...")
		}
	}

	return in, nil
}

func EditPunishment(in models.Punishment) (models.Punishment, error) {
	if err := database.C.Save(&in).Error; err != nil {
		return in, err
	} else {
		var moderator = "System"
		if in.Moderator != nil {
			moderator = fmt.Sprintf("@%s", in.Moderator.Name)
		}
		err = NewNotification(models.Notification{
			Topic:     "passport.punishments",
			Title:     localize.L.GetLocalizedString("subjectPunishmentUpdated", in.Account.Language),
			Subtitle:  fmt.Sprintf(localize.L.GetLocalizedString("subtitlePunishment", in.Account.Language), in.ID, moderator),
			Body:      fmt.Sprintf(localize.L.GetLocalizedString("shortBodyPunishmentUpdated", in.Account.Language), in.ID),
			Account:   in.Account,
			AccountID: in.Account.ID,
			Metadata:  map[string]any{"punishment": in},
		})
		if err != nil {
			log.Warn().Err(err).Uint("case", in.ID).Msg("Failed to delivery punishment via notify...")
		}
	}
	return in, nil
}

func DeletePunishment(in models.Punishment) error {
	if err := database.C.Delete(&in).Error; err != nil {
		return err
	} else {
		var moderator = "System"
		if in.Moderator != nil {
			moderator = fmt.Sprintf("@%s", in.Moderator.Name)
		}
		err = NewNotification(models.Notification{
			Topic:     "passport.punishments",
			Title:     localize.L.GetLocalizedString("subjectPunishmentDeleted", in.Account.Language),
			Subtitle:  fmt.Sprintf(localize.L.GetLocalizedString("subtitlePunishment", in.Account.Language), in.ID, moderator),
			Body:      fmt.Sprintf(localize.L.GetLocalizedString("shortBodyPunishmentDeleted", in.Account.Language), in.ID),
			Account:   in.Account,
			AccountID: in.Account.ID,
			Metadata:  map[string]any{"punishment": in},
		})
		if err != nil {
			log.Warn().Err(err).Uint("case", in.ID).Msg("Failed to delivery punishment via notify...")
		}
	}
	return nil
}

func GetPunishment(id uint, preload ...bool) (models.Punishment, error) {
	tx := database.C
	if len(preload) > 0 && preload[0] {
		tx = tx.Preload("Moderator").Preload("Account")
	}

	var punishment models.Punishment
	if err := tx.First(&punishment, id).Error; err != nil {
		return punishment, err
	}
	return punishment, nil
}

func GetMadePunishment(id uint, moderator models.Account) (models.Punishment, error) {
	var punishment models.Punishment
	if err := database.C.Where("id = ? AND moderator_id = ?", id, moderator.ID).First(&punishment).Error; err != nil {
		return punishment, err
	}
	return punishment, nil
}

func ListPunishments(user models.Account) ([]models.Punishment, error) {
	var punishments []models.Punishment
	if err := database.C.
		Where("account_id = ? AND (expired_at IS NULL OR expired_at <= ?)", user.ID, time.Now()).
		Preload("Moderator").
		Order("created_at DESC").
		Find(&punishments).Error; err != nil {
		return nil, err
	}
	return punishments, nil
}

func CountAllPunishments() (int64, error) {
	var count int64
	if err := database.C.
		Model(&models.Punishment{}).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func ListAllPunishments(take, offset int) ([]models.Punishment, error) {
	var punishments []models.Punishment
	if err := database.C.
		Preload("Account").
		Preload("Moderator").
		Order("created_at DESC").
		Take(take).Offset(offset).
		Find(&punishments).Error; err != nil {
		return nil, err
	}
	return punishments, nil
}

func CountMadePunishments(moderator models.Account) (int64, error) {
	var count int64
	if err := database.C.
		Model(&models.Punishment{}).
		Where("moderator_id = ?", moderator.ID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func ListMadePunishments(moderator models.Account, take, offset int) ([]models.Punishment, error) {
	var punishments []models.Punishment
	if err := database.C.
		Where("moderator_id = ?", moderator.ID).
		Preload("Account").
		Order("created_at DESC").
		Take(take).Offset(offset).
		Find(&punishments).Error; err != nil {
		return nil, err
	}
	return punishments, nil
}

func CheckLoginAbility(user models.Account) error {
	var punishments []models.Punishment
	if err := database.C.Where("account_id = ? AND (expired_at IS NULL OR expired_at <= ?)", user.ID, time.Now()).
		Find(&punishments).Error; err != nil {
		return fmt.Errorf("failed to get punishments: %v", err)
	}

	for _, punishment := range punishments {
		if punishment.Type == models.PunishmentTypeDisabled {
			return fmt.Errorf("account has been fully disabled due to: %s (case #%d)", punishment.Reason, punishment.ID)
		}
		// Limited punishment with no permissions override is fully limited
		// Refer https://solsynth.dev/terms/basic-law#provision-and-discontinuation-of-services
		if punishment.Type == models.PunishmentTypeLimited && len(punishment.PermNodes) == 0 {
			return fmt.Errorf("account has been limited login due to: %s (case #%d)", punishment.Reason, punishment.ID)
		}
	}

	return nil
}
