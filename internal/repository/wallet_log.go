package repository

import (
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

type WalletLogQuery struct {
	UID        uint64
	AgentID    uint64
	GameID     uint32
	Type       string
	RefOrderNo string

	Page     int
	PageSize int
}

// Create 新增钱包流水
func (r *WalletLogRepository) Create(log *model.WalletLog) error {
	return r.db.Create(log).Error
}

// Update 更新钱包流水
func (r *WalletLogRepository) Update(log *model.WalletLog) error {
	return r.db.
		Model(&model.WalletLog{}).
		Where("id = ?", log.ID).
		Updates(log).
		Error
}

// Delete 删除钱包流水
func (r *WalletLogRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.WalletLog{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *WalletLogRepository) GetByID(id uint64) (*model.WalletLog, error) {

	var log model.WalletLog

	err := r.db.
		Where("id = ?", id).
		First(&log).
		Error

	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetByOrderNo 根据订单号查询
func (r *WalletLogRepository) GetByOrderNo(orderNo string) ([]model.WalletLog, error) {

	var logs []model.WalletLog

	err := r.db.
		Where("ref_order_no = ?", orderNo).
		Find(&logs).
		Error

	return logs, err
}

// GetByUID 根据UID查询
func (r *WalletLogRepository) GetByUID(uid uint64) ([]model.WalletLog, error) {

	var logs []model.WalletLog

	err := r.db.
		Where("uid = ?", uid).
		Order("id DESC").
		Find(&logs).
		Error

	return logs, err
}

// List 钱包流水列表
func (r *WalletLogRepository) List(q WalletLogQuery) ([]model.WalletLog, error) {

	var logs []model.WalletLog

	db := r.db

	if q.UID > 0 {
		db = db.Where("uid = ?", q.UID)
	}

	if q.AgentID > 0 {
		db = db.Where("agent_id = ?", q.AgentID)
	}

	if q.GameID > 0 {
		db = db.Where("game_id = ?", q.GameID)
	}

	if q.Type != "" {
		db = db.Where("type = ?", q.Type)
	}

	if q.RefOrderNo != "" {
		db = db.Where("ref_order_no = ?", q.RefOrderNo)
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