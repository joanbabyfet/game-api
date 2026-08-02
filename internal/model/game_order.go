package model

// GameOrder 游戏注单
type GameOrder struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 幂等请求ID
	RequestID string `gorm:"column:request_id;size:64" json:"request_id"`

	// 局号
	RoundID string `gorm:"column:round_id;size:64" json:"round_id"`

	// 注单号
	OrderNo string `gorm:"column:order_no;size:64" json:"order_no"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 创建注单时的钱包模式快照：1=单一钱包 2=转账钱包
	WalletMode int8 `gorm:"column:wallet_mode" json:"wallet_mode"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 下注金额（最小货币单位）
	BetAmount int64 `gorm:"column:bet_amount" json:"bet_amount"`

	// 实际派奖金额
	WinAmount int64 `gorm:"column:win_amount" json:"win_amount"`

	// 平台盈亏
	Profit int64 `gorm:"column:profit" json:"profit"`

	// 下注前余额
	BalanceBefore int64 `gorm:"column:balance_before" json:"balance_before"`

	// 结算后余额
	BalanceAfter int64 `gorm:"column:balance_after" json:"balance_after"`

	// 默认币种
	Currency string `gorm:"column:currency;size:10" json:"currency"`

	// 回滚原因
	RollbackReason string `gorm:"column:rollback_reason;size:100" json:"rollback_reason"`

	// Spin 类型
	SpinType uint8 `gorm:"column:spin_type" json:"spin_type"`

	// Free Spin 批次ID
	FreeSpinID string `gorm:"column:free_spin_id;size:64" json:"free_spin_id"`

	// 第几次 Free Spin
	FreeSpinIndex uint32 `gorm:"column:free_spin_index" json:"free_spin_index"`

	// 状态
	Status int8 `gorm:"column:status" json:"status"`

	// 结算时间
	SettleTime int64 `gorm:"column:settle_time" json:"settle_time"`

	// 回滚时间
	RollbackTime int64 `gorm:"column:rollback_time" json:"rollback_time"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime int64 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (GameOrder) TableName() string {
	return "game_order"
}

// Spin 类型
const (
	SpinTypeNormal   uint8 = 1 // 普通 Spin
	SpinTypeFreeSpin uint8 = 2 // Free Spin
	SpinTypeBonus    uint8 = 3 // Bonus
	SpinTypeRespin   uint8 = 4 // Respin
)

// 注单状态
const (
	OrderStatusPending      int8 = 0 // 处理中 (Provider 已创建 game_order，Operator 下注结果尚未确认)
	OrderStatusBetSuccess   int8 = 1 // 扣款成功 (等待 Skynet 游戏结果)
	OrderStatusWaitSettle   int8 = 2 // 待派奖 (等待 Operator 派奖)
	OrderStatusSettled      int8 = 3 // 已结算 (Operator 派奖成功，整笔注单完成)
	OrderStatusWaitRollback int8 = 4 // 待回滚 (Skynet 明确失败，等待 Operator 回滚)
	OrderStatusRolledBack   int8 = 5 // 已回滚 (Operator 回滚成功)
	OrderStatusFailed       int8 = 6 // 失败 (明确失败且不再进行补偿)
)
