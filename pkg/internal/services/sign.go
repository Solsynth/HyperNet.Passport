package services

import (
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"time"
)

func CheckDailyCanSign(user uint) error {
	probe := time.Now().Format("2006-01-02")

	var record models.SignRecord
	if err := database.C.Where("account_id = ? AND created_at::date = ?", user, probe).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("unable check daliy sign record: %v", err)
	}
	return fmt.Errorf("daliy sign record exists")
}

func GetTodayDailySign(user uint) (models.SignRecord, error) {
	probe := time.Now().Format("2006-01-02")

	var record models.SignRecord
	if err := database.C.Where("account_id = ? AND created_at::date = ?", user, probe).First(&record).Error; err != nil {
		return record, fmt.Errorf("unable get daliy sign record: %v", err)
	}
	return record, nil
}

func DailySign(user uint) (models.SignRecord, error) {
	tier := rand.Intn(5)
	record := models.SignRecord{
		ResultTier:       tier,
		ResultExperience: rand.Intn(int(math.Max(float64(tier), 1)*100)+1-100) + 100,
		AccountID:        user,
	}

	if err := CheckDailyCanSign(user); err != nil {
		return record, fmt.Errorf("today already signed")
	}

	var profile models.AccountProfile
	if err := database.C.Where("account_id = ?", user).First(&profile).Error; err != nil {
		return record, fmt.Errorf("unable get account profile: %v", err)
	} else {
		profile.Experience += uint64(record.ResultExperience)
		database.C.Save(&profile)
	}

	if err := database.C.Save(&record).Error; err != nil {
		return record, fmt.Errorf("unable do daliy sign: %v", err)
	}

	return record, nil
}
