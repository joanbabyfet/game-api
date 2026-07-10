package provider

type GameOrderListReq struct {
	Page      int `form:"page"`
	PageSize  int `form:"page_size"`

	OrderNo   string `form:"order_no"`
	RoundID   string `form:"round_id"`

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

	CreateTime uint32 `json:"create_time"`
}