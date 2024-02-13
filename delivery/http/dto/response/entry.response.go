package response

import "time"

type EntryResponse struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Date      string    `json:"date"`
	Amount    int64     `json:"amount"`
	EntryType string    `json:"entry_type"`
	CreatedAt time.Time `json:"created_at"`
}

type EntryResponses struct {
	Entries []EntryResponse `json:"entries"`
	Paging  *PageMetaData   `json:"paging"`
}
