package model

import "gorm.io/datatypes"

// 请求状态
const (
	RequestStatusFailed  int8 = 0 // 失败
	RequestStatusSuccess int8 = 1 // 成功
)

// RequestLog API 请求日志
type RequestLog struct {
	// 主键
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 请求唯一ID（幂等）
	RequestID string `gorm:"column:request_id;size:64" json:"request_id"`

	// 注单号
	OrderNo string `gorm:"column:order_no;size:64" json:"order_no"`

	// 代理ID
	AgentID uint32 `gorm:"column:agent_id" json:"agent_id"`

	// 玩家ID
	UID uint64 `gorm:"column:uid" json:"uid"`

	// 游戏ID
	GameID uint32 `gorm:"column:game_id" json:"game_id"`

	// API 名称
	API string `gorm:"column:api;size:50" json:"api"`

	// 请求内容
	Request datatypes.JSON `gorm:"column:request" json:"request"`

	// 响应内容
	Response datatypes.JSON `gorm:"column:response" json:"response"`

	// 请求状态
	Status int8 `gorm:"column:status" json:"status"`

	// 业务错误码
	ErrorCode int32 `gorm:"column:error_code" json:"error_code"`

	// 耗时(ms)
	CostTime uint32 `gorm:"column:cost_time" json:"cost_time"`

	// 请求IP
	IP string `gorm:"column:ip;size:45" json:"ip"`

	// 创建时间
	CreateTime uint32 `gorm:"column:create_time" json:"create_time"`
}

// TableName 表名
func (RequestLog) TableName() string {
	return "request_log"
}