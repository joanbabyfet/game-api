package repository

import (
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

type RollbackLogQuery struct {
	RollbackNo string
	OrderNo    string
	RequestID  string

	AgentID uint64
	UID      uint64
	GameID   uint32

	RollbackType *int8
	Status       *int8

	Page     int
	PageSize int
}

// Create 新增回滚日志
func (r *RollbackLogRepository) Create(log *model.RollbackLog) error {
	return r.db.Create(log).Error
}

// Update 更新回滚日志
func (r *RollbackLogRepository) Update(log *model.RollbackLog) error {
	return r.db.
		Model(&model.RollbackLog{}).
		Where("id = ?", log.ID).
		Updates(log).
		Error
}

// Delete 删除回滚日志
func (r *RollbackLogRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.RollbackLog{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *RollbackLogRepository) GetByID(id uint64) (*model.RollbackLog, error) {

	var log model.RollbackLog

	err := r.db.
		Where("id = ?", id).
		First(&log).
		Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetByRollbackNo 根据回滚单号查询
func (r *RollbackLogRepository) GetByRollbackNo(rollbackNo string) (*model.RollbackLog, error) {

	var log model.RollbackLog

	err := r.db.
		Where("rollback_no = ?", rollbackNo).
		First(&log).
		Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetByOrderNo 根据订单号查询
func (r *RollbackLogRepository) GetByOrderNo(orderNo string) ([]model.RollbackLog, error) {

	var logs []model.RollbackLog

	err := r.db.
		Where("order_no = ?", orderNo).
		Order("id DESC").
		Find(&logs).
		Error

	return logs, err
}

// List 回滚日志列表
func (r *RollbackLogRepository) List(q RollbackLogQuery) ([]model.RollbackLog, error) {

	var logs []model.RollbackLog

	db := r.db

	if q.RollbackNo != "" {
		db = db.Where("rollback_no = ?", q.RollbackNo)
	}

	if q.OrderNo != "" {
		db = db.Where("order_no = ?", q.OrderNo)
	}

	if q.RequestID != "" {
		db = db.Where("request_id = ?", q.RequestID)
	}

	if q.AgentID > 0 {
		db = db.Where("agent_id = ?", q.AgentID)
	}

	if q.UID > 0 {
		db = db.Where("uid = ?", q.UID)
	}

	if q.GameID > 0 {
		db = db.Where("game_id = ?", q.GameID)
	}

	if q.RollbackType != nil {
		db = db.Where("rollback_type = ?", q.RollbackType)
	}

	if q.Status != nil {
		db = db.Where("status = ?", q.Status)
	}

	offset := (q.Page - 1) * q.PageSize

	err := db.
		Order("id DESC").
		Offset(offset).
		Limit(q.PageSize).
		Find(&logs).
		Error

	return logs, err
}