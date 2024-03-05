package request

import (
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/google/uuid"
	"time"
)

type SessionRequest struct {
	ID           string    `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	IsBlocked    bool      `json:"is_blocked"`
	ClientIP     string    `json:"client_ip"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *SessionRequest) ToEntity() *entity.Session {
	return &entity.Session{
		ID:           r.ID,
		UserID:       r.UserID,
		RefreshToken: r.RefreshToken,
		UserAgent:    r.UserAgent,
		ClientIP:     r.ClientIP,
		ExpiredAt:    r.ExpiredAt,
	}
}
