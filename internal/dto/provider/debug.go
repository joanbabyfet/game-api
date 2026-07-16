package provider

type DebugSignReq struct {
	Fields map[string]any 	`json:"fields" binding:"required"`
}

type DebugSignResp struct {
	Data string `json:"data"`
	Sign string `json:"sign"`
}