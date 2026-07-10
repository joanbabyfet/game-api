package admin

type UserListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	UID      uint64 `form:"uid"`
	AgentID  uint32 `form:"agent_id"`
	Username string `form:"username"`
}