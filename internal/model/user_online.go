package model

// UserOnline 在线玩家
type UserOnline struct {
	// 玩家ID
	UID uint64 `gorm:"column:uid;primaryKey" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 当前游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 登录时间
	LoginTime uint32 `gorm:"column:login_time" json:"login_time"`

	// 最后活跃时间
	LastActiveTime uint32 `gorm:"column:last_active_time" json:"last_active_time"`

	// 所在游戏服务器ID
	ServerID string `gorm:"column:server_id;size:50" json:"server_id"`
}

// TableName 表名
func (UserOnline) TableName() string {
	return "user_online"
}