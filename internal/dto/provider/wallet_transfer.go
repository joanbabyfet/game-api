package provider

type TransferReq struct {
	BaseReq
	RequestID    string  `json:"request_id" binding:"required"`
	PlayerID     string  `json:"player_id" binding:"required"`
	GameCode     string  `json:"game_code" binding:"required"`
	ThirdOrderNo string  `json:"third_order_no" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Currency     string  `json:"currency" binding:"required"`
}

type TransferStatusReq struct {
	BaseReq
	ThirdOrderNo string `json:"third_order_no" binding:"required"`
}

type TransferResp struct {
	OrderNo      string  `json:"order_no"`
	ThirdOrderNo string  `json:"third_order_no"`
	TransferType int8    `json:"transfer_type"`
	Amount       float64 `json:"amount"`
	Balance      float64 `json:"balance"`
	Currency     string  `json:"currency"`
	Status       int8    `json:"status"`
	ErrorCode    int     `json:"error_code,omitempty"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

type TransferStatusResp = TransferResp
