package operator

import (
	"context"
	mock "game-api/internal/dto/mock"
	"game-api/internal/model"
	"game-api/pkg"
	"log"
)

//查玩家余额
func (c *Client) Balance(ctx context.Context, baseURL string, agent *model.Agent, playerID string) (*mock.BalanceResp, error) {

	//组装 sign 字段
	now := pkg.Timestamp()
	sign := pkg.GenerateSign(
		pkg.BuildSignData(map[string]any{
			"app_id":    agent.AppID,
			"timestamp": now,
			"player_id": playerID,
		}),
		agent.SecretKey,
	)

	req := &mock.BalanceReq{
		BaseReq: mock.BaseReq{
			AppID:     agent.AppID,
			Timestamp: now,
			Sign:      sign,
		},
		PlayerID: playerID,
	}

	var resp mock.BalanceResp
	//调用运营商api
	if err := c.post(ctx, baseURL+"/balance", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//下注(扣钱)
func (c *Client) Bet(ctx context.Context, baseURL string, agent *model.Agent, playerID string, orderNo string, roundID string, gameCode string, betAmount float64) (*mock.BetResp, error) {

	// 组装 sign
	now := pkg.Timestamp()
	sign := pkg.GenerateSign(
		pkg.BuildSignData(map[string]any{
			"app_id":    agent.AppID,
			"timestamp": now,
			"player_id": playerID,
			"order_no":  orderNo,
			"round_id":  roundID,
			"game_code": gameCode,
			"bet_amount": betAmount,
		}),
		agent.SecretKey,
	)

	req := &mock.BetReq{
		BaseReq: mock.BaseReq{
			AppID:     agent.AppID,
			Timestamp: now,
			Sign:      sign,
		},
		PlayerID: playerID,
		OrderNo:  orderNo,
		RoundID:  roundID,
		GameCode: gameCode,
		BetAmount:  betAmount, //下注金额
	}

	log.Printf("====== Operator Bet Request ======")
    log.Printf("URL: %s", baseURL+"/bet")
    log.Printf("Request: %+v", req)

	var resp mock.BetResp
	// 调用 Operator API
	if err := c.post(ctx, baseURL+"/bet", req, &resp); err != nil {
		return nil, err
	}

	log.Printf("====== Operator Bet Response ======")
    log.Printf("Response: %+v", resp)

	return &resp, nil
}

//结算(加钱)
func (c *Client) Settle(ctx context.Context, baseURL string, agent *model.Agent, playerID string, orderNo string, roundID string, gameCode string, winAmount float64) (*mock.SettleResp, error) {

	// 组装 sign
	now := pkg.Timestamp()
	sign := pkg.GenerateSign(
		pkg.BuildSignData(map[string]any{
			"app_id":    agent.AppID,
			"timestamp": now,
			"player_id": playerID,
			"order_no":  orderNo,
			"round_id":  roundID,
			"game_code": gameCode,
			"win_amount": winAmount,
		}),
		agent.SecretKey,
	)

	req := &mock.SettleReq{
		BaseReq: mock.BaseReq{
			AppID:     agent.AppID,
			Timestamp: now,
			Sign:      sign,
		},
		PlayerID: playerID,
		OrderNo:  orderNo,
		RoundID:  roundID,
		GameCode: gameCode,
		WinAmount: winAmount,
	}

	log.Printf("====== Operator Settle Request ======")
    log.Printf("URL: %s", baseURL+"/settle")
    log.Printf("Request: %+v", req)

	var resp mock.SettleResp

	// 调用 Operator API
	if err := c.post(ctx, baseURL+"/settle", req, &resp); err != nil {
		return nil, err
	}

	log.Printf("====== Operator Settle Response ======")
    log.Printf("Response: %+v", resp)

	return &resp, nil
}

// 取消下注
func (c *Client) Rollback(ctx context.Context, baseURL string, agent *model.Agent, playerID string, orderNo string, roundID string, gameCode string, betAmount float64) (*mock.RollbackResp, error) {

	// 组装 sign
	now := pkg.Timestamp()
	sign := pkg.GenerateSign(
		pkg.BuildSignData(map[string]any{
			"app_id":    agent.AppID,
			"timestamp": now,
			"player_id": playerID,
			"order_no":  orderNo,
			"round_id":  roundID,
			"game_code": gameCode,
			"bet_amount": betAmount,
		}),
		agent.SecretKey,
	)

	req := &mock.RollbackReq{
		BaseReq: mock.BaseReq{
			AppID:     agent.AppID,
			Timestamp: now,
			Sign:      sign,
		},
		PlayerID: playerID,
		OrderNo:  orderNo,
		RoundID:  roundID,
		GameCode: gameCode,
		BetAmount:  betAmount, //下注金额
	}

	log.Printf("====== Operator Rollback Request ======")
    log.Printf("URL: %s", baseURL+"/rollback")
    log.Printf("Request: %+v", req)

	var resp mock.RollbackResp

	// 调用 Operator API
	if err := c.post(ctx, baseURL+"/rollback", req, &resp); err != nil {
		return nil, err
	}

	log.Printf("====== Operator Rollback Response ======")
    log.Printf("Response: %+v", resp)

	return &resp, nil
}