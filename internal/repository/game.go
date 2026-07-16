package repository

import (
	"game-api/internal/model"
	"game-api/pkg"

	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

// GameQuery Repository 查询条件
type GameQuery struct {
	Provider string
	Status   *int8
	Page     int
	PageSize int
}

// Create 新增游戏
func (r *GameRepository) Create(game *model.Game) error {
	return r.db.Create(game).Error
}

// Update 更新游戏
func (r *GameRepository) Update(game *model.Game) error {
	return r.db.
		Model(&model.Game{}).
		Where("id = ?", game.ID).
		Updates(game).
		Error
}

// Delete 删除游戏
func (r *GameRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.Game{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *GameRepository) GetByID(id uint64) (*model.Game, error) {

	var game model.Game

	err := r.db.
		Where("id = ?", id).
		First(&game).
		Error

	if err != nil {
		return nil, err
	}

	return &game, nil
}

// GetByGameID 根据游戏ID查询
func (r *GameRepository) GetByGameID(gameID uint32) (*model.Game, error) {

	var game model.Game

	err := r.db.
		Where("game_id = ?", gameID).
		First(&game).
		Error

	if err != nil {
		return nil, err
	}

	return &game, nil
}

// GetByCode 根据游戏编码查询
func (r *GameRepository) GetByCode(code string) (*model.Game, error) {

	var game model.Game

	err := r.db.
		Where("game_code = ?", code).
		First(&game).
		Error

	if err != nil {
		return nil, err
	}

	return &game, nil
}

// List 游戏列表
func (r *GameRepository) List(q GameQuery) ([]model.Game, error) {

	var games []model.Game

	db := r.db

	if q.Provider != "" {
		db = db.Where("provider = ?", q.Provider)
	}

	if q.Status != nil {
		db = db.Where("status = ?", q.Status)
	}

	page, pageSize := pkg.Page(q.Page, q.PageSize)

	offset := (page - 1) * pageSize

	err := db.
		Debug().
		Order("sort ASC,id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&games).
		Error

	return games, err
}

// ListByIDs 根据ID查询游戏
func (r *GameRepository) ListByIDs(ids []uint32) ([]model.Game, error) {

	var list []model.Game

	if len(ids) == 0 {
		return list, nil
	}

	err := r.db.
		Where("id IN ?", ids).
		Where("status = ?", model.GameStatusEnable).
		Order("sort ASC").
		Find(&list).
		Error

	if err != nil {
		return nil, err
	}

	return list, nil
}