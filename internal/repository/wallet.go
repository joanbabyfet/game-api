package repository

import (
	"context"
	"errors"

	"game-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// DB 返回当前 Repository 使用的 DB
func (r *WalletRepository) DB() *gorm.DB {
	return r.db
}

// WithTx 创建使用事务连接的 Repository
func (r *WalletRepository) WithTx(tx *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: tx,
	}
}

// Create 新增钱包
func (r *WalletRepository) Create(
	ctx context.Context,
	wallet *model.Wallet,
) error {
	return r.db.
		WithContext(ctx).
		Create(wallet).
		Error
}

// Update 更新钱包
func (r *WalletRepository) Update(
	ctx context.Context,
	wallet *model.Wallet,
) error {
	result := r.db.
		WithContext(ctx).
		Model(&model.Wallet{}).
		Where("uid = ?", wallet.UID).
		Updates(wallet)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete 删除钱包
func (r *WalletRepository) Delete(
	ctx context.Context,
	uid uint64,
) error {
	result := r.db.
		WithContext(ctx).
		Delete(&model.Wallet{}, "uid = ?", uid)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetByUID 根据 UID 查询钱包
func (r *WalletRepository) GetByUID(
	ctx context.Context,
	uid uint64,
) (*model.Wallet, error) {
	var wallet model.Wallet

	err := r.db.
		WithContext(ctx).
		Where("uid = ?", uid).
		First(&wallet).
		Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetByUIDForUpdate 查询并锁定钱包
//
// 必须在事务中调用：
//
//	tx.Transaction(func(tx *gorm.DB) error {
//	    walletRepo := repo.WithTx(tx)
//	    wallet, err := walletRepo.GetByUIDForUpdate(ctx, uid)
//	})
func (r *WalletRepository) GetByUIDForUpdate(
	ctx context.Context,
	uid uint64,
) (*model.Wallet, error) {
	var wallet model.Wallet

	err := r.db.
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("uid = ?", uid).
		First(&wallet).
		Error

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// AddBalance 增加余额
func (r *WalletRepository) AddBalance(
	ctx context.Context,
	uid uint64,
	amount int64,
) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	result := r.db.
		WithContext(ctx).
		Model(&model.Wallet{}).
		Where("uid = ?", uid).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// SubBalance 扣除余额
//
// WHERE balance >= amount 可以保证余额不会扣成负数。
func (r *WalletRepository) SubBalance(
	ctx context.Context,
	uid uint64,
	amount int64,
) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	result := r.db.
		WithContext(ctx).
		Model(&model.Wallet{}).
		Where("uid = ?", uid).
		Where("balance >= ?", amount).
		UpdateColumn("balance", gorm.Expr("balance - ?", amount))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// List 钱包列表
func (r *WalletRepository) List(
	ctx context.Context,
	q WalletQuery,
) ([]model.Wallet, error) {
	var wallets []model.Wallet

	if q.Page <= 0 {
		q.Page = 1
	}

	if q.PageSize <= 0 {
		q.PageSize = 20
	}

	offset := (q.Page - 1) * q.PageSize

	err := r.db.
		WithContext(ctx).
		Offset(offset).
		Limit(q.PageSize).
		Find(&wallets).
		Error

	return wallets, err
}