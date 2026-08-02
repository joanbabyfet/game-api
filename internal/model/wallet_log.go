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
	Amount int64 `gorm:"column:amount" json:"amount"`

	// 变动前余额（最小货币单位）
	BalanceBefore int64 `gorm:"column:balance_before" json:"balance_before"`

	// 变动后余额（最小货币单位）
	BalanceAfter int64 `gorm:"column:balance_after" json:"balance_after"`

	// 关联注单号
	RefOrderNo string `gorm:"column:ref_order_no;size:64" json:"ref_order_no"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (WalletLog) TableName() string {
	return "wallet_log"
}

// Wallet Log Type
const (
	WalletLogTypeBet       = "BET"        // 下注
	WalletLogTypeWin       = "WIN"        // 派奖
	WalletLogTypeBonus     = "BONUS"      // Bonus 奖励
	WalletLogTypeFreeSpin  = "FREE_SPIN"  // 免费游戏奖励
	WalletLogTypeJackpot   = "JACKPOT"    // 奖池奖励
	WalletLogTypeDeposit   = "DEPOSIT"    // 充值
	WalletLogTypeWithdraw  = "WITHDRAW"   // 提现
	WalletLogTypeAdminAdd  = "ADMIN_ADD"  // 后台加钱
	WalletLogTypeAdminSub  = "ADMIN_SUB"  // 后台扣钱
	WalletLogTypeRollback  = "ROLLBACK"   // 回滚
)