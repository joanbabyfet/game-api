package model

// 回滚类型
const (
	RollbackTypeProvider int8 = 1 // Provider 自动回滚
	RollbackTypeAdmin    int8 = 2 // Admin 手动回滚
	RollbackTypeRetry    int8 = 3 // Retry 重试回滚
)

// 回滚状态
const (
	RollbackStatusPending int8 = 0 // 待执行
	RollbackStatusSuccess int8 = 1 // 成功
	RollbackStatusFailed  int8 = 2 // 失败
)

// RollbackLog 注单回滚记录
type RollbackLog struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 回滚类型
	RollbackType int8 `gorm:"column:rollback_type" json:"rollback_type"`

	// 回滚单号
	RollbackNo string `gorm:"column:rollback_no;size:64" json:"rollback_no"`

	// 原始订单号
	OrderNo string `gorm:"column:order_no;size:64" json:"order_no"`

	// 局号
	RoundID string `gorm:"column:round_id;size:64" json:"round_id"`

	// 幂等请求ID
	RequestID string `gorm:"column:request_id;size:64" json:"request_id"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 回滚金额（最小货币单位）
	Amount int64 `gorm:"column:amount" json:"amount"`

	// 回滚原因
	Reason string `gorm:"column:reason;size:255" json:"reason"`

	// 回滚状态
	Status int8 `gorm:"column:status" json:"status"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (RollbackLog) TableName() string {
	return "rollback_log"
}