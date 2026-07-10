package provider

type BalanceReq struct {
	Token string `json:"token" binding:"required"`
}

type BalanceResp struct {
	Balance uint64 `json:"balance"`
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