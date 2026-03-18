package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	ZitadelUserID string         `gorm:"type:varchar(255);uniqueIndex;not null;column:zitadel_user_id" json:"zitadel_user_id"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null;column:email" json:"email"`
	Username      string         `gorm:"type:varchar(100);not null;column:username" json:"username"`
	Roles         pq.StringArray `gorm:"type:text[]" json:"roles"`
	ImageURL      *string        `gorm:"type:text;column:image_url" json:"image_url"`
	IsActive      bool           `gorm:"type:boolean;not null;default:true;column:is_active" json:"is_active"`
	Address       *string        `gorm:"type:varchar(500);column:address" json:"address"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (u User) ToDTO() models.User {
	return models.User{
		ID:            u.ID,
		ZitadelUserID: u.ZitadelUserID,
		Email:         u.Email,
		Username:      u.Username,
		Roles:         []string(u.Roles),
		ImageURL:      u.ImageURL,
		IsActive:      u.IsActive,
		Address:       u.Address,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func (u User) ToEntity() User {
	return User{
		ID:            u.ID,
		ZitadelUserID: u.ZitadelUserID,
		Email:         u.Email,
		Username:      u.Username,
		Roles:         pq.StringArray(u.Roles),
		ImageURL:      u.ImageURL,
		IsActive:      u.IsActive,
		Address:       u.Address,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func FromDTO(u models.User) User {
	return User{
		ZitadelUserID: u.ZitadelUserID,
		Email:         u.Email,
		Username:      u.Username,
		Roles:         pq.StringArray(u.Roles),
		ImageURL:      u.ImageURL,
		Address:       u.Address,
		IsActive:      u.IsActive,
	}
}
