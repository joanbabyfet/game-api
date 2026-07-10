package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

type WalletQuery struct {
	Page     int
	PageSize int
}

// Create 新增钱包
func (r *WalletRepository) Create(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

// Update 更新钱包
func (r *WalletRepository) Update(wallet *model.Wallet) error {
	return r.db.
		Model(&model.Wallet{}).
		Where("uid = ?", wallet.UID).
		Updates(wallet).
		Error
}

// Delete 删除钱包
func (r *WalletRepository) Delete(uid uint64) error {
	return r.db.
		Delete(&model.Wallet{}, "uid = ?", uid).
		Error
}

// GetByUID 根据 UID 查询钱包
func (r *WalletRepository) GetByUID(uid uint64) (*model.Wallet, error) {

	var wallet model.Wallet

	err := r.db.
		Where("uid = ?", uid).
		First(&wallet).
		Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// List 钱包列表
func (r *WalletRepository) List(q WalletQuery) ([]model.Wallet, error) {

	var wallets []model.Wallet

	offset := (q.Page - 1) * q.PageSize

	err := r.db.
		Offset(offset).
		Limit(q.PageSize).
		Find(&wallets).
		Error

	return wallets, err
}