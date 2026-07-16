package model

// Jackpot 类型
const (
	JackpotTypeMini  = "mini"
	JackpotTypeMinor = "minor"
	JackpotTypeMajor = "major"
	JackpotTypeGrand = "grand"
)

// JackpotLog Jackpot 中奖记录
type JackpotLog struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// Jackpot 类型
	JackpotType string `gorm:"column:jackpot_type;size:20" json:"jackpot_type"`

	// 中奖金额（最小货币单位）
	Amount uint64 `gorm:"column:amount" json:"amount"`

	// 中奖时间
	CreateTime int64 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (JackpotLog) TableName() string {
	return "jackpot_log"
}