package mock

const (
	OperatorOrderStatusBet        int8 = 1
	OperatorOrderStatusSettled    int8 = 2
	OperatorOrderStatusRolledBack int8 = 3
	OperatorOrderStatusFailed     int8 = 4
)

type OperatorOrder struct {
	AppID           string
	OrderNo         string
	PlayerID        string
	GameCode        string
	Status          int8
	BetAmount       float64
	WinAmount       float64
	BalanceBefore   float64
	BetBalanceAfter float64
	BalanceAfter    float64
}
