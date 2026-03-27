package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Notification struct {
	ID        uuid.UUID                  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	UserID    string                     `gorm:"type:varchar(255);not null;column:user_id" json:"user_id"`
	Type      models.NotificationType    `gorm:"type:varchar(100);not null;column:type" json:"type"`
	Channel   models.NotificationChannel `gorm:"type:varchar(20);not null;column:channel" json:"channel"`
	Status    models.NotificationStatus  `gorm:"type:varchar(20);not null;column:status" json:"status"`
	Subject   string                     `gorm:"type:text;not null;column:subject" json:"subject"`
	Body      string                     `gorm:"type:text;not null;column:body" json:"body"`
	Metadata  models.Metadata            `gorm:"type:jsonb;column:metadata" json:"metadata"`
	Error     string                     `gorm:"type:text;column:error" json:"error,omitempty"`
	SentAt    *time.Time                 `gorm:"column:sent_at" json:"sent_at,omitempty"`
	CreatedAt time.Time                  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (n Notification) ToDTO() models.Notification {
	return models.Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Channel:   n.Channel,
		Status:    n.Status,
		Subject:   n.Subject,
		Body:      n.Body,
		Metadata:  n.Metadata,
		Error:     n.Error,
		SentAt:    n.SentAt,
		CreatedAt: n.CreatedAt,
	}
}

func FromNotificationDTO(n models.Notification) Notification {
	return Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Channel:   n.Channel,
		Status:    n.Status,
		Subject:   n.Subject,
		Body:      n.Body,
		Metadata:  n.Metadata,
		Error:     n.Error,
		SentAt:    n.SentAt,
		CreatedAt: n.CreatedAt,
	}
}

func (Notification) TableName() string {
	return "notifications"
}
