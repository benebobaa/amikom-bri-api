package request

type SearchPaginationRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page" validate:"min=1"`
	Size    int    `json:"size" validate:"min=1,max=100"`
}
