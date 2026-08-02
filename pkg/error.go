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

// 预定义错误
var (
	ErrInvalidParam            = NewError(INVALID_PARAM, "invalid parameter")
	ErrUnauthorized            = NewError(UNAUTHORIZED, "unauthorized")
	ErrForbidden               = NewError(FORBIDDEN, "forbidden")
	ErrPlayerNotOnline         = NewError(PLAYER_NOT_ONLINE, "player not online")
	ErrGameNotFound            = NewError(GAME_NOT_FOUND, "game not found")
	ErrAgentNotFound           = NewError(AGENT_NOT_FOUND, "agent not found")
	ErrUserNotFound            = NewError(USER_NOT_FOUND, "user not found")
	ErrWalletNotFound          = NewError(WALLET_NOT_FOUND, "wallet not found")
	ErrOrderNotFound           = NewError(ORDER_NOT_FOUND, "order not found")
	ErrOrderStatusInvalid      = NewError(ORDER_STATUS_INVALID, "order status invalid")
	ErrOrderProcessing         = NewError(ORDER_PROCESSING, "game order processing")
	ErrOrderUpdateFailed       = NewError(ORDER_UPDATE_FAILED, "order update failed")
	ErrOrderStatus             = NewError(ORDER_STATUS_ERROR, "order status error")
	ErrWalletAddFailed         = NewError(WALLET_ADD_FAILED, "wallet add failed")
	ErrGameDisabled            = NewError(GAME_DISABLED, "game not found")
	ErrInsufficientBalance     = NewError(INSUFFICIENT_BALANCE, "insufficient balance")
	ErrTransferOrderNotFound   = NewError(TRANSFER_ORDER_NOT_FOUND, "transfer order not found")
	ErrTransferOrderProcessing = NewError(TRANSFER_ORDER_PROCESSING, "transfer order processing")
	ErrTransferOrderFailed     = NewError(TRANSFER_ORDER_FAILED, "transfer order failed")
	ErrTransferOrderConflict   = NewError(TRANSFER_ORDER_CONFLICT, "transfer order conflict")
	ErrTransferInFailed        = NewError(TRANSFER_IN_FAILED, "transfer in failed")
	ErrTransferOutFailed       = NewError(TRANSFER_OUT_FAILED, "transfer out failed")
	ErrWalletModeInvalid       = NewError(WALLET_MODE_INVALID, "wallet mode invalid")
	ErrCurrencyMismatch        = NewError(CURRENCY_MISMATCH, "currency mismatch")
	ErrOperatorOrderNotFound   = NewError(OPERATOR_ORDER_NOT_FOUND, "operator order not found")
	ErrOperatorOrderConflict   = NewError(OPERATOR_ORDER_CONFLICT, "operator order conflict")
	ErrOperatorOrderInvalid    = NewError(OPERATOR_ORDER_INVALID, "operator order status invalid")
)
