package response

import "time"

type TranferInfo struct {
	Amount        int64     `json:"amount"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferResponse struct {
	Transfer    *TranferInfo     `json:"transfer"`
	FromAccount *AccountResponse `json:"from_account"`
	ToAccount   *AccountResponse `json:"to_account"`
}
