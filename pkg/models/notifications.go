package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationChannel string
type NotificationStatus string
type NotificationType string

type Metadata map[string]any

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Metadata) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Metadata: %v", value)
	}
	return json.Unmarshal(bytes, m)
}

const (
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"
	ChannelSMS   NotificationChannel = "sms"
)

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

const (
	TypeTransactionCreated NotificationType = "transaction.created"
	TypeAccountCreated     NotificationType = "account.created"
	TypeLowBalance         NotificationType = "account.low_balance"
)

type Notification struct {
	ID        uuid.UUID           `json:"id"`
	UserID    string              `json:"user_id"`
	Type      NotificationType    `json:"type"`
	Channel   NotificationChannel `json:"channel"`
	Status    NotificationStatus  `json:"status"`
	Subject   string              `json:"subject"`
	Body      string              `json:"body"`
	Metadata  Metadata            `json:"metadata"`
	Error     string              `json:"error,omitempty"`
	SentAt    *time.Time          `json:"sent_at,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
}
