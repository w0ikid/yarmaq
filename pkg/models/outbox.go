package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Outbox struct {
	ID          uuid.UUID       `json:"id"`
	EventType   string          `json:"event_type"`
	Payload     json.RawMessage `json:"payload"`
	AggregateID uuid.UUID       `json:"aggregate_id"`
	CreatedAt   time.Time       `json:"created_at"`
	SentAt      *time.Time      `json:"sent_at,omitempty"`
}
