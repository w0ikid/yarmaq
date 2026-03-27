package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/errs"
	"github.com/w0ikid/yarmaq/pkg/models"
	"github.com/w0ikid/yarmaq/pkg/smtpclient"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, notification models.Notification) (*models.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Notification, error)
	Send(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type implementation struct {
	repo       NotificationRepo
	smtpClient *smtpclient.Client
	logger     *zap.SugaredLogger
}

func NewService(repo NotificationRepo, smtpClient *smtpclient.Client, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:       repo,
		smtpClient: smtpClient,
		logger:     logger.Named("notification_service"),
	}
}

func (s *implementation) Create(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	if notification.Channel == "" {
		notification.Channel = models.ChannelEmail
	}
	if notification.Status == "" {
		notification.Status = models.StatusPending
	}

	if err := validateNotification(notification); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, notification)
}

func (s *implementation) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *implementation) GetByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *implementation) Send(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error) {
	notification, err := s.repo.GetByID(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, fmt.Errorf("%w: notification not found: %s", errs.ErrNotFound, notificationID)
	}

	if notification.Channel != models.ChannelEmail {
		return nil, fmt.Errorf("%w: unsupported notification channel: %s", errs.ErrValidation, notification.Channel)
	}

	recipient, err := extractRecipient(notification.Metadata)
	if err != nil {
		return nil, err
	}

	err = s.smtpClient.Send(smtpclient.Message{
		To:      []string{recipient},
		Subject: notification.Subject,
		Body:    notification.Body,
		IsHTML:  true,
	})
	if err != nil {
		s.logger.Errorw("failed to send notification", "id", notificationID, "error", err)
		notification.Status = models.StatusFailed
		notification.Error = err.Error()
		notification.SentAt = nil

		updated, updateErr := s.repo.Update(ctx, *notification)
		if updateErr != nil {
			return nil, updateErr
		}

		return updated, err
	}

	now := time.Now().UTC()
	notification.Status = models.StatusSent
	notification.Error = ""
	notification.SentAt = &now

	return s.repo.Update(ctx, *notification)
}

func (s *implementation) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func validateNotification(notification models.Notification) error {
	if notification.UserID == "" {
		return fmt.Errorf("%w: user_id is required", errs.ErrValidation)
	}
	if notification.Type == "" {
		return fmt.Errorf("%w: type is required", errs.ErrValidation)
	}
	if notification.Channel != models.ChannelEmail {
		return fmt.Errorf("%w: unsupported notification channel: %s", errs.ErrValidation, notification.Channel)
	}
	if notification.Subject == "" {
		return fmt.Errorf("%w: subject is required", errs.ErrValidation)
	}
	if notification.Body == "" {
		return fmt.Errorf("%w: body is required", errs.ErrValidation)
	}

	_, err := extractRecipient(notification.Metadata)
	return err
}

func extractRecipient(metadata map[string]any) (string, error) {
	if metadata == nil {
		return "", fmt.Errorf("%w: notification recipient is required in metadata", errs.ErrValidation)
	}

	for _, key := range []string{"email", "to", "recipient_email"} {
		value, ok := metadata[key]
		if !ok {
			continue
		}

		recipient, ok := value.(string)
		if ok && recipient != "" {
			return recipient, nil
		}
	}

	return "", fmt.Errorf("%w: notification recipient is required in metadata", errs.ErrValidation)
}
