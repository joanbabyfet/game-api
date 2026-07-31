package model

const (
	// 免费游戏进行中
	FreeSpinStatusProcessing int8 = 0

	// 免费游戏完成
	FreeSpinStatusFinished int8 = 1

	// 免费游戏取消
	FreeSpinStatusCanceled int8 = 2
)

// FreeSpin 免费旋转批次
type FreeSpin struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// Free Spin 批次ID
	FreeSpinID string `gorm:"column:free_spin_id;size:64" json:"free_spin_id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 触发免费游戏的注单号
	TriggerOrderNo string `gorm:"column:trigger_order_no;size:64" json:"trigger_order_no"`

	// 免费游戏下注金额（用于计算派奖）
	BetAmount int64 `gorm:"column:bet_amount" json:"bet_amount"`

	// 总免费次数
	TotalCount int32 `gorm:"column:total_count" json:"total_count"`

	// 剩余免费次数
	RemainCount int32 `gorm:"column:remain_count" json:"remain_count"`

	// 累计中奖金额
	TotalWinAmount int64 `gorm:"column:total_win_amount" json:"total_win_amount"`

	// 状态
	Status int8 `gorm:"column:status" json:"status"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 完成时间
	FinishTime uint32 `gorm:"column:finish_time" json:"finish_time"`
}

// TableName 表名
func (FreeSpin) TableName() string {
	return "free_spin"
}