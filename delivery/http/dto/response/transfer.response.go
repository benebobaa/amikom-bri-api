package response

type TranferResponse struct {
	Amount        int64           `json:"amount"`
	FromAccountID AccountResponse `json:"from_account_id"`
	ToAccountID   AccountResponse `json:"to_account_id"`
}
