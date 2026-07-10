package model

// 转账类型
const (
	GameTransferTypeIn int8 = 1 // 转入游戏
	GameTransferTypeOut int8 = 2 // 转出游戏
)

// 转账状态
const (
	GameTransferStatusPending int8 = 0 // 处理中
	GameTransferStatusSuccess int8 = 1 // 成功
	GameTransferStatusFailed  int8 = 2 // 失败
)

// GameTransferOrder 游戏转账订单
type GameTransferOrder struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 转账订单号
	OrderNo string `gorm:"column:order_no;size:64" json:"order_no"`

	// 第三方订单号
	ThirdOrderNo string `gorm:"column:third_order_no;size:64" json:"third_order_no"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 转账类型
	TransferType int8 `gorm:"column:transfer_type" json:"transfer_type"`

	// 转账金额（最小货币单位）
	Amount uint64 `gorm:"column:amount" json:"amount"`

	// 状态
	Status int8 `gorm:"column:status" json:"status"`

	// 完成时间
	FinishTime uint32 `gorm:"column:finish_time" json:"finish_time"`

	// 创建时间
	CreateTime uint32 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (GameTransferOrder) TableName() string {
	return "game_transfer_order"
}