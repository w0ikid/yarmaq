package ctxkeys

import "context"

type contextKey string

const (
	UserID   contextKey = "userID"
	Roles    contextKey = "roles"
	Claims   contextKey = "claims"
	ClientID contextKey = "clientID"
)

func WithUserContext(ctx context.Context, userID string, roles []string) context.Context {
	ctx = context.WithValue(ctx, UserID, userID)
	ctx = context.WithValue(ctx, Roles, roles)
	return ctx
}

func WithServiceContext(ctx context.Context, clientID string, claims map[string]interface{}) context.Context {
	ctx = context.WithValue(ctx, ClientID, clientID)
	ctx = context.WithValue(ctx, Claims, claims)
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

func GetClientID(ctx context.Context) string {
	if val, ok := ctx.Value(ClientID).(string); ok {
		return val
	}
	return ""
}

func GetClaims(ctx context.Context) map[string]interface{} {
	if val, ok := ctx.Value(Claims).(map[string]interface{}); ok {
		return val
	}
	return nil
}
