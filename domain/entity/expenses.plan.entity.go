package entity

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"gorm.io/gorm"
	"time"
)

type ExpensesPlan struct {
	ID          int64          `gorm:"column:id"`
	UserID      string         `gorm:"column:user_id"`
	Title       string         `gorm:"title"`
	Description string         `gorm:"description"`
	Amount      int64          `gorm:"amount"`
	Date        string         `gorm:"column:date"`
	User        User           `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (e *ExpensesPlan) TableName() string {
	return "expenses_plans"
}

func (e *ExpensesPlan) ToExpensesResponse() *response.ExpensesPlanResponse {
	return &response.ExpensesPlanResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Amount:      e.Amount,
		Date:        e.Date,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func ToExpensesResponses(expensePlans []ExpensesPlan, pagingMetadata *response.PageMetaData) *response.ExpensesPlanResponses {
	var expenseResponses []response.ExpensesPlanResponse
	for _, expense := range expensePlans {
		expenseResponses = append(expenseResponses, *expense.ToExpensesResponse())
	}
	return &response.ExpensesPlanResponses{
		ExpensesPlans: expenseResponses,
		Paging:        pagingMetadata,
	}
}
