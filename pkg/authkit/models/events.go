package models

type ActionEvent struct {
	BaseModel

	Type      string `json:"type"`
	Target    string `json:"target"`
	Location  string `json:"location"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}
