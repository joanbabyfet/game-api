package provider

//查玩家余額
type BalanceReq struct {
	Token string `json:"token" binding:"required"` // JWT Token
}

type BalanceResp struct {
	Balance float64 `json:"balance"`
	Currency string `json:"currency"`
}

//Spin (玩家点击 Spin 时调用)
type SpinReq struct {
	Token 		string `json:"token" binding:"required"` // JWT Token
	
	GameCode	string `json:"game_code" binding:"required"`
	BetAmount	float64 `json:"bet_amount" binding:"required"`
	
	DebugFail 	bool `json:"debug_fail"` //测试取消下注用
}

type SpinResp struct {
    Balance float64 `json:"balance"`
	Currency string `json:"currency"`
}