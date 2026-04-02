package zitadel

import (
    "context"
    "fmt"
	"net/http"
    "github.com/zitadel/oidc/v3/pkg/oidc"
    "github.com/zitadel/zitadel-go/v3/pkg/client/management"
    "github.com/zitadel/zitadel-go/v3/pkg/client/middleware"
    "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
)

type Client struct {
    Mgmt *management.Client
    auth *middleware.AuthInterceptor
}

func New(ctx context.Context, domain, api, keyPath string) (*Client, error) {
    scopes := []string{
        oidc.ScopeOpenID,
        "urn:zitadel:iam:org:project:id:zitadel:aud",
    }

    tokenSource := middleware.JWTProfileFromPath(ctx, keyPath)

    auth, err := middleware.NewAuthenticator(domain, tokenSource, scopes...)
    if err != nil {
        return nil, fmt.Errorf("failed to create authenticator: %w", err)
    }

    mgmt, err := management.NewClient(
        ctx,
        domain,
        api,
        scopes,
        zitadel.WithInsecure(),
        zitadel.WithJWTProfileTokenSource(tokenSource),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create management client: %w", err)
    }

    return &Client{
        Mgmt: mgmt,
        auth: auth,
    }, nil
}

func (c *Client) GetServiceToken(ctx context.Context) (string, error) {
    token, err := c.auth.Token()
    if err != nil {
        return "", fmt.Errorf("get service token: %w", err)
    }
    return token.AccessToken, nil
}

func (c *Client) NewHTTPClient(ctx context.Context) (*http.Client, error) {
    token, err := c.GetServiceToken(ctx)
    if err != nil {
        return nil, err
    }

    return &http.Client{
        Transport: &tokenTransport{token: token},
    }, nil
}

type tokenTransport struct {
    token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    req.Header.Set("Authorization", "Bearer "+t.token)
    return http.DefaultTransport.RoundTrip(req)
}

