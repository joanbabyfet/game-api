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

type GameOrderQuery struct {
	OrderNo   string
	RoundID   string

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
	return r.db.Create(order).Error
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

// GetByOrderNo 根据订单号查询
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

// GetByOrderNoForUpdate 根据订单号查询并加行锁
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
// 可以防止并发请求把订单状态覆盖掉。
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
		// 先检查订单是否存在，区分“不存在”和“状态已经变化”。
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

func (r *GameOrderRepository) UpdateGameResult(
	ctx context.Context,
	orderID uint64,
	roundID string,
	winAmount int64,
	profit int64,
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
			"round_id":    roundID,
			"win_amount":  winAmount,
			"profit":      profit,
			"status":      model.OrderStatusWaitSettle,
			"update_time": now,
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
			"balance_after": balanceAfter,
			"currency":      currency,
			"status":        model.OrderStatusSettled,
			"settle_time":   now,
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