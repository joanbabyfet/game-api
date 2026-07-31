package repository

import (
	"context"

	"game-api/internal/model"

	"gorm.io/gorm"
)

type RollbackLogRepository struct {
	db *gorm.DB
}

func NewRollbackLogRepository(db *gorm.DB) *RollbackLogRepository {
	return &RollbackLogRepository{
		db: db,
	}
}

func (r *RollbackLogRepository) WithTx(tx *gorm.DB) *RollbackLogRepository {
	return &RollbackLogRepository{
		db: tx,
	}
}

func (r *RollbackLogRepository) Create(
	ctx context.Context,
	log *model.RollbackLog,
) error {
	return r.db.
		WithContext(ctx).
		Create(log).
		Error
}