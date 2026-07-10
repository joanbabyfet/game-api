package admin

type AgentListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	AgentID uint32 `form:"agent_id"`

	Code string `form:"code"`
	Name string `form:"name"`

	Status int8 `form:"status"`
}