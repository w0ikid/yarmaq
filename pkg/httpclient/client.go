package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"log"
	"time"
	"bytes"

	"github.com/w0ikid/yarmaq/pkg/zitadel"
)

type Client struct {
	http    *http.Client
	zitadel *zitadel.Client
	BaseURL string
}

func New(baseURL string, zitadelClient *zitadel.Client) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
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

	// ---- Логирование ----
	log.Printf("[httpclient] outgoing request: %s %s token=%s...", req.Method, req.URL, token[:8]) // первые 8 символов токена
	// ---------------------

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	// ---- Логирование ответа ----
	if resp.StatusCode != http.StatusOK {
		var respBody bytes.Buffer
		_, _ = respBody.ReadFrom(resp.Body)
		log.Printf("[httpclient] response code=%d body=%s", resp.StatusCode, respBody.String())
		// возвращаем ошибку дальше
		return resp, fmt.Errorf("accounts-service returned %d", resp.StatusCode)
	}

	return resp, nil
}
