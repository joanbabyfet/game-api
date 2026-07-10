package admin

type JackpotListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	GameID uint32 `form:"game_id"`
}