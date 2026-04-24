package middleware

import (
	"context"

	"github.com/w0ikid/yarmaq/pkg/auth/jwks"
	"github.com/w0ikid/yarmaq/pkg/core/ctxkeys"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCAuthInterceptor validates the JWT token in gRPC metadata and populates context with claims.
func GRPCAuthInterceptor(j *jwks.JWKS) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		claims, err := j.Validate(values[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		clientID, _ := claims["client_id"].(string)
		if clientID == "" {
			clientID, _ = claims["azp"].(string)
		}

		// Update context with claims and clientID
		ctx = ctxkeys.WithServiceContext(ctx, clientID, claims)

		// Also try to set UserID and Roles if it's a user token
		if sub, ok := claims["sub"].(string); ok {
			var roles []string
			if rolesRaw, ok := claims["urn:zitadel:iam:org:project:roles"].(map[string]interface{}); ok {
				for role := range rolesRaw {
					roles = append(roles, role)
				}
			}
			ctx = ctxkeys.WithUserContext(ctx, sub, roles)
		}

		return handler(ctx, req)
	}
}

// GRPCServiceOnlyInterceptor restricts access to specific services based on clientID in claims.
func GRPCServiceOnlyInterceptor(logger *zap.SugaredLogger, allowedServices ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		clientID := ctxkeys.GetClientID(ctx)
		if clientID == "" {
			return nil, status.Errorf(codes.PermissionDenied, "service token required")
		}

		allowed := false
		for _, svc := range allowedServices {
			if svc == clientID {
				allowed = true
				break
			}
		}

		if !allowed {
			logger.Warnw("unauthorized gRPC service call",
				"client_id", clientID,
				"method", info.FullMethod,
			)
			return nil, status.Errorf(codes.PermissionDenied, "service %s is not allowed to call this method", clientID)
		}

		logger.Infow("internal gRPC service call", "from", clientID, "method", info.FullMethod)
		return handler(ctx, req)
	}
}
