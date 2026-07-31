package mock

//查余额
type BalanceReq struct {
	BaseReq

	PlayerID string `json:"player_id" binding:"required"`
}

type BalanceResp struct {
	Balance float64  `json:"balance"`
	Currency string `json:"currency"`
}

//下注
type BetReq struct {
	BaseReq

	PlayerID 	string  `json:"player_id" binding:"required"`
	OrderNo   	string `json:"order_no" binding:"required"`
	GameCode 	string `json:"game_code" binding:"required"` // 游戏标识
	BetAmount 	float64 `json:"bet_amount" binding:"required"`
}

type BetResp struct {
	Balance float64  `json:"balance"`
	Currency string `json:"currency"`
}

//结算
type SettleReq struct {
	BaseReq

	PlayerID 	string  `json:"player_id" binding:"required"`
	OrderNo   	string `json:"order_no" binding:"required"`
	GameCode 	string `json:"game_code" binding:"required"` // 游戏标识
	WinAmount 	float64 `json:"win_amount" binding:"gte=0"` //可以是 0，且不能小于 0
}

type SettleResp struct {
	Balance float64  `json:"balance"`
	Currency string `json:"currency"`
}

//取消下注
type RollbackReq struct {
	BaseReq

	PlayerID 	string  `json:"player_id" binding:"required"`
	OrderNo		string `json:"order_no" binding:"required"`
	GameCode 	string `json:"game_code" binding:"required"` // 游戏标识
	BetAmount	float64 `json:"bet_amount" binding:"required"`
}

type RollbackResp struct {
	Balance float64  `json:"balance"`
	Currency string `json:"currency"`
}