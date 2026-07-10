package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type FreeSpinRepository struct {
	db *gorm.DB
}

func NewFreeSpinRepository(db *gorm.DB) *FreeSpinRepository {
	return &FreeSpinRepository{
		db: db,
	}
}

type FreeSpinQuery struct {
	FreeSpinID string
	UID        uint64
	AgentID    uint64
	GameID     uint32
	Status     *int8

	Page     int
	PageSize int
}

// Create 新增 Free Spin
func (r *FreeSpinRepository) Create(freeSpin *model.FreeSpin) error {
	return r.db.Create(freeSpin).Error
}

// Update 更新 Free Spin
func (r *FreeSpinRepository) Update(freeSpin *model.FreeSpin) error {
	return r.db.
		Model(&model.FreeSpin{}).
		Where("id = ?", freeSpin.ID).
		Updates(freeSpin).
		Error
}

// Delete 删除 Free Spin
func (r *FreeSpinRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.FreeSpin{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *FreeSpinRepository) GetByID(id uint64) (*model.FreeSpin, error) {

	var freeSpin model.FreeSpin

	err := r.db.
		Where("id = ?", id).
		First(&freeSpin).
		Error

	if err != nil {
		return nil, err
	}

	return &freeSpin, nil
}

// GetByFreeSpinID 根据 FreeSpinID 查询
func (r *FreeSpinRepository) GetByFreeSpinID(freeSpinID string) (*model.FreeSpin, error) {

	var freeSpin model.FreeSpin

	err := r.db.
		Where("free_spin_id = ?", freeSpinID).
		First(&freeSpin).
		Error

	if err != nil {
		return nil, err
	}

	return &freeSpin, nil
}

// GetByTriggerOrderNo 根据触发订单查询
func (r *FreeSpinRepository) GetByTriggerOrderNo(orderNo string) (*model.FreeSpin, error) {

	var freeSpin model.FreeSpin

	err := r.db.
		Where("trigger_order_no = ?", orderNo).
		First(&freeSpin).
		Error

	if err != nil {
		return nil, err
	}

	return &freeSpin, nil
}

// List Free Spin 列表
func (r *FreeSpinRepository) List(q FreeSpinQuery) ([]model.FreeSpin, error) {

	var list []model.FreeSpin

	db := r.db

	if q.FreeSpinID != "" {
		db = db.Where("free_spin_id = ?", q.FreeSpinID)
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

	offset := (q.Page - 1) * q.PageSize

	err := db.
		Order("id DESC").
		Offset(offset).
		Limit(q.PageSize).
		Find(&list).
		Error

	return list, err
}