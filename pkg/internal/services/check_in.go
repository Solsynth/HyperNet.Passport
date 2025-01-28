package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/wallet/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"gorm.io/gorm"
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
	var record models.CheckInRecord
	if err := CheckCanCheckIn(user); err != nil {
		return record, fmt.Errorf("today already signed")
	}

	tier := rand.Intn(5)

	expMax := 100 * (tier + 1)
	expMin := 100
	exp := rand.Intn(expMax+1-expMin) + expMin

	coinMax := 10.0 * float64(tier+1)
	coinMin := 10.0
	rawCoins := coinMax + rand.Float64()*(coinMax-coinMin)

	record = models.CheckInRecord{
		ResultTier:       tier,
		ResultExperience: exp,
		ResultCoin:       float64(int(rawCoins*100)) / 100,
		AccountID:        user.ID,
	}

	modifiers := make([]int, CheckInResultModifiersLength)
	for i := 0; i < CheckInResultModifiersLength; i++ {
		modifiers[i] = rand.Intn(1025) // from 0 to 1024 as the comment said
	}
	record.ResultModifiers = modifiers

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

	conn, err := gap.Nx.GetClientGrpcConn("wa")
	if err != nil {
		log.Warn().Err(err).Msg("Unable to connect with wallet to send daily rewards")
		record.ResultCoin = 0
	}
	wc := proto.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err = wc.MakeTransactionWithAccount(ctx, &proto.MakeTransactionWithAccountRequest{
		PayeeAccountId: lo.ToPtr(uint64(user.ID)),
		Amount:         record.ResultCoin,
		Remark:         "Daily Check-In Rewards",
	})
	if err != nil {
		log.Warn().Err(err).Msg("Unable to make transaction with account to send daily rewards")
		record.ResultCoin = 0
	}

	if err := tx.Save(&record).Error; err != nil {
		return record, fmt.Errorf("unable do check in: %v", err)
	}

	tx.Commit()

	return record, nil
}
