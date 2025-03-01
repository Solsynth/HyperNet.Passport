package services

import (
	"time"

	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/rs/zerolog/log"
)

func DoAutoDatabaseCleanup() {
	log.Debug().Msg("Now cleaning up entire database...")

	var count int64

	deadline := time.Now().Add(-30 * 24 * time.Hour)
	seenDeadline := time.Now().Add(-7 * 24 * time.Hour)
	tx := database.C.Unscoped().Where("created_at <= ? OR read_at <= ?", deadline, seenDeadline).Delete(&models.Notification{})
	count += tx.RowsAffected

	log.Debug().Int64("affected", count).Msg("Clean up entire database accomplished.")
}
