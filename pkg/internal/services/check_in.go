package services

import (
	"errors"
	"fmt"
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func CheckCanCheckIn(user models.Account) error {
	probe := time.Now().Format("2006-01-02")

	var record models.CheckInRecord
	if err := database.C.Where("account_id = ? AND created_at::date = ?", user.ID, probe).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("unable get check in record: %v", err)
	}
	return fmt.Errorf("today's check in record exists")
}

func GetTodayCheckIn(user models.Account) (models.CheckInRecord, error) {
	probe := time.Now().Format("2006-01-02")

	var record models.CheckInRecord
	if err := database.C.Where("account_id = ? AND created_at::date = ?", user.ID, probe).First(&record).Error; err != nil {
		return record, fmt.Errorf("unable get check in record: %v", err)
	}
	return record, nil
}

const CheckInResultModifiersLength = 4

func CheckIn(user models.Account) (models.CheckInRecord, error) {
	tier := rand.Intn(5)
	expMax := 100 * (tier + 1)
	expMin := 100
	record := models.CheckInRecord{
		ResultTier:       tier,
		ResultExperience: rand.Intn(expMax-expMin) + expMin,
		AccountID:        user.ID,
	}

	modifiers := make([]int, CheckInResultModifiersLength)
	for i := 0; i < CheckInResultModifiersLength; i++ {
		modifiers[i] = rand.Intn(1025) // from 0 to 1024 as the comment said
	}
	record.ResultModifiers = modifiers

	if err := CheckCanCheckIn(user); err != nil {
		return record, fmt.Errorf("today already signed")
	}

	tx := database.C.Begin()

	var profile models.AccountProfile
	if err := database.C.Where("account_id = ?", user.ID).First(&profile).Error; err != nil {
		return record, fmt.Errorf("unable get account profile: %v", err)
	} else {
		profile.Experience += uint64(record.ResultExperience)
		if err := tx.Save(&profile).Error; err != nil {
			tx.Rollback()
			return record, fmt.Errorf("unable update account profile: %v", err)
		}
	}

	if err := tx.Save(&record).Error; err != nil {
		return record, fmt.Errorf("unable do check in: %v", err)
	}

	tx.Commit()

	return record, nil
}
