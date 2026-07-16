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
}

type SpinResp struct {
    Balance float64 `json:"balance"`
	Currency string `json:"currency"`
}

type RollbackReq struct {
	Token     string `json:"token" binding:"required"`
	OrderNo   string `json:"order_no" binding:"required"`
	RequestID string `json:"request_id" binding:"required"`
	Reason    string `json:"reason"`
}

type RollbackResp struct {
    Balance uint64 `json:"balance"`
}