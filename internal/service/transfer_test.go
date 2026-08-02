package service

import (
	"errors"
	providerdto "game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/pkg"
	"testing"
)

func TestReplayTransferSuccess(t *testing.T) {
	req := &providerdto.TransferReq{
		RequestID: "request-1", ThirdOrderNo: "third-1", Currency: "usd", Amount: 10,
	}
	order := &model.WalletTransfer{
		RequestID: "request-1", ThirdOrderNo: "third-1", UID: 100, GameID: 1,
		TransferType: model.GameTransferTypeIn, Amount: 1000, Currency: "USD",
		BalanceAfter: 5000, Status: model.GameTransferStatusSuccess,
	}

	resp, err := replayTransfer(order, req, 100, 1, model.GameTransferTypeIn, 1000)
	if err != nil {
		t.Fatalf("replayTransfer() error = %v", err)
	}
	if resp.Balance != 50 || resp.Amount != 10 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestReplayTransferConflict(t *testing.T) {
	req := &providerdto.TransferReq{RequestID: "request-1", ThirdOrderNo: "third-1", Currency: "USD"}
	order := &model.WalletTransfer{
		RequestID: "request-1", ThirdOrderNo: "third-1", UID: 100, GameID: 1,
		TransferType: model.GameTransferTypeIn, Amount: 1000, Currency: "USD",
		Status: model.GameTransferStatusSuccess,
	}

	_, err := replayTransfer(order, req, 100, 1, model.GameTransferTypeIn, 2000)
	if !errors.Is(err, pkg.ErrTransferOrderConflict) {
		t.Fatalf("expected transfer conflict, got %v", err)
	}
}

func TestReplayTransferFailedReturnsOriginalError(t *testing.T) {
	req := &providerdto.TransferReq{RequestID: "request-1", ThirdOrderNo: "third-1", Currency: "USD"}
	order := &model.WalletTransfer{
		RequestID: "request-1", ThirdOrderNo: "third-1", UID: 100, GameID: 1,
		TransferType: model.GameTransferTypeOut, Amount: 1000, Currency: "USD",
		Status: model.GameTransferStatusFailed, ErrorCode: pkg.INSUFFICIENT_BALANCE,
		ErrorMessage: "insufficient balance",
	}

	_, err := replayTransfer(order, req, 100, 1, model.GameTransferTypeOut, 1000)
	var appErr *pkg.AppError
	if !errors.As(err, &appErr) || appErr.Code != pkg.INSUFFICIENT_BALANCE {
		t.Fatalf("expected original insufficient balance error, got %v", err)
	}
}
