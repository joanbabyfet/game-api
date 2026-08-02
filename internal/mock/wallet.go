package mock

import (
	"game-api/pkg"
	"sync"
)

type Wallet struct {
	mu       sync.RWMutex
	balances map[string]float64
	orders   map[string]*OperatorOrder
}

func NewWallet() *Wallet {
	return &Wallet{
		balances: map[string]float64{
			"chris1": 5000.00,
		},
		orders: make(map[string]*OperatorOrder),
	}
}

// 查询余额
func (w *Wallet) Balance(player_id string) float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.balances[player_id]
}

// 下注扣款
func (w *Wallet) Bet(appID, playerID, orderNo, gameCode string, amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	key := operatorOrderKey(appID, orderNo)
	if order, ok := w.orders[key]; ok {
		if order.PlayerID != playerID || order.GameCode != gameCode || order.BetAmount != amount {
			return 0, pkg.ErrOperatorOrderConflict
		}
		switch order.Status {
		case OperatorOrderStatusBet, OperatorOrderStatusSettled:
			return order.BetBalanceAfter, nil
		case OperatorOrderStatusRolledBack:
			return order.BalanceAfter, pkg.ErrOperatorOrderInvalid
		default:
			return 0, pkg.ErrOperatorOrderInvalid
		}
	}

	balance, ok := w.balances[playerID]
	if !ok {
		return 0, pkg.ErrUserNotFound
	}

	if balance < amount {
		return balance, pkg.ErrInsufficientBalance
	}

	balance -= amount
	w.balances[playerID] = balance
	w.orders[key] = &OperatorOrder{
		AppID: appID, OrderNo: orderNo, PlayerID: playerID, GameCode: gameCode,
		Status: OperatorOrderStatusBet, BetAmount: amount,
		BalanceBefore: balance + amount, BetBalanceAfter: balance, BalanceAfter: balance,
	}

	return balance, nil
}

// 派彩
func (w *Wallet) Settle(appID, playerID, orderNo, gameCode string, amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	order, ok := w.orders[operatorOrderKey(appID, orderNo)]
	if !ok {
		return 0, pkg.ErrOperatorOrderNotFound
	}
	if order.PlayerID != playerID || order.GameCode != gameCode {
		return 0, pkg.ErrOperatorOrderConflict
	}
	if order.Status == OperatorOrderStatusSettled {
		if order.WinAmount != amount {
			return 0, pkg.ErrOperatorOrderConflict
		}
		return order.BalanceAfter, nil
	}
	if order.Status != OperatorOrderStatusBet {
		return 0, pkg.ErrOperatorOrderInvalid
	}

	balance := w.balances[playerID]
	balance += amount
	w.balances[playerID] = balance
	order.Status = OperatorOrderStatusSettled
	order.WinAmount = amount
	order.BalanceAfter = balance

	return balance, nil
}

// 下注并结算
func (w *Wallet) BetSettle(player_id string, bet, win float64) (float64, error) {

	_, err := w.Bet("legacy", player_id, "legacy", "legacy", bet)
	if err != nil {
		return 0, err
	}

	return w.Settle("legacy", player_id, "legacy", "legacy", win)
}

// 回滚
func (w *Wallet) Rollback(appID, playerID, orderNo, gameCode string, amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	order, ok := w.orders[operatorOrderKey(appID, orderNo)]
	if !ok {
		balance, playerExists := w.balances[playerID]
		if !playerExists {
			return 0, pkg.ErrUserNotFound
		}
		// V2.8 cancel semantics: a missing/failed bet is a successful no-op.
		return balance, nil
	}
	if order.PlayerID != playerID || order.GameCode != gameCode || order.BetAmount != amount {
		return 0, pkg.ErrOperatorOrderConflict
	}
	if order.Status == OperatorOrderStatusRolledBack {
		return order.BalanceAfter, nil
	}
	if order.Status != OperatorOrderStatusBet {
		return 0, pkg.ErrOperatorOrderInvalid
	}

	balance := w.balances[playerID]
	balance += amount
	w.balances[playerID] = balance
	order.Status = OperatorOrderStatusRolledBack
	order.BalanceAfter = balance

	return balance, nil
}

func operatorOrderKey(appID, orderNo string) string {
	return appID + ":" + orderNo
}

// 设置余额（测试用）
func (w *Wallet) SetBalance(player_id string, balance float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.balances[player_id] = balance
}

// 重置所有玩家余额（测试用）
func (w *Wallet) Reset(balance float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for player_id := range w.balances {
		w.balances[player_id] = balance
	}
}
