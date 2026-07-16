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
	RoundID  	string `json:"round_id" binding:"required"`
	GameCode 	string `json:"game_code" binding:"required"` // 游戏标识
	BetAmount 	float64 `json:"bet_amount" binding:"required"`
	BetTime 	int64 `json:"bet_time" binding:"required"` //下注时间
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
	RoundID  	string `json:"round_id" binding:"required"`
	GameCode 	string `json:"game_code" binding:"required"` // 游戏标识
	BetAmount 	float64 `json:"bet_amount" binding:"required"`
	WinAmount 	float64 `json:"win_amount" binding:"required"`
	SettleTime 	int64  `json:"settle_time" binding:"required"`
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
	RoundID  	string `json:"round_id" binding:"required"`
	GameCode 	string `json:"game_code" binding:"required"` // 游戏标识
	BetAmount	float64 `json:"bet_amount" binding:"required"`
}

type RollbackResp struct {
	Balance float64  `json:"balance"`
	Currency string `json:"currency"`
}