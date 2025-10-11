package sleeper

type transactionType string

const (
	TransactionTypeTrade     transactionType = "trade"
	TransactionTypeWaiver    transactionType = "waiver"
	TransactionTypeFreeAgent transactionType = "free_agent"
)

// Transaction represents a transaction in a league
type Transaction struct {
	Type          string               `json:"type,omitempty"`
	TransactionID string               `json:"transaction_id,omitempty"`
	StatusUpdated int64                `json:"status_updated,omitempty"`
	Status        string               `json:"status,omitempty"`
	Settings      *TransactionSettings `json:"settings,omitempty"`   // trades do not use this field
	RosterIDs     []int                `json:"roster_ids,omitempty"` // roster_ids involved in this transaction
	Metadata      *TransactionMetadata `json:"metadata,omitempty"`
	Leg           int                  `json:"leg,omitempty"` // in football, this is the week
	Drops         map[string]int       `json:"drops,omitempty"`
	DraftPicks    []*TradedDraftPick   `json:"draft_picks,omitempty"` // picks that were traded
	Creator       string               `json:"creator,omitempty"`     // user id who initiated the transaction
	Created       int64                `json:"created,omitempty"`
	ConsenterIDs  []int                `json:"consenter_ids,omitempty"` // roster_ids of the people who agreed to this transaction
	Adds          map[string]int       `json:"adds,omitempty"`
	WaiverBudget  []*WaiverBudget      `json:"waiver_budget,omitempty"` // roster_id 2 sends 55 FAAB dollars to roster_id 3
}

// TransactionSettings holds settings for waiver transactions
type TransactionSettings struct {
	WaiverBid int `json:"waiver_bid,omitempty"`
}

// TransactionMetadata can contain notes about why a transaction didn't go through
type TransactionMetadata struct {
	Notes string `json:"notes,omitempty"`
}

// TradedDraftPick represents a draft pick that was traded
type TradedDraftPick struct {
	Season          string `json:"season,omitempty"`            // the season this draft pick belongs to
	Round           int    `json:"round,omitempty"`             // which round this draft pick is
	RosterID        int    `json:"roster_id,omitempty"`         // original owner's roster_id
	PreviousOwnerID int    `json:"previous_owner_id,omitempty"` // previous owner's roster id (in this trade)
	OwnerID         int    `json:"owner_id,omitempty"`          // the new owner of this pick after the trade
}

// WaiverBudget represents a waiver amount involved in a trade
type WaiverBudget struct {
	Sender   int `json:"sender,omitempty"`
	Receiver int `json:"receiver,omitempty"`
	Amount   int `json:"amount,omitempty"`
}
