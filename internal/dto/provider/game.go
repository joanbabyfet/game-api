package provider

type GameListReq struct {
	BaseReq
}

// GameListResp 游戏列表
type GameListResp struct {
	GameCode string `json:"game_code"` // 游戏标识
	GameName string `json:"game_name"` // 游戏名称
	Provider string `json:"provider"`  // 游戏厂商
	Icon     string `json:"icon"`      // 游戏图标
}

type GameURLReq struct {
	BaseReq

	GameCode	string `json:"game_code" form:"game_code" binding:"required"`
	PlayerID	string `json:"player_id" form:"player_id" binding:"required"`
	Lang  		string `json:"lang" form:"lang" binding:"required"`
}

type GameURLResp struct {
	URL      string `json:"url"`
}