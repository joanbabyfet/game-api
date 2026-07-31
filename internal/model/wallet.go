package model

// Wallet 钱包
type Wallet struct {
	// 主键
	ID uint32 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 可用余额（最小货币单位）
	Balance int64 `gorm:"column:balance" json:"balance"`

	// 冻结余额（最小货币单位）
	FreezeBalance int64 `gorm:"column:freeze_balance" json:"freeze_balance"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime int64 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (Wallet) TableName() string {
	return "wallet"
}