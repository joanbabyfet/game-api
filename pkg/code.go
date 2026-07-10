package pkg

//统一错误码
const (
	// 成功
	SUCCESS = 0

	// 系统错误
	UNKNOWN_ERROR = -1

	// 通用错误
	INVALID_PARAM = -1001
	UNAUTHORIZED  = -1002
	FORBIDDEN     = -1003

	ORDER_NOT_FOUND      = -2001
    INSUFFICIENT_BALANCE = -2002

	PLAYER_NOT_ONLINE	= -2003
)