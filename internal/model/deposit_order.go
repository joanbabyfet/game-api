package model

// 充值注单状态
const (
	DepositOrderStatusPending int8 = 0 // 待处理
	DepositOrderStatusSuccess int8 = 1 // 成功
	DepositOrderStatusFailed  int8 = 2 // 失败
)

// DepositOrder 充值注单
type DepositOrder struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 充值注单号
	OrderNo string `gorm:"column:order_no;size:64" json:"order_no"`

	// 第三方注单号
	ThirdOrderNo string `gorm:"column:third_order_no;size:64" json:"third_order_no"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 金额（最小货币单位）
	Amount int64 `gorm:"column:amount" json:"amount"`

	// 状态
	Status int8 `gorm:"column:status" json:"status"`

	// 完成时间
	FinishTime uint32 `gorm:"column:finish_time" json:"finish_time"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (DepositOrder) TableName() string {
	return "deposit_order"
}