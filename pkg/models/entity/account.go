package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Account struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	UserID    *string    `gorm:"type:varchar(255);uniqueIndex;column:user_id" json:"user_id"`
	Type      string     `gorm:"type:varchar(20);not null;default:'USER';column:type" json:"type"`
	Number    string     `gorm:"type:varchar(20);uniqueIndex;not null;column:number" json:"number"`
	Balance   int64      `gorm:"type:bigint;not null;default:0;column:balance" json:"balance"`
	Currency  string     `gorm:"type:varchar(3);not null;default:'KZT';column:currency" json:"currency"`
	Status    string     `gorm:"type:varchar(20);not null;default:'ACTIVE';column:status" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (a Account) ToDTO() models.Account {
	return models.Account{
		ID:        a.ID,
		UserID:    a.UserID,
		Type:      a.Type,
		Number:    a.Number,
		Balance:   a.Balance,
		Currency:  a.Currency,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func FromAccountDTO(a models.Account) Account {
	return Account{
		ID:        a.ID,
		UserID:    a.UserID,
		Type:      a.Type,
		Number:    a.Number,
		Balance:   a.Balance,
		Currency:  a.Currency,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
