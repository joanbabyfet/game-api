package repository

import (
	"context"

	"game-api/internal/model"

	"gorm.io/gorm"
)

type WalletLogRepository struct {
	db *gorm.DB
}

func NewWalletLogRepository(db *gorm.DB) *WalletLogRepository {
	return &WalletLogRepository{
		db: db,
	}
}

func (r *WalletLogRepository) WithTx(tx *gorm.DB) *WalletLogRepository {
	return &WalletLogRepository{
		db: tx,
	}
}

// Create 写入钱包流水
func (r *WalletLogRepository) Create(
	ctx context.Context,
	log *model.WalletLog,
) error {
	return r.db.
		WithContext(ctx).
		Create(log).
		Error
}