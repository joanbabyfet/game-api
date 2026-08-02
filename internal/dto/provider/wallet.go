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
	BetAmount   float64 `json:"bet_amount"`
	
	RequestID	string `json:"request_id" binding:"required"`
	FreeSpinID string  `json:"free_spin_id"` //是否免费旋转

	DebugFail 	bool `json:"debug_fail"` //测试取消下注用
}

type SpinResp struct {
    Balance float64 `json:"balance"`
	Currency string `json:"currency"`

	RoundID   string  `json:"round_id"`
	WinAmount float64 `json:"win_amount"`
	SpinType  uint8   `json:"spin_type"`

	FreeSpinID          string `json:"free_spin_id,omitempty"`
	FreeSpinIndex       uint32 `json:"free_spin_index,omitempty"`
	FreeSpinTotalCount  uint32 `json:"free_spin_total_count,omitempty"`
	FreeSpinRemainCount uint32 `json:"free_spin_remain_count,omitempty"`
}

// ChangeBalanceReq 钱包余额变更请求
type ChangeBalanceReq struct {
	UID        uint64 `json:"uid" binding:"required"`
	AgentID    uint32 `json:"agent_id" binding:"required"`
	GameID     uint32 `json:"game_id" binding:"required"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
	Type       string `json:"type" binding:"required"`
	RefOrderNo string `json:"ref_order_no" binding:"required"`
}

// RollbackReq 钱包回滚请求
type RollbackReq struct {
	OrderNo     string `json:"order_no" binding:"required"`
	RollbackType int8    `json:"rollback_type" binding:"required"`
	RequestID   string `json:"request_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
}