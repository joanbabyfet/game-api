package provider

type GameListReq struct {
	Provider string `form:"provider" binding:"required"`
}

type GameListResp struct {
	ID       uint64 `json:"id"`
	GameCode string `json:"game_code"`
	Name 	 string `json:"name"`
	Provider string `json:"provider"`
	Status   int8   `json:"status"`
	Icon     string `json:"icon"`
}

type BetReq struct {
    Token   string `json:"token"`
    RoundID string `json:"round_id"`
	GameID  uint32 `json:"game_id"`
    BetAmount int64  `json:"bet_amount"`
}

type BetResp struct {
    Balance int64 `json:"balance"`
}