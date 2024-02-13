package response

import "time"

type AccountResponse struct {
	ID        int64     `json:"id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
