package services

import (
	"git.solsynth.dev/hypernet/passport/pkg/authkit/models"
	"git.solsynth.dev/hypernet/passport/pkg/internal/database"
	"github.com/rs/zerolog/log"
)

var (
	writeEventQueue []models.ActionEvent
	writeAuditQueue []models.AuditRecord
)

// AddEvent to keep operation logs by user themselves clear to query
func AddEvent(user uint, event string, meta map[string]any, ip, ua string) {
	writeEventQueue = append(writeEventQueue, models.ActionEvent{
		Type:      event,
		Metadata:  meta,
		IpAddress: ip,
		UserAgent: ua,
		AccountID: user,
	})
}

// AddAuditRecord to keep logs to make administrators' operations clear to query
func AddAuditRecord(operator models.Account, act, ip, ua string, metadata map[string]any) {
	writeAuditQueue = append(writeAuditQueue, models.AuditRecord{
		Action:    act,
		Metadata:  metadata,
		IpAddress: ip,
		UserAgent: ua,
		AccountID: operator.ID,
	})
}

// SaveEventChanges runs every 60 seconds to save events / audits changes into the database
func SaveEventChanges() {
	if len(writeEventQueue) > 0 {
		count := len(writeEventQueue)
		database.C.CreateInBatches(writeEventQueue, min(count, 1000))
		log.Info().Int("count", count).Msg("Saved action events changes into database...")
		writeEventQueue = nil
	}
	if len(writeAuditQueue) > 0 {
		count := len(writeAuditQueue)
		database.C.CreateInBatches(writeAuditQueue, min(count, 1000))
		log.Info().Int("count", count).Msg("Saved audit records changes into database...")
		writeAuditQueue = nil
	}
}
