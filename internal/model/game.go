package model

// 游戏状态
const (
	GameStatusDisable int8 = 0 // 禁用
	GameStatusEnable  int8 = 1 // 启用
)

// 游戏维护状态
const (
	GameMaintenanceOff int8 = 0 // 正常
	GameMaintenanceOn  int8 = 1 // 维护中
)

// 是否支持试玩
const (
	GameDemoNo  int8 = 0 // 不支持
	GameDemoYes int8 = 1 // 支持
)

// Game 游戏
type Game struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement"` // 主键ID

	Provider string `gorm:"column:provider"` // 游戏厂商
	GameCode string `gorm:"column:game_code"` // 游戏标识
	GameName string `gorm:"column:game_name"` // 游戏名称

	RTP float32 `gorm:"column:rtp"` // 理论RTP

	Icon   string `gorm:"column:icon"`   // 游戏图标
	Banner string `gorm:"column:banner"` // Banner图片

	Sort uint `gorm:"column:sort"` // 排序

	SupportDemo int8 `gorm:"column:support_demo"` // 是否支持试玩：0否 1是

	Maintenance    int8   `gorm:"column:maintenance"`     // 是否维护：0否 1是
	MaintenanceMsg string `gorm:"column:maintenance_msg"` // 维护提示

	Status int8 `gorm:"column:status"` // 状态：0禁用 1启用

	CreateTime uint32 `gorm:"column:create_time"` // 创建时间
	UpdateTime uint32 `gorm:"column:update_time"` // 更新时间
}

// TableName 表名
func (Game) TableName() string {
	return "game"
}