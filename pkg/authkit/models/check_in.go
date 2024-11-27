package models

type CheckInRecord struct {
	BaseModel

	ResultTier       int     `json:"result_tier"`
	ResultExperience int     `json:"result_experience"`
	Account          Account `json:"account"`
	AccountID        uint    `json:"account_id"`
}
