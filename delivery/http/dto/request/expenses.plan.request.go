package request

import (
	"github.com/benebobaa/amikom-bri-api/domain/entity"
)

type ExpensesPlanRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,gte=0"`
	Date        string `json:"date" validate:"required,datetime=2006-01-02"`
}

type ExpensesPlanUpdateRequest struct {
	Title       string `json:"title" validate:"min=3,max=100"`
	Description string `json:"description" `
	Amount      int64  `json:"amount" validate:"gte=0,numeric,number"`
	Date        string `json:"date" validate:"datetime=2006-01-02"`
}

func (e *ExpensesPlanRequest) ToEntity(userId string) *entity.ExpensesPlan {
	return &entity.ExpensesPlan{
		UserID:      userId,
		Title:       e.Title,
		Description: e.Description,
		Amount:      e.Amount,
		Date:        e.Date,
	}
}

func (e *ExpensesPlanUpdateRequest) ToEntity(userId string) *entity.ExpensesPlan {
	return &entity.ExpensesPlan{
		Title:       e.Title,
		UserID:      userId,
		Description: e.Description,
		Amount:      e.Amount,
		Date:        e.Date,
	}
}
