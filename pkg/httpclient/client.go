package httpclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/w0ikid/yarmaq/pkg/zitadel"
)

type Client struct {
	http    *http.Client
	zitadel *zitadel.Client
	BaseURL string
}

func New(baseURL string, zitadelClient *zitadel.Client) *Client {
	return &Client{
		http:    &http.Client{},
		zitadel: zitadelClient,
		BaseURL: baseURL,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	token, err := c.zitadel.GetServiceToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get service token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}
