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
	if val, ok := ctx.Value(UserID).(string); ok {
		return val
	}
	return ""
}

func GetRoles(ctx context.Context) []string {
	if val, ok := ctx.Value(Roles).([]string); ok {
		return val
	}
	return nil
}