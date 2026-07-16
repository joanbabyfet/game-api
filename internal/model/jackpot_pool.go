package model

// JackpotPool Jackpot 奖池
type JackpotPool struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// Mini 奖池金额
	Mini uint64 `gorm:"column:mini" json:"mini"`

	// Minor 奖池金额
	Minor uint64 `gorm:"column:minor" json:"minor"`

	// Major 奖池金额
	Major uint64 `gorm:"column:major" json:"major"`

	// Grand 奖池金额
	Grand uint64 `gorm:"column:grand" json:"grand"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime int64 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (JackpotPool) TableName() string {
	return "jackpot_pool"
}