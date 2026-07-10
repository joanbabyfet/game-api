package admin

type WalletListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	UID      uint64 `form:"uid"`
	AgentID  uint32 `form:"agent_id"`
}