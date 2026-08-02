package service

import (
	"errors"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/pkg"
	"testing"
)

func TestValidateSpinReplay(t *testing.T) {
	claims := &pkg.JWTClaims{UID: 100, AgentID: 1}
	game := &model.Game{ID: 7}
	req := &provider.SpinReq{RequestID: "request-1", GameCode: "slot", BetAmount: 10}
	order := &model.GameOrder{
		RequestID: "request-1", UID: 100, AgentID: 1, WalletMode: model.WalletModeTransfer, GameID: 7,
		SpinType: model.SpinTypeNormal, BetAmount: 1000,
	}

	if err := validateSpinReplay(order, req, claims, game, model.SpinTypeNormal, model.WalletModeTransfer); err != nil {
		t.Fatalf("validateSpinReplay() error = %v", err)
	}
}

func TestValidateSpinReplayRejectsBetConflict(t *testing.T) {
	claims := &pkg.JWTClaims{UID: 100, AgentID: 1}
	game := &model.Game{ID: 7}
	req := &provider.SpinReq{RequestID: "request-1", BetAmount: 20}
	order := &model.GameOrder{
		RequestID: "request-1", UID: 100, AgentID: 1, WalletMode: model.WalletModeTransfer, GameID: 7,
		SpinType: model.SpinTypeNormal, BetAmount: 1000,
	}

	err := validateSpinReplay(order, req, claims, game, model.SpinTypeNormal, model.WalletModeTransfer)
	var appErr *pkg.AppError
	if !errors.As(err, &appErr) || appErr.Code != pkg.ORDER_STATUS_ERROR {
		t.Fatalf("expected request conflict, got %v", err)
	}
}

func TestValidateFreeSpinReplay(t *testing.T) {
	claims := &pkg.JWTClaims{UID: 100, AgentID: 1}
	game := &model.Game{ID: 7}
	req := &provider.SpinReq{RequestID: "request-1", FreeSpinID: "FS-1"}
	order := &model.GameOrder{
		RequestID: "request-1", UID: 100, AgentID: 1, WalletMode: model.WalletModeTransfer, GameID: 7,
		SpinType: model.SpinTypeFreeSpin, FreeSpinID: "FS-1", BetAmount: 1000,
	}

	if err := validateSpinReplay(order, req, claims, game, model.SpinTypeFreeSpin, model.WalletModeTransfer); err != nil {
		t.Fatalf("validateSpinReplay() error = %v", err)
	}
}

func TestValidateSpinReplayRejectsWalletModeConflict(t *testing.T) {
	claims := &pkg.JWTClaims{UID: 100, AgentID: 1}
	game := &model.Game{ID: 7}
	req := &provider.SpinReq{RequestID: "request-1", BetAmount: 10}
	order := &model.GameOrder{
		RequestID: "request-1", UID: 100, AgentID: 1,
		WalletMode: model.WalletModeSingle, GameID: 7,
		SpinType: model.SpinTypeNormal, BetAmount: 1000,
	}

	err := validateSpinReplay(order, req, claims, game, model.SpinTypeNormal, model.WalletModeTransfer)
	var appErr *pkg.AppError
	if !errors.As(err, &appErr) || appErr.Code != pkg.ORDER_STATUS_ERROR {
		t.Fatalf("expected wallet mode conflict, got %v", err)
	}
}
