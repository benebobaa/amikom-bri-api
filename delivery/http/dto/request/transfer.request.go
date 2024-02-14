package request

import "github.com/benebobaa/amikom-bri-api/domain/entity"

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,numeric"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,numeric,nefield=FromAccountID"`
	Amount        int64  `json:"amount" validate:"required,gt=0,numeric"`
	Pin           string `json:"pin" validate:"required,len=6,number"`
}

func (t *TransferRequest) ToEntity() *entity.Transfer {
	return &entity.Transfer{
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
	}
}
