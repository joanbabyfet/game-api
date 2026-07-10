package admin

type GameListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
	Provider string `form:"provider" binding:"required"`
	Status   *int8   `json:"status"`
}

type GameListResp struct {
	ID       uint64 `json:"id"`
	GameCode string `json:"game_code"`
	Name 	 string `json:"name"`
	Provider string `json:"provider"`
	Status   int8   `json:"status"`
	Icon     string `json:"icon"`
}