package model

const (
	// AgentGame Status
	AgentGameStatusDisable int8 = 0 // 禁用
	AgentGameStatusEnable  int8 = 1 // 启用
)

const (
	// AgentGame Risk
	AgentGameRiskDisable int8 = 0 // 关闭
	AgentGameRiskEnable  int8 = 1 // 开启
)

// AgentGame 代理游戏配置
type AgentGame struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 代理专属 RTP
	RTP float64 `gorm:"column:rtp" json:"rtp"`

	// 佣金比例
	CommissionRate float64 `gorm:"column:commission_rate" json:"commission_rate"`

	// 最小下注（最小货币单位）
	MinBet int64 `gorm:"column:min_bet" json:"min_bet"`

	// 最大下注（最小货币单位）
	MaxBet int64 `gorm:"column:max_bet" json:"max_bet"`

	// 单局最大派奖
	MaxWin int64 `gorm:"column:max_win" json:"max_win"`

	// 单局最大倍率
	MaxMultiple int32 `gorm:"column:max_multiple" json:"max_multiple"`

	// 是否开启代理风控（0=关闭 1=开启）
	RiskEnable int8 `gorm:"column:risk_enable" json:"risk_enable"`

	// 状态（0=关闭 1=开启）
	Status int8 `gorm:"column:status" json:"status"`

	// 创建时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime int64 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (AgentGame) TableName() string {
	return "agent_game"
}