package pkg

// AppError 带错误码的业务错误
type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string {
	return e.Msg
}

// NewError 创建业务错误
func NewError(code int, msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  msg,
	}
}

//预定义错误
var (
	ErrInvalidParam 		= NewError(INVALID_PARAM, "invalid parameter")
	ErrUnauthorized 		= NewError(UNAUTHORIZED, "unauthorized")
	ErrForbidden    		= NewError(FORBIDDEN, "forbidden")
	ErrPlayerNotOnline		= NewError(PLAYER_NOT_ONLINE, "player not online")
	ErrGameNotFound			= NewError(GAME_NOT_FOUND, "game not found")
	ErrAgentNotFound		= NewError(AGENT_NOT_FOUND, "agent not found")
	ErrUserNotFound			= NewError(USER_NOT_FOUND, "user not found")
	ErrGameDisabled			= NewError(GAME_DISABLED, "game not found")
	ErrInsufficientBalance 	= NewError(INSUFFICIENT_BALANCE, "insufficient balance")
)