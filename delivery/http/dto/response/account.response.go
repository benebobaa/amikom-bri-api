package response

import "time"

type AccountResponse struct {
	ID        int       `json:"id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
