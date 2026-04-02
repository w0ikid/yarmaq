package users

import (
    "context"
    "fmt"

    "github.com/w0ikid/yarmaq/pkg/zitadel"
    management_pb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
    "go.uber.org/zap"
	"encoding/base64"
	"strings"
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
    // --- Отладка токена перед RPC ---
    if s.client != nil {
        token, err := s.client.GetServiceToken(ctx)
        if err != nil {
            s.logger.Errorw("failed to get service token", "error", err)
        } else {
            shortToken := token
            if len(token) > 20 {
                shortToken = token[:20] + "..."
            }
            s.logger.Infow("service token (shortened)", "token", shortToken)

            // Декодируем payload JWT
            parts := strings.Split(token, ".")
            if len(parts) == 3 {
                payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
                if err != nil {
                    s.logger.Errorw("failed to decode JWT payload", "error", err)
                } else {
                    s.logger.Infow("JWT payload", "payload", string(payloadBytes))
                }
            } else {
                s.logger.Warnw("token is not valid JWT format")
            }
        }
    }
    // -------------------------------

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
