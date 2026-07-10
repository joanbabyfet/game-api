package model

// RTPStat RTP统计
type RTPStat struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// 总Spin次数
	TotalSpin uint64 `gorm:"column:total_spin" json:"total_spin"`

	// 总下注（最小货币单位）
	TotalBet uint64 `gorm:"column:total_bet" json:"total_bet"`

	// 总派奖（最小货币单位）
	TotalWin uint64 `gorm:"column:total_win" json:"total_win"`

	// 创建时间
	CreateTime uint32 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime uint32 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (RTPStat) TableName() string {
	return "rtp_stat"
}