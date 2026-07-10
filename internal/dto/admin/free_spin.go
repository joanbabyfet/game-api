package admin

type FreeSpinListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	FreeSpinID string `form:"free_spin_id"`

	UID     uint64 `form:"uid"`
	AgentID uint32 `form:"agent_id"`
	GameID  uint32 `form:"game_id"`

	Status int8 `form:"status"`
}