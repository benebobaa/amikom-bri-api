package response

import (
	"time"
)

type AccountResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
