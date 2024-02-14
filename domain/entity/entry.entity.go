package entity

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"gorm.io/gorm"
	"time"
)

const (
	EntryIn  = "IN"
	EntryOut = "OUT"
)

type Entry struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	AccountID int64          `gorm:"column:account_id"`
	Date      time.Time      `gorm:"column:date;type:date; not null; default:now()"`
	Amount    int64          `gorm:"column:amount"`
	EntryType string         `gorm:"column:entry_type"`
	Account   Account        `gorm:"foreignKey:AccountID;references:ID"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (e *Entry) TableName() string {
	return "entries"
}

func (e *Entry) ToEntryResponse() *response.EntryResponse {
	return &response.EntryResponse{
		ID:        e.ID,
		AccountID: e.AccountID,
		Date:      e.Date.String(),
		Amount:    e.Amount,
		EntryType: e.EntryType,
		CreatedAt: e.CreatedAt,
	}
}

func ToEntryResponses(entries []Entry, pagingMetaData *response.PageMetaData) *response.EntryResponses {
	var entryResponses []response.EntryResponse
	for _, entry := range entries {
		entryResponses = append(entryResponses, *entry.ToEntryResponse())
	}
	return &response.EntryResponses{
		Entries: entryResponses,
		Paging:  pagingMetaData,
	}
}

func (e *Entry) ToString() string {
	return fmt.Sprintf("%d ")
}
