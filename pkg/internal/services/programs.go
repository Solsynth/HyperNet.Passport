package services

import (
	"context"
	"fmt"
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"git.solsynth.dev/hypernet/passport/pkg/internal/gap"
	"git.solsynth.dev/hypernet/wallet/pkg/proto"
	"github.com/samber/lo"
)

func JoinProgram(user models.Account, program models.Program) (models.ProgramMember, error) {
	var member models.ProgramMember
	if err := database.C.Where("account_id = ? AND program_id = ?", user.ID, program.ID).First(&member).Error; err == nil {
		return member, fmt.Errorf("program member already exists")
	}
	var profile models.AccountProfile
	if err := database.C.Where("account_id = ?", user.ID).Select("Experience").First(&profile).Error; err != nil {
		return member, err
	}
	if program.ExpRequirement < int64(profile.Experience) {
		return member, fmt.Errorf("insufficient experience")
	}
	member = models.ProgramMember{
		LastPaid:  lo.ToPtr(time.Now()),
		Account:   user,
		AccountID: user.ID,
		Program:   program,
		ProgramID: program.ID,
	}
	if err := ChargeForProgram(member); err != nil {
		return member, err
	}
	if err := database.C.Create(&member).Error; err != nil {
		return member, err
	}
	return member, nil
}

func LeaveProgram(user models.Account, program models.Program) error {
	var member models.ProgramMember
	if err := database.C.Where("account_id = ? AND program_id = ?", user.ID, program.ID).First(&member).Error; err != nil {
		return err
	}
	if err := database.C.Delete(&member).Error; err != nil {
		return err
	}
	return nil
}

func ChargeForProgram(member models.ProgramMember) error {
	pricing := member.Program.Price.Data()
	if pricing.Amount == 0 {
		return nil
	}
	conn, err := gap.Nx.GetClientGrpcConn("wa")
	if err != nil {
		return err
	}
	wc := proto.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err = wc.MakeTransactionWithAccount(ctx, &proto.MakeTransactionWithAccountRequest{
		PayeeAccountId: lo.ToPtr(uint64(member.AccountID)),
		Amount:         pricing.Amount,
		Currency:       pricing.Currency,
		Remark:         fmt.Sprintf("Program Membership: %s", member.Program.Name),
	})
	return err
}

func PeriodicChargeProgramFee() {
	var members []models.ProgramMember
	if err := database.C.Preload("Program").Find(&members).Error; err != nil {
		return
	}
	for _, member := range members {
		// every month paid once
		if member.LastPaid == nil || time.Since(*member.LastPaid) < time.Hour*24*30 {
			if err := ChargeForProgram(member); err == nil {
				database.C.Model(&member).Update("last_paid", time.Now())
			}
		}
	}
}
