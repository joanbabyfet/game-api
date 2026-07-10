package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type RequestLogRepository struct {
	db *gorm.DB
}

func NewRequestLogRepository(db *gorm.DB) *RequestLogRepository {
	return &RequestLogRepository{
		db: db,
	}
}

type RequestLogQuery struct {
	RequestID string
	OrderNo   string

	AgentID uint64
	UID      uint64
	GameID   uint32

	API    string
	Status *int8

	Page     int
	PageSize int
}

// Create 新增请求日志
func (r *RequestLogRepository) Create(log *model.RequestLog) error {
	return r.db.Create(log).Error
}

// Update 更新请求日志
func (r *RequestLogRepository) Update(log *model.RequestLog) error {
	return r.db.
		Model(&model.RequestLog{}).
		Where("id = ?", log.ID).
		Updates(log).
		Error
}

// Delete 删除请求日志
func (r *RequestLogRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.RequestLog{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *RequestLogRepository) GetByID(id uint64) (*model.RequestLog, error) {

	var log model.RequestLog

	err := r.db.
		Where("id = ?", id).
		First(&log).
		Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetByRequestID 根据RequestID查询
func (r *RequestLogRepository) GetByRequestID(requestID string) (*model.RequestLog, error) {

	var log model.RequestLog

	err := r.db.
		Where("request_id = ?", requestID).
		First(&log).
		Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetByOrderNo 根据订单号查询
func (r *RequestLogRepository) GetByOrderNo(orderNo string) ([]model.RequestLog, error) {

	var logs []model.RequestLog

	err := r.db.
		Where("order_no = ?", orderNo).
		Order("id DESC").
		Find(&logs).
		Error

	return logs, err
}

// List 请求日志列表
func (r *RequestLogRepository) List(q RequestLogQuery) ([]model.RequestLog, error) {

	var logs []model.RequestLog

	db := r.db

	if q.RequestID != "" {
		db = db.Where("request_id = ?", q.RequestID)
	}

	if q.OrderNo != "" {
		db = db.Where("order_no = ?", q.OrderNo)
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

	if q.API != "" {
		db = db.Where("api = ?", q.API)
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