package models

import (
	"time"

	"github.com/google/uuid"
)

const (
    OperationTypeHold        = "hold"
    OperationTypeHoldRelease = "hold_release"
    OperationTypeDeposit     = "deposit"
    OperationTypeWithdraw    = "withdraw"
    OperationTypeRefund      = "refund"
    OperationTypeTransferIn  = "transfer_in"
    OperationTypeTransferOut = "transfer_out"
)

type Ledger struct {
	ID            uuid.UUID  `json:"id"`
	AccountID     uuid.UUID  `json:"account_id"`
	Amount        int64      `json:"amount"`
	OperationType string     `json:"operation_type"`
	ReferenceID   *uuid.UUID `json:"reference_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
