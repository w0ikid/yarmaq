package consumers

import (
	"context"
	"fmt"

	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	managementpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
)

type UserEmailResolver interface {
	ResolveEmail(ctx context.Context, userID string) (string, error)
}

type AccountUserResolver interface {
	ResolveUserID(ctx context.Context, accountID string) (string, error)
}

type Resolver struct {
	accountsClient *accounts.Client
	zitadelClient  *zitadel.Client
}

func NewResolver(accountsClient *accounts.Client, zitadelClient *zitadel.Client) *Resolver {
	return &Resolver{
		accountsClient: accountsClient,
		zitadelClient:  zitadelClient,
	}
}

func (r *Resolver) ResolveUserID(ctx context.Context, accountID string) (string, error) {
	account, err := r.accountsClient.GetAccount(ctx, accountID)
	if err != nil {
		return "", fmt.Errorf("get account %s: %w", accountID, err)
	}
	if account == nil {
		return "", nil
	}

	return account.UserID, nil
}

func (r *Resolver) ResolveEmail(ctx context.Context, userID string) (string, error) {
	resp, err := r.zitadelClient.Mgmt.GetHumanEmail(ctx, &managementpb.GetHumanEmailRequest{
		UserId: userID,
	})
	if err != nil {
		return "", fmt.Errorf("get human email for user %s: %w", userID, err)
	}
	if resp == nil || resp.GetEmail() == nil || resp.GetEmail().GetEmail() == "" {
		return "", fmt.Errorf("email not found for user %s", userID)
	}

	return resp.GetEmail().GetEmail(), nil
}
