package models

type ApiKey struct {
	BaseModel

	Name        string     `json:"name"`
	Description string     `json:"description"`
	Lifecycle   *int64     `json:"lifecycle"`
	Ticket      AuthTicket `json:"ticket" gorm:"TicketID"`
	TicketID    uint       `json:"ticket_id"`
	Account     Account    `json:"account"`
	AccountID   uint       `json:"account_id"`
}
