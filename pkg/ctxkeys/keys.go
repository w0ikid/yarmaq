package ctxkeys

import "context"

type contextKey string

const (
    UserID contextKey = "userID"
    Roles  contextKey = "roles"
)

func WithUserContext(ctx context.Context, userID string, roles []string) context.Context {
    ctx = context.WithValue(ctx, UserID, userID)
    ctx = context.WithValue(ctx, Roles, roles)
    return ctx
}

func GetUserID(ctx context.Context) string {
    return ctx.Value(UserID).(string)
}

func GetRoles(ctx context.Context) []string {
    return ctx.Value(Roles).([]string)
}