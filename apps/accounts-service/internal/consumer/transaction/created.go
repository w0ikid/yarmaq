package transaction

import (
    "context"
    "encoding/json"

    "github.com/segmentio/kafka-go"
    "go.uber.org/zap"

    "github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
)

type CreatedHandlerDeps struct {
    AccountUsecase account.AccountDomain
    Logger         *zap.SugaredLogger
}

type CreatedHandler struct {
    deps CreatedHandlerDeps
}

func NewCreatedHandler(deps CreatedHandlerDeps) *CreatedHandler {
    return &CreatedHandler{deps: deps}
}

func (h *CreatedHandler) Handle(ctx context.Context, msg kafka.Message) error {
    var event struct {
        ID     string `json:"id"`
        UserID string `json:"user_id"`
    }
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return err
    }

    h.deps.Logger.Infow("account.created received", "id", event.ID)
    // бизнес логика
    return nil
}