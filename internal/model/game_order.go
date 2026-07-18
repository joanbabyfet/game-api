package model

import "gorm.io/datatypes"

// GameOrder 游戏注单
type GameOrder struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 注单号
	OrderNo string `gorm:"column:order_no;size:64" json:"order_no"`

	// 局号
	RoundID string `gorm:"column:round_id;size:64" json:"round_id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 下注金额（最小货币单位）
	BetAmount uint64 `gorm:"column:bet_amount" json:"bet_amount"`

	// 风控前派奖金额
	OriginalWin uint64 `gorm:"column:original_win" json:"original_win"`

	// 实际派奖金额
	WinAmount uint64 `gorm:"column:win_amount" json:"win_amount"`

	// 是否命中风控
	RiskHit int8 `gorm:"column:risk_hit" json:"risk_hit"`

	// 风控原因
	RiskReason string `gorm:"column:risk_reason;size:50" json:"risk_reason"`

	// 平台盈亏
	Profit int64 `gorm:"column:profit" json:"profit"`

	// 下注前余额
	BalanceBefore int64 `gorm:"column:balance_before" json:"balance_before"`

	// 结算后余额
	BalanceAfter int64 `gorm:"column:balance_after" json:"balance_after"`

	// 默认币种
	Currency string `gorm:"column:currency;size:10" json:"currency"`

	// 开奖结果
	Reels datatypes.JSON `gorm:"column:reels" json:"reels"`

	// 中奖线路
	WinLines datatypes.JSON `gorm:"column:win_lines" json:"win_lines"`

	// 是否免费游戏
	IsFreeSpin int8 `gorm:"column:is_free_spin" json:"is_free_spin"`

	// Free Spin 批次ID
	FreeSpinID string `gorm:"column:free_spin_id;size:64" json:"free_spin_id"`

	// 第几次 Free Spin
	FreeSpinIndex int32 `gorm:"column:free_spin_index" json:"free_spin_index"`

	// 回滚原因
	RollbackReason string `gorm:"column:rollback_reason;size:100" json:"rollback_reason"`

	// 幂等请求ID
	RequestID string `gorm:"column:request_id;size:64" json:"request_id"`

	// 状态
	Status int8 `gorm:"column:status" json:"status"`

	// 结算时间
	SettleTime int64 `gorm:"column:settle_time" json:"settle_time"`

	// 回滚时间
	RollbackTime int64 `gorm:"column:rollback_time" json:"rollback_time"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (GameOrder) TableName() string {
	return "game_order"
}