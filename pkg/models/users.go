package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id"`
	ZitadelUserID string     `json:"zitadel_user_id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	Roles         []string   `json:"roles"`
	ImageURL      *string    `json:"image_url"`
	Address       *string    `json:"address,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
