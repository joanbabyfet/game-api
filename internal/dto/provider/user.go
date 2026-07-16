package provider

// KickReq 踢玩家请求
type KickReq struct {
	BaseReq
	
	GameCode 	string `json:"game_code" form:"game_code" binding:"required"` // 遊戲ID
	PlayerID	string `json:"player_id" form:"player_id" binding:"required"`     // 玩家UID
}

type KickResp struct {
}


//登录
type LoginReq struct {
	BaseReq

	GameCode	string `json:"game_code" form:"game_code" binding:"required"`
	PlayerID	string `json:"player_id" form:"player_id" binding:"required"` //运营商用户标识
	//Nickname	string `json:"nickname" form:"nickname"` //运营商用户昵称(选填)
}

type LoginResp struct {
}