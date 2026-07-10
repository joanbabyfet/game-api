package admin

type GameConfigListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`

	GameID uint32 `form:"game_id"`

	ConfigKey string `form:"config_key"`
}