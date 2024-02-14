package response

import "time"

type ExpensesPlanResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpensesPlanResponses struct {
	ExpensesPlans []ExpensesPlanResponse `json:"expenses_plans"`
	Paging        *PageMetaData          `json:"paging"`
}
