package model

// Agent 代理
type Agent struct {
	// 代理ID
	ID uint32 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 代理编码（唯一）
	AgentCode string `gorm:"column:agent_code;size:50" json:"agent_code"`

	// 代理名称
	AgentName string `gorm:"column:agent_name;size:100" json:"agent_name"`

	// 上级代理ID
	PID uint32 `gorm:"column:pid" json:"pid"`

	// 默认币种
	Currency string `gorm:"column:currency;size:10" json:"currency"`

	// 国家代码
	Country string `gorm:"column:country;size:2" json:"country"`

	// 默认语言
	Language string `gorm:"column:language;size:10" json:"language"`

	// 时区
	Timezone string `gorm:"column:timezone;size:50" json:"timezone"`

	// API Key
	APIKey string `gorm:"column:api_key;size:64" json:"api_key"`

	// API Secret
	SecretKey string `gorm:"column:secret_key;size:128" json:"secret_key"`

	// 回调地址
	CallbackURL string `gorm:"column:callback_url;size:255" json:"callback_url"`

	// 状态：0=禁用 1=启用
	Status int8 `gorm:"column:status" json:"status"`

	// 备注
	Remark string `gorm:"column:remark;size:255" json:"remark"`

	// 创建时间
	CreateTime uint32 `gorm:"column:create_time" json:"create_time"`

	// 修改时间
	UpdateTime uint32 `gorm:"column:update_time" json:"update_time"`
}

// TableName 表名
func (Agent) TableName() string {
	return "agent"
}