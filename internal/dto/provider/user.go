package provider

// AuthReq 登录认证
type AuthReq struct {
	Token string `json:"token" binding:"required"` // JWT Token
}

// Service -> Adapter (内部dto)
type AuthenticateReq struct {
    UID      uint64
    AgentID  uint32
    Username string
}

type AuthResp struct {
	UID      uint64 `json:"uid"`
	AgentID  uint32 `json:"agent_id"`
	Username string `json:"username"`
	Balance  uint64  `json:"balance"`
	Currency string `json:"currency"`
}

// KickReq 踢玩家请求
type KickReq struct {
	UID uint64 `json:"uid" binding:"required" validate:"required"`
}