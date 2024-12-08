package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/proto"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	localCache "git.solsynth.dev/hypernet/passport/pkg/internal/cache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"time"
	"unicode"

	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"

	"gorm.io/gorm/clause"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/datatypes"

	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func GetAccountCacheKey(query any) string {
	return fmt.Sprintf("account-query#%v", query)
}

func CacheAccount(account models.Account) {
	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	ctx := context.Background()

	_ = marshal.Set(
		ctx,
		GetAccountCacheKey(account.Name),
		account,
		store.WithExpiration(30*time.Minute),
		store.WithTags([]string{"account", fmt.Sprintf("user#%d", account.ID)}),
	)
	_ = marshal.Set(
		ctx,
		GetAccountCacheKey(account.ID),
		account,
		store.WithExpiration(30*time.Minute),
		store.WithTags([]string{"account", fmt.Sprintf("user#%d", account.ID)}),
	)
}

func ValidateAccountName(val string, min, max int) bool {
	actualLength := 0
	for _, r := range val {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Hangul, r) {
			actualLength += 2
		} else {
			actualLength += 1
		}
	}
	return actualLength >= min && max >= actualLength
}

func GetAccount(id uint) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{
		BaseModel: models.BaseModel{ID: id},
	}).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func GetAccountList(id []uint) ([]models.Account, error) {
	var accounts []models.Account
	if err := database.C.Where("id IN ?", id).Find(&accounts).Error; err != nil {
		return accounts, err
	}

	return accounts, nil
}

func GetAccountWithName(alias string) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{
		Name: alias,
	}).First(&account).Error; err != nil {
		return account, err
	}

	return account, nil
}

func LookupAccount(probe string) (models.Account, error) {
	var account models.Account
	if err := database.C.Where(models.Account{Name: probe}).First(&account).Error; err == nil {
		return account, nil
	}

	var contact models.AccountContact
	if err := database.C.Where(models.AccountContact{Content: probe}).First(&contact).Error; err == nil {
		if err := database.C.
			Where(models.Account{
				BaseModel: models.BaseModel{ID: contact.AccountID},
			}).First(&account).Error; err == nil {
			return account, err
		}
	}

	return account, fmt.Errorf("account was not found")
}

func SearchAccount(probe string) ([]models.Account, error) {
	probe = "%" + probe + "%"
	var accounts []models.Account
	if err := database.C.Where("name LIKE ? OR nick LIKE ?", probe, probe).Find(&accounts).Error; err != nil {
		return accounts, err
	}
	return accounts, nil
}

func CreateAccount(name, nick, email, password string) (models.Account, error) {
	user := models.Account{
		Name: name,
		Nick: nick,
		Profile: models.AccountProfile{
			Experience: 100,
		},
		Factors: []models.AuthFactor{
			{
				Type:   models.PasswordAuthFactor,
				Secret: HashPassword(password),
			},
			{
				Type:   models.EmailPasswordFactor,
				Secret: uuid.NewString()[:8],
			},
		},
		Contacts: []models.AccountContact{
			{
				Type:       models.EmailAccountContact,
				Content:    email,
				IsPrimary:  true,
				VerifiedAt: nil,
			},
		},
		PermNodes:   datatypes.JSONMap{},
		ConfirmedAt: nil,
	}

	if err := database.C.Create(&user).Error; err != nil {
		return user, err
	} else if viper.GetInt("default_user_group") > 0 {
		database.C.Create(&models.AccountGroupMember{
			AccountID: user.ID,
			GroupID:   uint(viper.GetInt("default_user_group")),
		})
	}

	if tk, err := NewMagicToken(models.ConfirmMagicToken, &user, nil); err != nil {
		return user, err
	} else if err := NotifyMagicToken(tk); err != nil {
		return user, err
	}

	return user, nil
}

func ConfirmAccount(code string) error {
	token, err := ValidateMagicToken(code, models.ConfirmMagicToken)
	if err != nil {
		return err
	} else if token.AccountID == nil {
		return fmt.Errorf("magic token didn't assign a valid account")
	}

	var user models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: *token.AccountID},
	}).First(&user).Error; err != nil {
		return err
	}

	if err = ForceConfirmAccount(user); err != nil {
		return err
	} else {
		database.C.Delete(&token)
	}

	return nil
}

func ForceConfirmAccount(user models.Account) error {
	user.ConfirmedAt = lo.ToPtr(time.Now())

	for k, v := range viper.GetStringMap("permissions.verified") {
		if val, ok := user.PermNodes[k]; !ok {
			user.PermNodes[k] = v
		} else {
			user.PermNodes[k] = val
		}
	}

	if err := database.C.Save(&user).Error; err != nil {
		return err
	}

	InvalidAuthCacheWithUser(user.ID)

	return nil
}

func CheckAbleToDeleteAccount(user models.Account) error {
	if user.AutomatedID != nil {
		return fmt.Errorf("bot cannot request delete account, head to developer portal and dispose bot")
	}

	var count int64
	if err := database.C.
		Where("account_id = ?", user.ID).
		Where("expired_at < ?", time.Now()).
		Where("type = ?", models.DeleteAccountMagicToken).
		Model(&models.MagicToken{}).
		Count(&count).Error; err != nil {
		return fmt.Errorf("unable to check delete account ability: %v", err)
	} else if count > 0 {
		return fmt.Errorf("you requested delete account recently")
	}

	return nil
}

func RequestDeleteAccount(user models.Account) error {
	if tk, err := NewMagicToken(
		models.DeleteAccountMagicToken,
		&user,
		lo.ToPtr(time.Now().Add(24*time.Hour)),
	); err != nil {
		return err
	} else if err := NotifyMagicToken(tk); err != nil {
		log.Error().
			Err(err).
			Str("code", tk.Code).
			Uint("user", user.ID).
			Msg("Failed to notify delete account magic token...")
	}

	return nil
}

func ConfirmDeleteAccount(code string) error {
	token, err := ValidateMagicToken(code, models.DeleteAccountMagicToken)
	if err != nil {
		return err
	} else if token.AccountID == nil {
		return fmt.Errorf("magic token didn't assign a valid account")
	}

	if err := DeleteAccount(*token.AccountID); err != nil {
		return err
	} else {
		database.C.Delete(&token)
	}

	return nil
}

func CheckAbleToResetPassword(user models.Account) error {
	var count int64
	if err := database.C.
		Where("account_id = ?", user.ID).
		Where("expired_at < ?", time.Now()).
		Where("type = ?", models.ResetPasswordMagicToken).
		Model(&models.MagicToken{}).
		Count(&count).Error; err != nil {
		return fmt.Errorf("unable to check reset password ability: %v", err)
	} else if count > 0 {
		return fmt.Errorf("you requested reset password recently")
	}

	return nil
}

func RequestResetPassword(user models.Account) error {
	if tk, err := NewMagicToken(
		models.ResetPasswordMagicToken,
		&user,
		lo.ToPtr(time.Now().Add(24*time.Hour)),
	); err != nil {
		return err
	} else if err := NotifyMagicToken(tk); err != nil {
		log.Error().
			Err(err).
			Str("code", tk.Code).
			Uint("user", user.ID).
			Msg("Failed to notify password reset magic token...")
	}

	return nil
}

func ConfirmResetPassword(code, newPassword string) error {
	token, err := ValidateMagicToken(code, models.ResetPasswordMagicToken)
	if err != nil {
		return err
	} else if token.AccountID == nil {
		return fmt.Errorf("magic token didn't assign a valid account")
	}

	factor, err := GetPasswordTypeFactor(*token.AccountID)
	if err != nil {
		factor = models.AuthFactor{
			Type:      models.PasswordAuthFactor,
			Secret:    HashPassword(newPassword),
			AccountID: *token.AccountID,
		}
	} else {
		factor.Secret = HashPassword(newPassword)
	}

	if err = database.C.Save(&factor).Error; err != nil {
		return err
	} else {
		database.C.Delete(&token)
	}

	return nil
}

func DeleteAccount(id uint) error {
	tx := database.C.Begin()

	if err := tx.Select(clause.Associations).Delete(&models.Account{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	} else {
		InvalidAuthCacheWithUser(id)
		_, _ = proto.NewDirectoryServiceClient(gap.Nx.GetNexusGrpcConn()).BroadcastEvent(context.Background(), &proto.EventInfo{
			Event: "deletion",
			Data: nex.EncodeMap(map[string]any{
				"type": "account",
				"id":   fmt.Sprintf("%d", id),
			}),
		})
	}

	return nil
}

func RecycleUnConfirmAccount() {
	deadline := time.Now().Add(-24 * time.Hour)

	var hitList []models.Account
	if err := database.C.Where("confirmed_at IS NULL AND created_at <= ?", deadline).Find(&hitList).Error; err != nil {
		log.Error().Err(err).Msg("An error occurred while recycling accounts...")
		return
	}

	if len(hitList) > 0 {
		log.Info().Int("count", len(hitList)).Msg("Going to recycle those un-confirmed accounts...")
		for _, entry := range hitList {
			if err := DeleteAccount(entry.ID); err != nil {
				log.Error().Err(err).Msg("An error occurred while recycling accounts...")
			}
		}
	}
}

func SetAccountLastSeen(uid uint) error {
	var profile models.AccountProfile
	if err := database.C.Where("account_id = ?", uid).First(&profile).Error; err != nil {
		return err
	}

	profile.LastSeenAt = lo.ToPtr(time.Now())

	return database.C.Save(&profile).Error
}
