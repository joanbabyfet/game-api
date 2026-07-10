package repository

import (
	"gorm.io/gorm"

	"game-api/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type UserQuery struct {
	UID      uint64
	AgentID  uint64
	Username string

	Page     int
	PageSize int
}

// List 玩家列表
func (r *UserRepository) List(q UserQuery) ([]model.User, error) {

	var users []model.User

	db := r.db

	if q.UID > 0 {
		db = db.Where("uid = ?", q.UID)
	}

	if q.AgentID > 0 {
		db = db.Where("agent_id = ?", q.AgentID)
	}

	if q.Username != "" {
		db = db.Where("username LIKE ?", "%"+q.Username+"%")
	}

	offset := (q.Page - 1) * q.PageSize

	err := db.
		Order("uid DESC").
		Offset(offset).
		Limit(q.PageSize).
		Find(&users).
		Error

	return users, err
}

// 根据 UID 查询
func (r *UserRepository) GetByUID(uid uint64) (*model.User, error) {

	var user model.User

	err := r.db.
		Where("uid = ?", uid).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByUsername 根据用户名查询
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {

	var user model.User

	err := r.db.
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByAgentID 查询代理下所有玩家
func (r *UserRepository) GetByAgentID(agentID uint32) ([]model.User, error) {

	var users []model.User

	err := r.db.
		Where("agent_id = ?", agentID).
		Find(&users).Error

	return users, err
}

// Exists 判断玩家是否存在
func (r *UserRepository) Exists(uid uint64) (bool, error) {

	var count int64

	err := r.db.
		Model(&model.User{}).
		Where("uid = ?", uid).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 新增
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// 更新
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// 删除
func (r *UserRepository) Delete(id int64) error {
	return r.db.Delete(&model.User{}, id).Error
}