package service

import (
	"context"
	"game-api/internal/model"
	"testing"
)

func TestRecoverOrderIgnoresTerminalStatus(t *testing.T) {
	s := &WalletService{}
	statuses := []int8{
		model.OrderStatusSettled,
		model.OrderStatusRolledBack,
		model.OrderStatusFailed,
		model.OrderStatusPending,
	}
	for _, status := range statuses {
		if err := s.RecoverOrder(context.Background(), &model.GameOrder{Status: status}); err != nil {
			t.Fatalf("RecoverOrder(status=%d) error = %v", status, err)
		}
	}
}
