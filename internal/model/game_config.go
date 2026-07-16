package model

import "gorm.io/datatypes"

// GameConfig 游戏配置
type GameConfig struct {
	// 主键
	ID uint32 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 配置Key
	ConfigKey string `gorm:"column:config_key;size:50" json:"config_key"`

	// 配置内容(JSON)
	ConfigValue datatypes.JSON `gorm:"column:config_value" json:"config_value"`

	// 配置版本
	Version uint32 `gorm:"column:version" json:"version"`

	// 备注
	Remark string `gorm:"column:remark;size:255" json:"remark"`

	// 状态：0=禁用 1=启用
	Status int8 `gorm:"column:status" json:"status"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime int64 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (GameConfig) TableName() string {
	return "game_config"
}