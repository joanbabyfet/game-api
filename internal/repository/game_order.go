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

type GameOrderRepository struct {
	db *gorm.DB
}

func NewGameOrderRepository(db *gorm.DB) *GameOrderRepository {
	return &GameOrderRepository{
		db: db,
	}
}

// ExistsProcessingTransferOrder 查询玩家是否存在尚未完成的转账钱包游戏注单。
func (r *GameOrderRepository) ExistsProcessingTransferOrder(ctx context.Context, uid uint64, agentID uint32) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("uid = ? AND agent_id = ? AND wallet_mode = ?", uid, agentID, model.WalletModeTransfer).
		Where("status IN ?", []int8{
			model.OrderStatusPending,
			model.OrderStatusBetSuccess,
			model.OrderStatusWaitSettle,
			model.OrderStatusWaitRollback,
		}).
		Limit(1).
		Count(&count).Error
	return count > 0, err
}

type GameOrderQuery struct {
	OrderNo string
	RoundID string

	UID     uint64
	AgentID uint32
	GameID  uint32

	Status *int8

	SettleStartTime uint32
	SettleEndTime   uint32

	Page     int
	PageSize int
}

// List 注单列表
func (r *GameOrderRepository) List(q GameOrderQuery) ([]model.GameOrder, error) {

	var orders []model.GameOrder

	db := r.db

	if q.OrderNo != "" {
		db = db.Where("order_no = ?", q.OrderNo)
	}

	if q.RoundID != "" {
		db = db.Where("round_id = ?", q.RoundID)
	}

	if q.UID > 0 {
		db = db.Where("uid = ?", q.UID)
	}

	if q.AgentID > 0 {
		db = db.Where("agent_id = ?", q.AgentID)
	}

	if q.GameID > 0 {
		db = db.Where("game_id = ?", q.GameID)
	}

	if q.Status != nil {
		db = db.Where("status = ?", q.Status)
	}

	// 这里用结算时间查询，而不是下注时间
	if q.SettleStartTime > 0 && q.SettleEndTime > 0 {
		db = db.Where("settle_time BETWEEN ? AND ?", q.SettleStartTime, q.SettleEndTime)
	} else if q.SettleStartTime > 0 {
		db = db.Where("settle_time >= ?", q.SettleStartTime)
	} else if q.SettleEndTime > 0 {
		db = db.Where("settle_time <= ?", q.SettleEndTime)
	}

	page, pageSize := pkg.Page(q.Page, q.PageSize)

	offset := (page - 1) * pageSize

	err := db.
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).
		Error

	return orders, err
}

// Create 新增注单
func (r *GameOrderRepository) Create(ctx context.Context, order *model.GameOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByIDForUpdate 查询并锁定注单，必须在事务中调用。
func (r *GameOrderRepository) GetByIDForUpdate(ctx context.Context, id uint64) (*model.GameOrder, error) {
	var order model.GameOrder
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateRolledBackFromWaitRollback 原子地将待退款注单更新为已退款。
func (r *GameOrderRepository) UpdateRolledBackFromWaitRollback(
	ctx context.Context,
	orderID uint64,
	balanceAfter int64,
	currency string,
	reason string,
) error {
	now := time.Now().Unix()
	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("id = ? AND status = ?", orderID, model.OrderStatusWaitRollback).
		Updates(map[string]any{
			"balance_after":   balanceAfter,
			"currency":        currency,
			"rollback_reason": reason,
			"rollback_time":   now,
			"status":          model.OrderStatusRolledBack,
			"update_time":     now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}
	return nil
}

// ListRecoverable 查询到期且未被其他Worker锁定的补偿订单。
func (r *GameOrderRepository) ListRecoverable(
	ctx context.Context,
	now int64,
	updatedBefore int64,
	maxRetry uint32,
	limit int,
) ([]model.GameOrder, error) {
	var orders []model.GameOrder
	err := r.db.WithContext(ctx).
		Where("(status IN ?) OR (status = ? AND wallet_mode = ?)", []int8{
			model.OrderStatusBetSuccess,
			model.OrderStatusWaitSettle,
			model.OrderStatusWaitRollback,
		}, model.OrderStatusPending, model.WalletModeSingle).
		Where("next_retry_time <= ?", now).
		Where("locked_until <= ?", now).
		Where("retry_count < ?", maxRetry).
		Where("update_time <= ?", updatedBefore).
		Order("update_time ASC, id ASC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// ClaimRecovery 使用条件更新抢占一笔订单。
func (r *GameOrderRepository) ClaimRecovery(ctx context.Context, orderID uint64, now, lockedUntil int64) (bool, error) {
	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("id = ? AND locked_until <= ?", orderID, now).
		Where("(status IN ?) OR (status = ? AND wallet_mode = ?)", []int8{
			model.OrderStatusBetSuccess,
			model.OrderStatusWaitSettle,
			model.OrderStatusWaitRollback,
		}, model.OrderStatusPending, model.WalletModeSingle).
		Update("locked_until", lockedUntil)
	return result.RowsAffected == 1, result.Error
}

func (r *GameOrderRepository) MarkPendingRolledBack(ctx context.Context, orderID uint64, balanceAfter int64, currency string) error {
	now := time.Now().Unix()
	result := r.db.WithContext(ctx).Model(&model.GameOrder{}).
		Where("id = ? AND status = ?", orderID, model.OrderStatusPending).
		Updates(map[string]any{
			"status":        model.OrderStatusRolledBack,
			"balance_after": balanceAfter,
			"currency":      currency,
			"rollback_time": now,
			"update_time":   now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}
	return nil
}

func (r *GameOrderRepository) MarkRetryFailed(ctx context.Context, orderID uint64, lastError string, nextRetryTime int64) error {
	return r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("id = ?", orderID).
		Updates(map[string]any{
			"retry_count":     gorm.Expr("retry_count + 1"),
			"next_retry_time": nextRetryTime,
			"locked_until":    0,
			"last_error":      lastError,
		}).Error
}

func (r *GameOrderRepository) ReleaseRecovery(ctx context.Context, orderID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("id = ?", orderID).
		Updates(map[string]any{
			"next_retry_time": 0,
			"locked_until":    0,
			"last_error":      "",
		}).Error
}

// Update 更新注单
func (r *GameOrderRepository) Update(order *model.GameOrder) error {
	return r.db.
		Model(&model.GameOrder{}).
		Where("id = ?", order.ID).
		Updates(order).
		Error
}

// Delete 删除注单
func (r *GameOrderRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.GameOrder{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *GameOrderRepository) GetByID(id uint64) (*model.GameOrder, error) {

	var order model.GameOrder

	err := r.db.
		Where("id = ?", id).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetByOrderNo 根据注单号查询
func (r *GameOrderRepository) GetByOrderNo(orderNo string) (*model.GameOrder, error) {

	var order model.GameOrder

	err := r.db.
		Where("order_no = ?", orderNo).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetByRequestID 根据请求ID查询（幂等）
func (r *GameOrderRepository) GetByRequestID(ctx context.Context, requestID string) (*model.GameOrder, error) {

	var order model.GameOrder

	err := r.db.
		Where("request_id = ?", requestID).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetByRoundID 根据局号查询
func (r *GameOrderRepository) GetByRoundID(roundID string) ([]model.GameOrder, error) {

	var orders []model.GameOrder

	err := r.db.
		Where("round_id = ?", roundID).
		Find(&orders).
		Error

	return orders, err
}

// WithTx 使用事务
func (r *GameOrderRepository) WithTx(tx *gorm.DB) *GameOrderRepository {
	return &GameOrderRepository{
		db: tx,
	}
}

// GetByOrderNoForUpdate 根据注单号查询并加行锁
func (r *GameOrderRepository) GetByOrderNoForUpdate(
	ctx context.Context,
	orderNo string,
) (*model.GameOrder, error) {

	var order model.GameOrder

	err := r.db.
		WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("order_no = ?", orderNo).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *GameOrderRepository) Rollback(
	ctx context.Context,
	orderNo string,
	reason string,
) error {

	return r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("order_no = ?", orderNo).
		Update("status", model.OrderStatusRolledBack).
		Error
}

// UpdateStatus 按当前状态更新为目标状态。
//
// 使用：
// WHERE id = ? AND status = ?
//
// 可以防止并发请求把注单状态覆盖掉。
func (r *GameOrderRepository) UpdateStatus(
	ctx context.Context,
	orderID uint64,
	fromStatus int8,
	toStatus int8,
) error {

	now := time.Now().Unix()

	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where(
			"id = ? AND status = ?",
			orderID,
			fromStatus,
		).
		Updates(map[string]interface{}{
			"status":      toStatus,
			"update_time": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		// 先检查注单是否存在，区分“不存在”和“状态已经变化”。
		var order model.GameOrder

		err := r.db.WithContext(ctx).
			Select("id", "status").
			Where("id = ?", orderID).
			First(&order).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ErrOrderNotFound
		}

		if err != nil {
			return err
		}

		return pkg.NewError(
			pkg.ORDER_STATUS_ERROR,
			"order status has changed",
		)
	}

	return nil
}

func (r *GameOrderRepository) UpdateWaitRollback(ctx context.Context, orderID uint64, reason string) error {
	now := time.Now().Unix()
	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where("id = ? AND status = ?", orderID, model.OrderStatusBetSuccess).
		Updates(map[string]any{
			"status":          model.OrderStatusWaitRollback,
			"rollback_reason": reason,
			"update_time":     now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}
	return nil
}

func (r *GameOrderRepository) UpdateGameResult(
	ctx context.Context,
	orderID uint64,
	roundID string,
	betAmount int64,
	winAmount int64,
	profit int64,
	spinType uint8,
	freeSpinID string,
	freeSpinIndex uint32,
) error {

	now := time.Now().Unix()

	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where(
			"id = ? AND status = ?",
			orderID,
			model.OrderStatusBetSuccess,
		).
		Updates(map[string]interface{}{
			"round_id":        roundID,
			"bet_amount":      betAmount,
			"win_amount":      winAmount,
			"profit":          profit,
			"spin_type":       spinType,
			"free_spin_id":    freeSpinID,
			"free_spin_index": freeSpinIndex,
			"status":          model.OrderStatusWaitSettle,
			"update_time":     now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}

	return nil
}

func (r *GameOrderRepository) UpdateSettled(
	ctx context.Context,
	orderID uint64,
	balanceBefore int64,
	balanceAfter int64,
	currency string,
) error {
	now := time.Now().Unix()

	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where(
			"id = ? AND status = ?",
			orderID,
			model.OrderStatusWaitSettle,
		).
		Updates(map[string]interface{}{
			"balance_before": balanceBefore,
			"balance_after":  balanceAfter,
			"currency":       currency,
			"status":         model.OrderStatusSettled,
			"settle_time":    now,
			"update_time":    now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}

	return nil
}

func (r *GameOrderRepository) UpdateBetSuccess(
	ctx context.Context,
	orderID uint64,
	balanceBefore int64,
) error {
	now := time.Now().Unix()

	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where(
			"id = ? AND status = ?",
			orderID,
			model.OrderStatusPending,
		).
		Updates(map[string]interface{}{
			"balance_before": balanceBefore,
			"status":         model.OrderStatusBetSuccess,
			"update_time":    now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}

	return nil
}

func (r *GameOrderRepository) UpdateRolledBack(
	ctx context.Context,
	orderID uint64,
	balanceAfter int64,
	currency string,
) error {
	now := time.Now().Unix()

	result := r.db.WithContext(ctx).
		Model(&model.GameOrder{}).
		Where(
			"id = ? AND status = ?",
			orderID,
			model.OrderStatusWaitRollback,
		).
		Updates(map[string]interface{}{
			"balance_after": balanceAfter,
			"currency":      currency,
			"status":        model.OrderStatusRolledBack,
			"rollback_time": now,
			"update_time":   now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return pkg.ErrOrderStatus
	}

	return nil
}
