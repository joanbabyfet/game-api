package model

// WalletLog 钱包流水
type WalletLog struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 流水类型
	Type string `gorm:"column:type;size:30" json:"type"`

	// 变动金额（最小货币单位）
	Amount uint64 `gorm:"column:amount" json:"amount"`

	// 变动前余额（最小货币单位）
	BalanceBefore uint64 `gorm:"column:balance_before" json:"balance_before"`

	// 变动后余额（最小货币单位）
	BalanceAfter uint64 `gorm:"column:balance_after" json:"balance_after"`

	// 关联订单号
	RefOrderNo string `gorm:"column:ref_order_no;size:64" json:"ref_order_no"`

	// 创建时间
	CreateTime uint32 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (WalletLog) TableName() string {
	return "wallet_log"
}