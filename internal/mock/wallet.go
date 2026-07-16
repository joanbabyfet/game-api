package mock

import (
	"game-api/pkg"
	"sync"
)

type Wallet struct {
	mu       sync.RWMutex
	balances map[string]float64
}

func NewWallet() *Wallet {
	return &Wallet{
		balances: map[string]float64{
			"chris1": 1000.00,
			"chris2": 1000.00,
			"chris3": 1000.00,
		},
	}
}

// 查询余额
func (w *Wallet) Balance(player_id string) float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.balances[player_id]
}

// 下注扣款
func (w *Wallet) Bet(player_id string, amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	balance, ok := w.balances[player_id]
	if !ok {
		return 0, pkg.ErrUserNotFound
	}

	if balance < amount {
		return balance, pkg.ErrInsufficientBalance
	}

	balance -= amount
	w.balances[player_id] = balance

	return balance, nil
}

// 派彩
func (w *Wallet) Settle(player_id string, amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	balance, ok := w.balances[player_id]
	if !ok {
		return 0, pkg.ErrUserNotFound
	}

	balance += amount
	w.balances[player_id] = balance

	return balance, nil
}

// 下注并结算
func (w *Wallet) BetSettle(player_id string, bet, win float64) (float64, error) {

    _, err := w.Bet(player_id, bet)
    if err != nil {
        return 0, err
    }

    return w.Settle(player_id, win)
}

// 回滚
func (w *Wallet) Rollback(player_id string, amount float64) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	balance, ok := w.balances[player_id]
	if !ok {
		return 0, pkg.ErrUserNotFound
	}

	balance += amount
	w.balances[player_id] = balance

	return balance, nil
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