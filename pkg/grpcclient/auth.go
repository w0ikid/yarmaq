package grpcclient

import (
	"context"
	"fmt"

	"github.com/w0ikid/yarmaq/pkg/zitadel"
)

type zitadelCreds struct {
	zitadel *zitadel.Client
}

func (c *zitadelCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := c.zitadel.GetServiceToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("zitadel token: %w", err)
	}

	return map[string]string{
		"authorization": "Bearer " + token,
	}, nil
}

func (c *zitadelCreds) RequireTransportSecurity() bool {
	return false // Set to true if using TLS
}

func NewZitadelCreds(client *zitadel.Client) *zitadelCreds {
	return &zitadelCreds{zitadel: client}
}
