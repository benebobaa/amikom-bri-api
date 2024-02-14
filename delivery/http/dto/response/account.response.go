package response

type AccountResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
}
