package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Outbox struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	EventType   string          `gorm:"type:varchar(50);not null;column:event_type" json:"event_type"`
	Payload     json.RawMessage `gorm:"type:jsonb;not null;column:payload" json:"payload"`
	AggregateID uuid.UUID       `gorm:"type:uuid;not null;column:aggregate_id" json:"aggregate_id"`
	CreatedAt   time.Time       `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	SentAt      *time.Time      `gorm:"column:sent_at" json:"sent_at,omitempty"`
}

func (o Outbox) ToDTO() models.Outbox {
	return models.Outbox{
		ID:          o.ID,
		EventType:   o.EventType,
		Payload:     o.Payload,
		AggregateID: o.AggregateID,
		CreatedAt:   o.CreatedAt,
		SentAt:      o.SentAt,
	}
}

func FromOutboxDTO(o models.Outbox) Outbox {
	return Outbox{
		ID:          o.ID,
		EventType:   o.EventType,
		Payload:     o.Payload,
		AggregateID: o.AggregateID,
		CreatedAt:   o.CreatedAt,
		SentAt:      o.SentAt,
	}
}

func (Outbox) TableName() string {
    return "outbox"
}