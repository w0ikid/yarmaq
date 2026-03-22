package accounts

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/google/uuid"
    "github.com/w0ikid/yarmaq/pkg/httpclient"
    "github.com/w0ikid/yarmaq/pkg/models"
)

type Client struct {
    base *httpclient.Client
}

func New(baseURL string, base *httpclient.Client) *Client {
    return &Client{base: base}
}

func (c *Client) GetAccount(ctx context.Context, id string) (*models.AccountResponse, error) {
    req, err := http.NewRequestWithContext(ctx, "GET",
        fmt.Sprintf("%s/api/v1/accounts/%s", c.base.BaseURL, id), nil)
    if err != nil {
        return nil, err
    }
    resp, err := c.base.Do(ctx, req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("accounts-service returned %d", resp.StatusCode)
    }

    var account models.AccountResponse
    if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
        return nil, err
    }
    return &account, nil
}

func (c *Client) UpdateBalance(ctx context.Context, id string, req models.UpdateBalanceRequest) error {
    body, err := json.Marshal(req)
    if err != nil {
        return err
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST",
        fmt.Sprintf("%s/api/v1/internal/accounts/%s/balance", c.base.BaseURL, id),
        bytes.NewReader(body))
    if err != nil {
        return err
    }

    resp, err := c.base.Do(ctx, httpReq)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        var respBody bytes.Buffer
        _, _ = respBody.ReadFrom(resp.Body)
        return fmt.Errorf("accounts-service returned %d: %s", resp.StatusCode, respBody.String())
    }
    return nil
}

func (c *Client) Hold(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
    return c.UpdateBalance(ctx, accountID, models.UpdateBalanceRequest{
        Amount:        -amount,
        OperationType: models.OperationTypeHold,
        ReferenceID:   &transactionID,
    })
}

func (c *Client) Deposit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
    return c.UpdateBalance(ctx, accountID, models.UpdateBalanceRequest{
        Amount:        amount,
        OperationType: models.OperationTypeDeposit,
        ReferenceID:   &transactionID,
    })
}

func (c *Client) Refund(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
    return c.UpdateBalance(ctx, accountID, models.UpdateBalanceRequest{
        Amount:        amount,
        OperationType: models.OperationTypeRefund,
        ReferenceID:   &transactionID,
    })
}