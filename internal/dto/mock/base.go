package mock

type BaseReq struct {
    AppID     string `json:"app_id" form:"app_id" binding:"required"`      // Operator App ID
    Timestamp int64  `json:"timestamp" form:"timestamp" binding:"required"`
    Sign      string `json:"sign" form:"sign" binding:"required"`
}