package models

const (
	ReportStatusPending   = "pending"
	ReportStatusReviewing = "reviewing"
	ReportStatusConfirmed = "confirmed"
	ReportStatusRejected  = "rejected"
	ReportStatusProcessed = "processed"
)

type AbuseReport struct {
	BaseModel

	Resource  string  `json:"resource"`
	Reason    string  `json:"reason"`
	Status    string  `json:"status"`
	AccountID uint    `json:"account_id"`
	Account   Account `json:"account"`
}
