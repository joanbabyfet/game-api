package admin

type GameOrderListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	OrderNo  string `form:"order_no"`
	RoundID  string `form:"round_id"`
	RequestID string `form:"request_id"`

	UID      uint64 `form:"uid"`
	AgentID  uint32 `form:"agent_id"`
	GameID   uint32 `form:"game_id"`

	Status int8 `form:"status"`
}