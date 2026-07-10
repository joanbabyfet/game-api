package model

// 玩家状态
const (
	//禁用
	UserStatusDisable int8 = 0

	//启用
	UserStatusEnable int8 = 1
)

// User 玩家
type User struct {
	// 玩家唯一ID（雪花算法）
	UID uint64 `gorm:"column:uid;primaryKey" json:"uid"`
	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`
	// 用户名
	Username string `gorm:"column:username;size:50" json:"username"`
	// 昵称
	Nickname string `gorm:"column:nickname;size:50" json:"nickname"`
	// 状态：0=禁用 1=启用
	Status int8 `gorm:"column:status" json:"status"`
	// 最后登录时间
	LastLoginTime uint32 `gorm:"column:last_login_time" json:"last_login_time"`
	// 创建时间（Unix Timestamp）
	CreateTime uint32 `gorm:"column:create_time" json:"create_time"`
	// 修改时间（Unix Timestamp）
	UpdateTime uint32 `gorm:"column:update_time" json:"update_time"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}