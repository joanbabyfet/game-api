package repository

import (
	"context"
	"errors"
	"game-api/internal/model"
	"game-api/pkg"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletTransferRepository struct {
	db *gorm.DB
}

func NewWalletTransferRepository(db *gorm.DB) *WalletTransferRepository {
	return &WalletTransferRepository{db: db}
}

func (r *WalletTransferRepository) WithTx(tx *gorm.DB) *WalletTransferRepository {
	return &WalletTransferRepository{db: tx}
}

func (r *WalletTransferRepository) Create(ctx context.Context, order *model.WalletTransfer) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *WalletTransferRepository) GetByThirdOrderNo(ctx context.Context, agentID uint32, thirdOrderNo string) (*model.WalletTransfer, error) {
	var order model.WalletTransfer
	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND third_order_no = ?", agentID, thirdOrderNo).
		First(&order).Error
	return &order, err
}

func (r *WalletTransferRepository) GetByRequestID(ctx context.Context, agentID uint32, requestID string) (*model.WalletTransfer, error) {
	var order model.WalletTransfer
	err := r.db.WithContext(ctx).
		Where("agent_id = ? AND request_id = ?", agentID, requestID).
		First(&order).Error
	return &order, err
}

func (r *WalletTransferRepository) GetByThirdOrderNoForUpdate(ctx context.Context, agentID uint32, thirdOrderNo string) (*model.WalletTransfer, error) {
	var order model.WalletTransfer
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("agent_id = ? AND third_order_no = ?", agentID, thirdOrderNo).
		First(&order).Error
	return &order, err
}

func (r *WalletTransferRepository) UpdateSuccess(ctx context.Context, orderID uint64, before, after int64) error {
	now := time.Now().Unix()
	result := r.db.WithContext(ctx).Model(&model.WalletTransfer{}).
		Where("id = ? AND status = ?", orderID, model.GameTransferStatusPending).
		Updates(map[string]any{
			"balance_before": before,
			"balance_after":  after,
			"status":         model.GameTransferStatusSuccess,
			"finish_time":    now,
			"update_time":    now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return pkg.ErrTransferOrderConflict
	}
	return nil
}

func (r *WalletTransferRepository) UpdateFailed(ctx context.Context, orderID uint64, code int, message string) error {
	now := time.Now().Unix()
	result := r.db.WithContext(ctx).Model(&model.WalletTransfer{}).
		Where("id = ? AND status = ?", orderID, model.GameTransferStatusPending).
		Updates(map[string]any{
			"status":        model.GameTransferStatusFailed,
			"error_code":    code,
			"error_message": message,
			"finish_time":   now,
			"update_time":   now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return pkg.ErrTransferOrderConflict
	}
	return nil
}

func (r *WalletTransferRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.WalletTransfer, error) {
	var order model.WalletTransfer
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrTransferOrderNotFound
	}
	return &order, err
}
