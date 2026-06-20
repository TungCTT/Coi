package txmanager

import (
	"context"

	"gorm.io/gorm"
)

// Khai báo một key để giấu tx vào context
type contextKey string

const txKey = contextKey("tx")

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxManager {
	return &txManager{db: db}
}

func (tm *txManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		// Nhét tx vào một cái Context mới
		txCtx := context.WithValue(ctx, txKey, tx)
		// Trả context mới đó cho Service
		return fn(txCtx)
	})
}

func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db // Nếu không có transaction thì xài DB gốc
}
