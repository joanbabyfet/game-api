package model

const (
	// 已发放，尚未开始使用
	FreeSpinStatusPending int8 = 0

	// 免费游戏进行中
	FreeSpinStatusProcessing int8 = 1

	// 免费游戏完成
	FreeSpinStatusFinished int8 = 2

	// 免费游戏过期
	FreeSpinStatusExpired int8 = 3

	// 免费游戏取消
	FreeSpinStatusCanceled int8 = 4
)

const (
	// 游戏内触发
	FreeSpinSourceGameTrigger int8 = 1

	// 任务奖励
	FreeSpinSourceTaskReward int8 = 2

	// 后台人工发放
	FreeSpinSourceManual int8 = 3

	// 活动奖励
	FreeSpinSourceCampaign int8 = 4

	// 异常补偿
	FreeSpinSourceCompensation int8 = 5
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

	// 来源类型：1游戏触发 2任务奖励 3后台人工 4活动 5补偿
	SourceType int8 `gorm:"column:source_type" json:"source_type"`

	// 来源业务号
	// 游戏触发保存触发注单号
	// 任务奖励保存任务号
	// 活动奖励保存活动号
	// 后台人工保存操作号
	SourceNo string `gorm:"column:source_no;size:64" json:"source_no"`

	// 奖励发放号
	// 任务、活动、后台发放时作为幂等键
	GrantNo string `gorm:"column:grant_no;size:64" json:"grant_no"`

	// 触发免费游戏的注单号
	// 游戏内触发时有值，任务或后台发放时为空
	TriggerOrderNo string `gorm:"column:trigger_order_no;size:64" json:"trigger_order_no"`

	// 免费游戏下注金额
	// 不实际扣款，仅用于计算派奖
	BetAmount int64 `gorm:"column:bet_amount" json:"bet_amount"`

	// 总免费次数
	TotalCount int32 `gorm:"column:total_count" json:"total_count"`

	// 剩余免费次数
	RemainCount int32 `gorm:"column:remain_count" json:"remain_count"`

	// 累计中奖金额
	TotalWinAmount int64 `gorm:"column:total_win_amount" json:"total_win_amount"`

	// 状态：0待使用 1进行中 2完成 3过期 4取消
	Status int8 `gorm:"column:status" json:"status"`

	// 过期时间，0表示不过期
	ExpireTime int64 `gorm:"column:expire_time" json:"expire_time"`

	// 乐观锁版本
	Version uint32 `gorm:"column:version" json:"version"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime int64 `gorm:"column:update_time" json:"update_time"`

	// 完成时间
	FinishTime int64 `gorm:"column:finish_time" json:"finish_time"`
}

// TableName 表名
func (FreeSpin) TableName() string {
	return "free_spin"
}