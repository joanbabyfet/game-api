package repository

import (
	"game-api/internal/model"
	"game-api/pkg"

	"gorm.io/gorm"
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
func (r *GameOrderRepository) Create(order *model.GameOrder) error {
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
func (r *GameOrderRepository) GetByRequestID(requestID string) (*model.GameOrder, error) {

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