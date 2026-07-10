package admin

type WalletLogListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	UID        uint64 `form:"uid"`
	AgentID    uint32 `form:"agent_id"`
	GameID     uint32 `form:"game_id"`
	Type       string `form:"type"`
	RefOrderNo string `form:"ref_order_no"`
}