package admin

type RollbackLogListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	RollbackNo string `form:"rollback_no"`
	OrderNo    string `form:"order_no"`
	RequestID  string `form:"request_id"`

	AgentID uint32 `form:"agent_id"`
	UID      uint64 `form:"uid"`
	GameID   uint32 `form:"game_id"`

	RollbackType int8 `form:"rollback_type"`
	Status       int8 `form:"status"`
}