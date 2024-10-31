package services

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/rs/zerolog/log"
	"time"
)

func DoAutoDatabaseCleanup() {
	log.Debug().Msg("Now cleaning up entire database...")

	var count int64
	for _, model := range database.AutoMaintainRange {
		tx := database.C.Unscoped().Delete(model, "deleted_at IS NOT NULL")
		if tx.Error != nil {
			log.Error().Err(tx.Error).Msg("An error occurred when running cleaning up entire database...")
		}
		count += tx.RowsAffected
	}

	deadline := time.Now().Add(-30 * 24 * time.Hour)
	seenDeadline := time.Now().Add(-7 * 24 * time.Hour)
	database.C.Unscoped().Where("created_at <= ? OR read_at <= ?", deadline, seenDeadline).Delete(&models.Notification{})

	log.Debug().Int64("affected", count).Msg("Clean up entire database accomplished.")
}
