package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationChannel string
type NotificationStatus string
type NotificationType string

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
	UserID    uuid.UUID           `json:"user_id"`
	Type      NotificationType    `json:"type"`
	Channel   NotificationChannel `json:"channel"`
	Status    NotificationStatus  `json:"status"`
	Subject   string              `json:"subject"`
	Body      string              `json:"body"`
	Metadata  map[string]any      `json:"metadata"`
	Error     string              `json:"error,omitempty"`
	SentAt    *time.Time          `json:"sent_at,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
}
