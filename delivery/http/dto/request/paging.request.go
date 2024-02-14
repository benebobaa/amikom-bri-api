package request

import "time"

const (
	DAILY   = "daily"
	MONTHLY = "monthly"
	YEARLY  = "yearly"
)

type SearchPaginationRequest struct {
	Keyword    string    `json:"keyword"`
	Filter     string    `json:"filter"`
	Date       string    `json:"date"`
	ParsedDate time.Time `json:"parsed_date"`
	ExportPdf  bool      `json:"export_pdf"`
	Page       int       `json:"page" validate:"min=1"`
	Size       int       `json:"size" validate:"min=1,max=100"`
}

type FilterDatePaginationRequest struct {
	Page int `json:"page" validate:"min=1"`
	Size int `json:"size" validate:"min=1,max=100"`
}
