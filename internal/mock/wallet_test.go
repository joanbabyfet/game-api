package mock

import (
	"testing"
)

func TestWalletBetIsIdempotent(t *testing.T) {
	wallet := NewWallet()
	first, err := wallet.Bet("app", "chris1", "order-1", "slot", 10)
	if err != nil {
		t.Fatalf("first Bet() error = %v", err)
	}
	second, err := wallet.Bet("app", "chris1", "order-1", "slot", 10)
	if err != nil {
		t.Fatalf("second Bet() error = %v", err)
	}
	if first != 4990 || second != first || wallet.Balance("chris1") != 4990 {
		t.Fatalf("bet was not idempotent: first=%v second=%v balance=%v", first, second, wallet.Balance("chris1"))
	}
}

func TestWalletSettleIsIdempotent(t *testing.T) {
	wallet := NewWallet()
	if _, err := wallet.Bet("app", "chris1", "order-1", "slot", 10); err != nil {
		t.Fatal(err)
	}
	first, err := wallet.Settle("app", "chris1", "order-1", "slot", 20)
	if err != nil {
		t.Fatalf("first Settle() error = %v", err)
	}
	second, err := wallet.Settle("app", "chris1", "order-1", "slot", 20)
	if err != nil {
		t.Fatalf("second Settle() error = %v", err)
	}
	if first != 5010 || second != first || wallet.Balance("chris1") != 5010 {
		t.Fatalf("settle was not idempotent: first=%v second=%v balance=%v", first, second, wallet.Balance("chris1"))
	}
	order := wallet.orders[operatorOrderKey("app", "order-1")]
	if order == nil || order.Status != OperatorOrderStatusSettled || order.WinAmount != 20 {
		t.Fatalf("unexpected order: %+v", order)
	}
}

func TestWalletRollbackMissingOrderIsSuccess(t *testing.T) {
	wallet := NewWallet()
	balance, err := wallet.Rollback("app", "chris1", "missing", "slot", 10)
	if err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}
	if balance != 5000 || wallet.Balance("chris1") != 5000 {
		t.Fatalf("missing order rollback changed balance: returned=%v balance=%v", balance, wallet.Balance("chris1"))
	}
}

func TestWalletRollbackIsIdempotent(t *testing.T) {
	wallet := NewWallet()
	if _, err := wallet.Bet("app", "chris1", "order-1", "slot", 10); err != nil {
		t.Fatal(err)
	}
	first, err := wallet.Rollback("app", "chris1", "order-1", "slot", 10)
	if err != nil {
		t.Fatalf("first Rollback() error = %v", err)
	}
	second, err := wallet.Rollback("app", "chris1", "order-1", "slot", 10)
	if err != nil {
		t.Fatalf("second Rollback() error = %v", err)
	}
	if first != 5000 || second != first || wallet.Balance("chris1") != 5000 {
		t.Fatalf("rollback was not idempotent: first=%v second=%v balance=%v", first, second, wallet.Balance("chris1"))
	}
}
