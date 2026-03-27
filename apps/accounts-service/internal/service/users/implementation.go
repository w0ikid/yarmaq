package users

import (
    "context"
    "fmt"

    "github.com/w0ikid/yarmaq/pkg/zitadel"
    management_pb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
    "go.uber.org/zap"
)

type Service interface {
	GetEmailByID(ctx context.Context, id string) (string, error)
}

type implementation struct {
	client *zitadel.Client
	logger *zap.SugaredLogger
}

func NewService(client *zitadel.Client, logger *zap.SugaredLogger) Service {
	return &implementation{
		client: client,
		logger: logger.Named("users_service"),
	}
}

func (s *implementation) GetEmailByID(ctx context.Context, id string) (string, error) {
	resp, err := s.client.Mgmt.GetUserByID(ctx, &management_pb.GetUserByIDRequest{
		Id: id,
	})
	if err != nil {
		return "", fmt.Errorf("get email by id %s: %w", id, err)
	}

	email := resp.GetUser().GetHuman().GetEmail().GetEmail()
	if email == "" {
		return "", fmt.Errorf("email not found for user: %s", id)
	}

	return email, nil
}
