package model

import "gorm.io/datatypes"

// UserGameData 玩家游戏数据
type UserGameData struct {
	// 玩家ID
	UID uint64 `gorm:"column:uid;primaryKey" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id;primaryKey" json:"game_id"`

	// 剩余免费次数
	FreeSpin int32 `gorm:"column:free_spin" json:"free_spin"`

	// Bonus 状态
	BonusState datatypes.JSON `gorm:"column:bonus_state" json:"bonus_state"`

	// 最后进入游戏时间
	LastLoginTime uint32 `gorm:"column:last_login_time" json:"last_login_time"`
}

// TableName 表名
func (UserGameData) TableName() string {
	return "user_game_data"
}