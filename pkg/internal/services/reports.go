package services

import (
	"fmt"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"
)

func ListAbuseReport(account models.Account) ([]models.AbuseReport, error) {
	var reports []models.AbuseReport
	err := database.C.
		Where("account_id = ?", account.ID).
		Find(&reports).Error
	return reports, err
}

func GetAbuseReport(id uint) (models.AbuseReport, error) {
	var report models.AbuseReport
	err := database.C.
		Where("id = ?", id).
		First(&report).Error
	return report, err
}

func UpdateAbuseReportStatus(id uint, status, message string) error {
	var report models.AbuseReport
	err := database.C.
		Where("id = ?", id).
		Preload("Account").
		First(&report).Error
	if err != nil {
		return err
	}

	report.Status = status
	account := report.Account

	err = database.C.Save(&report).Error
	if err != nil {
		return err
	}

	NewNotification(models.Notification{
		Topic:     "reports.feedback",
		Title:     "Abuse report status has been changed.",
		Subtitle:  lo.ToPtr(fmt.Sprintf("The report #%d's status updated", id)),
		Body:      fmt.Sprintf("The report created by you with ID #%d's status has been changed to %s. Moderator message: %s", id, status, message),
		Account:   account,
		AccountID: account.ID,
	})

	return nil
}

func NewAbuseReport(resource string, reason string, account models.Account) (models.AbuseReport, error) {
	var report models.AbuseReport
	if err := database.C.
		Where(
			"resource = ? AND account_id = ? AND status IN ?",
			resource,
			account.ID,
			[]string{models.ReportStatusPending, models.ReportStatusReviewing},
		).First(&report).Error; err == nil {
		return report, fmt.Errorf("you already reported this resource and it still in process")
	}

	report = models.AbuseReport{
		Resource:  resource,
		Reason:    reason,
		AccountID: account.ID,
	}

	err := database.C.Create(&report).Error
	return report, err
}