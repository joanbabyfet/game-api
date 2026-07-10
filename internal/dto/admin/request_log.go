package admin

type RequestLogListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	RequestID string `form:"request_id"`
	OrderNo   string `form:"order_no"`

	AgentID uint32 `form:"agent_id"`
	UID      uint64 `form:"uid"`
	GameID   uint32 `form:"game_id"`

	API    string `form:"api"`
	Status int8   `form:"status"`
}