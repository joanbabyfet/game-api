package provider

type OrderLogReq struct {
	BaseReq

	StartTime string `json:"start_time" form:"start_time" binding:"required"` // 2024-10-23 00:00:00 (GMT+8)
	EndTime   string `json:"end_time" form:"end_time" binding:"required"`   // 2024-10-23 00:30:00 (GMT+8)
}

type OrderLogResp struct {
	OrderNo    string  `json:"order_no"`
	RoundID    string  `json:"round_id"`
	UID        uint64  `json:"uid"`
	GameCode   string  `json:"game_code"`
	BetAmount  float64 `json:"bet_amount"`
	WinAmount  float64 `json:"win_amount"`
	Status     int8    `json:"status"`
	SettleTime int64  `json:"settle_time"`
	CreateTime int64  `json:"create_time"`
}