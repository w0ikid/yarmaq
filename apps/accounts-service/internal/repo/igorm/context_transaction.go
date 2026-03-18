package igorm

import (
    "context"
    "gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "tx"

type ContextTransaction struct {
    db *gorm.DB
}

func NewContextTransaction(db *gorm.DB) *ContextTransaction {
    return &ContextTransaction{db: db}
}

func (ct *ContextTransaction) StartTransaction(ctx context.Context) (context.Context, error) {
    if _, ok := ctx.Value(txKey).(*gorm.DB); ok {
        return ctx, nil
    }
    tx := ct.db.Begin()
    if tx.Error != nil {
        return ctx, tx.Error
    }
    return context.WithValue(ctx, txKey, tx), nil
}

func (ct *ContextTransaction) FinalizeTransaction(ctx context.Context, err *error) error {
    tx := RetrieveTx(ctx)
    if tx == nil {
        return nil
    }
    if err != nil && *err != nil {
        if rbErr := tx.Rollback().Error; rbErr != nil {
            return rbErr
        }
        return nil
    }
    return tx.Commit().Error
}

func RetrieveTx(ctx context.Context) *gorm.DB {
    tx, _ := ctx.Value(txKey).(*gorm.DB)
    return tx
}