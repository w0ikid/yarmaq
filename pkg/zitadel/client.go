package zitadel

import (
	"context"
	"fmt"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
)

type Client struct {
	Mgmt *management.Client
}

func New(ctx context.Context, domain, api, keyPath string) (*Client, error) {
	scopes := []string{
		oidc.ScopeOpenID,
		"urn:zitadel:iam:org:project:id:zitadel:aud", // стандартный аудиенс для API
		"zitadel.api",
	}

	mgmt, err := management.NewClient(
        ctx,
        domain,  // issuer: "http://zitadel.localhost:8080"
        api,     // gRPC api: "zitadel.localhost:8080"
        scopes,
		zitadel.WithInsecure(),
        zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(ctx, keyPath)),
    )

	if err != nil {
		return nil, fmt.Errorf("failed to create management client: %w", err)
	}

	return &Client{
		Mgmt: mgmt,
	}, nil
}
