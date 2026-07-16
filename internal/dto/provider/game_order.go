package provider

type GameOrderListReq struct {
	AppID     string `json:"app_id" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
	Sign      string `json:"sign" binding:"required"`

	StartTime uint32 `form:"start_time"`
	EndTime   uint32 `form:"end_time"`
}

type GameOrderListResp struct {
	OrderNo   string `json:"order_no"`
	RoundID   string `json:"round_id"`

	GameID uint32 `json:"game_id"`

	BetAmount uint64 `json:"bet_amount"`
	WinAmount uint64 `json:"win_amount"`

	Status int8 `json:"status"`

	CreateTime int64 `json:"create_time"`
}