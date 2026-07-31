package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game-api/internal/adapter"
	"game-api/internal/client/operator"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
	"game-api/proto/slotpb"
	"log"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const walletCacheTTL = 24 * time.Hour

func walletCacheKey(uid uint64) string {
	return fmt.Sprintf("wallet:%d", uid)
}

type WalletService struct {
	db              *gorm.DB
	redisClient    *redis.Client
	repo *repository.WalletRepository
	agentRepo *repository.AgentRepository
	gameRepo *repository.GameRepository
	userRepo *repository.UserRepository
	walletRepo *repository.WalletRepository
	walletLogRepo *repository.WalletLogRepository
	orderRepo *repository.GameOrderRepository
	rollbackLogRepo *repository.RollbackLogRepository
	adapter *adapter.WalletAdapter
	slotAdapter *adapter.SlotAdapter
	authService *AuthService
	operatorClient *operator.Client
}

func NewWalletService(
	db *gorm.DB,
	redisClient    *redis.Client,
	repo *repository.WalletRepository,
	agentRepo *repository.AgentRepository,
	gameRepo *repository.GameRepository,
	userRepo *repository.UserRepository,
	walletRepo *repository.WalletRepository,
	walletLogRepo *repository.WalletLogRepository,
	orderRepo *repository.GameOrderRepository,
	rollbackLogRepo *repository.RollbackLogRepository,
	adapter *adapter.WalletAdapter,
	slotAdapter *adapter.SlotAdapter,
	authService *AuthService,
	operatorClient *operator.Client,
) *WalletService {
	return &WalletService{
		db: db,
		redisClient:    redisClient,
		repo: repo,
		agentRepo: agentRepo,
		gameRepo: gameRepo,
		userRepo: userRepo,
		walletRepo: walletRepo,
		walletLogRepo: walletLogRepo,
		orderRepo: orderRepo,
		rollbackLogRepo: rollbackLogRepo,
		adapter: adapter,
		slotAdapter: slotAdapter,
		authService: authService,
		operatorClient: operatorClient,
	}
}

// Create 新增钱包
func (s *WalletService) Create(ctx context.Context, wallet *model.Wallet) error {
	return s.repo.Create(ctx, wallet)
}

// Update 更新钱包
func (s *WalletService) Update(ctx context.Context, wallet *model.Wallet) error {
	return s.repo.Update(ctx, wallet)
}

// Delete 删除钱包
func (s *WalletService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// GetByUID 根据UID查询
func (s *WalletService) GetByUID(ctx context.Context, uid uint64) (*model.Wallet, error) {
	return s.repo.GetByUID(ctx, uid)
}

// List 钱包列表
func (s *WalletService) List(ctx context.Context, q repository.WalletQuery) ([]model.Wallet, error) {
	return s.repo.List(ctx, q)
}

// 查询玩家余额(单一钱包下，余额来源就是 Operator)
func (s *WalletService) Balance(ctx context.Context, req *provider.BalanceReq) (*provider.BalanceResp, error) {

	//解析 JWT (客户端已经登录过，不需要再验 app_id + sign)
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	//获取代理信息
	agent, err := s.agentRepo.GetByID(claims.AgentID)
	if err != nil {
		return nil, err
	}

	//调用 Operator
	resp, err := s.operatorClient.Balance(ctx, agent.OperatorURL, agent, claims.PlayerID)
	if err != nil {
		return nil, err
	}

	// DTO 轉換
	return &provider.BalanceResp{
		Balance: resp.Balance,
		Currency: resp.Currency,
	}, nil
}

// Spin (每点击一次 Spin 调用)
func (s *WalletService) Spin(ctx context.Context, req *provider.SpinReq) (*provider.SpinResp, error) {
	
	// 1. 解析 JWT
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		log.Printf("ParseToken error: %+v", err)
		return nil, pkg.ErrUnauthorized
	}

	// 2. 查询游戏
	game, err := s.gameRepo.GetByCode(req.GameCode)
	if err != nil {
		log.Printf("GetByCode error: %+v", err)
		return nil, pkg.ErrGameNotFound
	}

	// 3. 查询代理信息
	agent, err := s.agentRepo.GetByID(claims.AgentID)
	if err != nil {
		return nil, err
	}

	// 4. request_id 幂等检查
	order, err := s.orderRepo.GetByRequestID(ctx, req.RequestID)
	if err == nil {
		// 已存在订单，直接回放，不重新下注
		return s.replaySpin(order)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 查询数据库异常
		return nil, err
	}

	//这些由入口 Provider API 生成
	requestID := req.RequestID //幂等(由 Provider API 传入)
	orderNo := pkg.GenOrderNo() //注單号

	//5. 第一次请求，开始创建订单
	order = &model.GameOrder{
		RequestID:  requestID,
		OrderNo:    orderNo,
		UID:        claims.UID,
		AgentID:    claims.AgentID,
		GameID:     game.ID,
		BetAmount:  pkg.ToMoney(req.BetAmount),
		WinAmount:  0,
		Profit:     0,
		Currency:   claims.Currency,
		Status:     model.OrderStatusPending, //Provider 已创建 game_order，Operator 下注结果尚未确认
		CreateTime: time.Now().Unix(),
	}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	log.Printf("========== Provider Spin Start ==========")
	log.Printf("RequestID : %s", requestID)
	log.Printf("OrderNo   : %s", orderNo)
	log.Printf("PlayerID  : %s", claims.PlayerID)
	log.Printf("UID       : %d", claims.UID)
	log.Printf("AgentID   : %d", claims.AgentID)
	log.Printf("GameID    : %d", game.ID)
	log.Printf("GameCode  : %s", game.GameCode)
	log.Printf("BetAmount : %v", req.BetAmount)
	log.Printf("Currency  : %s", claims.Currency)

	// 6. Operator 扣款
	//
	// 单一钱包模式下：
	// Provider 不直接修改本地钱包余额，
	// 而是调用 Operator 的下注接口扣除玩家余额。
	log.Printf(
		"Operator Bet start: request_id=%s order_no=%s player=%s amount=%v",
		requestID,
		orderNo,
		claims.PlayerID,
		req.BetAmount,
	)
	betResp, err := s.operatorClient.Bet(
		ctx,
		agent.OperatorURL,
		agent,
		claims.PlayerID,
		orderNo,
		game.GameCode,
		req.BetAmount,
	)
	if err != nil {
		log.Printf(
			"Operator Bet error: request_id=%s order_no=%s player=%s err=%+v",
			requestID,
			orderNo,
			claims.PlayerID,
			err,
		)
		return nil, err
	}

	if betResp == nil {
		log.Printf(
			"Operator Bet returned nil response: request_id=%s order_no=%s",
			requestID,
			orderNo,
		)

		return nil, pkg.NewError(
			pkg.OPERATOR_BET_FAILED,
			"operator bet returned nil response",
		)
	}

	log.Printf(
		"Operator Bet success: request_id=%s order_no=%s resp=%+v",
		requestID,
		orderNo,
		betResp,
	)

	//Bet 成功后更新状态
	if err := s.orderRepo.UpdateStatus(
		ctx,
		order.ID,
		model.OrderStatusPending,
		model.OrderStatusBetSuccess,
	); err != nil {
		log.Printf(
			"Update order status failed: request_id=%s order_no=%s from=%d to=%d err=%+v",
			requestID,
			orderNo,
			model.OrderStatusPending,
			model.OrderStatusBetSuccess,
			err,
		)

		return nil, pkg.ErrOrderUpdateFailed
	}

	order.Status = model.OrderStatusBetSuccess

	// 7. 调用 Skynet 执行游戏逻辑
	pbReq := &slotpb.SpinReq{
		RequestId: requestID,
		OrderNo:   orderNo,
		Uid:       claims.UID,
		AgentId:   claims.AgentID,
		GameId:    game.ID,
		BetAmount: pkg.ToMoney(req.BetAmount),
		Currency:  claims.Currency,
		DebugFail: req.DebugFail, //测试取消下注用
	}

	log.Printf("========== Skynet Spin Request ==========")
	log.Printf("RequestID : %s", pbReq.RequestId)
	log.Printf("OrderNo   : %s", pbReq.OrderNo)
	log.Printf("UID       : %d", pbReq.Uid)
	log.Printf("AgentID   : %d", pbReq.AgentId)
	log.Printf("GameID    : %d", pbReq.GameId)
	log.Printf("BetAmount : %d", pbReq.BetAmount)
	log.Printf("Currency  : %s", pbReq.Currency)
	log.Printf("DebugFail : %v", pbReq.DebugFail)

	spinResp, spinErr := s.slotAdapter.Spin(ctx, pbReq)
	if spinErr != nil {
		log.Printf(
			"Skynet Spin error: request_id=%s order_no=%s err=%+v",
			requestID,
			orderNo,
			spinErr,
		)

		// Skynet 明确失败，订单进入待回滚状态
		if err := s.orderRepo.UpdateStatus(
			ctx,
			order.ID,
			model.OrderStatusBetSuccess,
			model.OrderStatusWaitRollback,
		); err != nil {
			log.Printf(
				"Update WaitRollback failed: request_id=%s order_no=%s err=%+v",
				requestID,
				orderNo,
				err,
			)

			return nil, pkg.ErrOrderUpdateFailed
		}

		order.Status = model.OrderStatusWaitRollback

		// Skynet 失败，回滚扣款
		_, rollbackErr := s.operatorClient.Rollback(
			ctx,
			agent.OperatorURL,
			agent,
			claims.PlayerID,
			orderNo,
			game.GameCode,
			req.BetAmount,
		)
		if rollbackErr != nil {
			log.Printf(
				"Operator rollback failed, order=%s player=%s err=%v",
				orderNo,
				claims.PlayerID,
				rollbackErr,
			)

			// 状态保持 WaitRollback，Worker 后续重试
			return nil, spinErr
		}

		// Operator 回滚成功
		if err := s.orderRepo.UpdateStatus(
			ctx,
			order.ID,
			model.OrderStatusWaitRollback,
			model.OrderStatusRolledBack,
		); err != nil {
			log.Printf(
				"Update RolledBack failed: request_id=%s order_no=%s err=%+v",
				requestID,
				orderNo,
				err,
			)

			return nil, pkg.ErrOrderUpdateFailed
		}

		order.Status = model.OrderStatusRolledBack

		return nil, spinErr
	}

	if spinResp == nil {
		log.Printf(
			"Skynet Spin returned nil response: request_id=%s order_no=%s",
			requestID,
			orderNo,
		)

		return nil, pkg.NewError(
			pkg.SPIN_FAILED,
			"skynet spin returned nil response",
		)
	}

	// 确认 spinResp 不为 nil 后，才访问里面的字段
	log.Printf("========== Skynet Spin Response ==========")
	log.Printf("RequestID : %s", requestID)
	log.Printf("OrderNo   : %s", orderNo)
	log.Printf("WinAmount : %d", spinResp.WinAmount)
	log.Printf("RawResp   : %+v", spinResp)

	if spinResp.WinAmount < 0 {
		log.Printf(
			"Invalid win amount: request_id=%s order_no=%s win_amount=%d",
			requestID,
			orderNo,
			spinResp.WinAmount,
		)
		return nil, pkg.ErrInvalidParam
	}

	// 玩家净输赢： 中奖金额 - 下注金额
	profit := spinResp.WinAmount - order.BetAmount

	// 保存 Skynet 游戏结果，并进入待派奖状态
	if err := s.orderRepo.UpdateGameResult(
		ctx,
		order.ID,
		spinResp.RoundId,
		spinResp.WinAmount,
		profit,
	); err != nil {
		log.Printf(
			"Update game result failed: request_id=%s order_no=%s round_id=%s err=%+v",
			requestID,
			orderNo,
			spinResp.RoundId,
			err,
		)

		// Skynet 已经产生结果，不能重新执行 Spin
		return nil, pkg.ErrOrderUpdateFailed
	}

	order.RoundID = spinResp.RoundId
	order.WinAmount = spinResp.WinAmount
	order.Profit = profit
	order.Status = model.OrderStatusWaitSettle

	winAmount := pkg.ToAmount(spinResp.WinAmount)

	// 8. Operator 派奖
	//
	// 单一钱包模式下，派奖也由 Operator 修改玩家余额。
	// orderNo 继续使用原下注订单号，作为该笔注单的关联标识。
	log.Printf(
		"Operator Settle start: request_id=%s order_no=%s player=%s win_amount=%v",
		requestID,
		orderNo,
		claims.PlayerID,
		winAmount,
	)
	settleResp, err := s.operatorClient.Settle(
		ctx,
		agent.OperatorURL,
		agent,
		claims.PlayerID,
		orderNo,
		game.GameCode,
		winAmount,
	)
	if err != nil {
		log.Printf(
			"Operator Settle error: request_id=%s order_no=%s player=%s win_amount=%v err=%+v",
			requestID,
			orderNo,
			claims.PlayerID,
			winAmount,
			err,
		)

		// 注意：
		// Skynet 游戏结果已经产生，不能再取消下注。
		//
		// 此时应将订单标记为 WAIT_SETTLE，
		// Worker 使用同一个 orderNo 和 winAmount 重试派彩。
		//
		// 不要重新执行 Spin，也不要重新生成订单号。
		return nil, err
	}

	if settleResp == nil {
		log.Printf(
			"Operator Settle returned nil response: request_id=%s order_no=%s",
			requestID,
			orderNo,
		)

		// 同样应该标记 WAIT_SETTLE，避免丢失派彩
		return nil, pkg.ErrInvalidParam
	}

	log.Printf("========== Operator Settle Success ==========")
	log.Printf("RequestID : %s", requestID)
	log.Printf("OrderNo   : %s", orderNo)
	log.Printf("Balance   : %v", settleResp.Balance)
	log.Printf("Currency  : %s", settleResp.Currency)

	// 更新订单为 Settled
	if err := s.orderRepo.UpdateSettled(
		ctx,
		order.ID,
		pkg.ToMoney(settleResp.Balance),
		settleResp.Currency,
	); err != nil {
		log.Printf(
			"Update Settled failed: request_id=%s order_no=%s err=%+v",
			requestID,
			orderNo,
			err,
		)

		return nil, pkg.ErrOrderUpdateFailed
	}

	order.BalanceAfter = pkg.ToMoney(settleResp.Balance)
	order.Currency = settleResp.Currency
	order.Status = model.OrderStatusSettled

	resp := &provider.SpinResp{
		// 单一钱包模式下，最终余额和币种以 Operator 返回为准
		Balance:  settleResp.Balance,
		Currency: settleResp.Currency,
	}

	log.Printf("========== Provider Spin Success ==========")
	log.Printf("RequestID : %s", requestID)
	log.Printf("OrderNo   : %s", orderNo)
	log.Printf("Response  : %+v", resp)

	return resp, nil
}

// getWallet 查询钱包，使用 Cache Aside
//
// 此方法只适合普通余额查询。
// 涉及加钱、扣钱、回滚时，必须直接查询数据库并加锁。
func (s *WalletService) getWallet(
	ctx context.Context,
	uid uint64,
) (*model.Wallet, error) {
	wallet, err := s.getWalletCache(ctx, uid)
	if err == nil {
		return wallet, nil
	}

	if !errors.Is(err, redis.Nil) {
		log.Printf(
			"[wallet] get cache failed uid=%d err=%v",
			uid,
			err,
		)
	}

	wallet, err = s.walletRepo.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 缓存失败不影响主流程
	if err := s.saveWalletCache(ctx, wallet); err != nil {
		log.Printf(
			"[wallet] set cache failed uid=%d err=%v",
			uid,
			err,
		)
	}

	return wallet, nil
}

// deleteCache 删除钱包缓存
//
// 删除失败只记录日志，不影响数据库事务结果。
func (s *WalletService) deleteCache(
	ctx context.Context,
	uid uint64,
) {
	if err := s.deleteWalletCache(ctx, uid); err != nil {
		log.Printf(
			"[wallet] delete cache failed uid=%d err=%v",
			uid,
			err,
		)
	}
}

// Add 钱包加钱
func (s *WalletService) Add(
	ctx context.Context,
	req *provider.ChangeBalanceReq,
) (int64, error) {
	if req.Amount <= 0 {
		return 0, pkg.ErrInvalidParam
	}

	var balance int64

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error

		balance, err = s.addWithTx(ctx, tx, req)
		return err
	})

	if err != nil {
		return 0, err
	}

	// 事务提交后删除缓存
	s.deleteCache(ctx, req.UID)

	return balance, nil
}

// Sub 钱包扣钱
func (s *WalletService) Sub(
	ctx context.Context,
	req *provider.ChangeBalanceReq,
) (int64, error) {
	if req.Amount <= 0 {
		return 0, pkg.ErrInvalidParam
	}

	var balance int64

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error

		balance, err = s.subWithTx(ctx, tx, req)
		return err
	})

	if err != nil {
		return 0, err
	}

	// 事务提交后删除缓存
	s.deleteCache(ctx, req.UID)

	return balance, nil
}

// addWithTx 在现有事务中增加余额
func (s *WalletService) addWithTx(
	ctx context.Context,
	tx *gorm.DB,
	req *provider.ChangeBalanceReq,
) (int64, error) {
	walletRepo := s.walletRepo.WithTx(tx)
	walletLogRepo := s.walletLogRepo.WithTx(tx)

	// 查询并锁定钱包
	wallet, err := walletRepo.GetByUIDForUpdate(ctx, req.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, pkg.ErrWalletNotFound
		}

		return 0, err
	}

	before := wallet.Balance
	after := before + req.Amount

	if err := walletRepo.AddBalance(
		ctx,
		req.UID,
		req.Amount,
	); err != nil {
		return 0, pkg.ErrWalletAddFailed
	}

	if err := walletLogRepo.Create(ctx, &model.WalletLog{
		UID:           req.UID,
		AgentID:       req.AgentID,
		GameID:        req.GameID,
		Type:          req.Type,
		Amount:        req.Amount,
		BalanceBefore: before,
		BalanceAfter:  after,
		RefOrderNo:    req.RefOrderNo,
	}); err != nil {
		return 0, err
	}

	return after, nil
}

// subWithTx 在现有事务中扣除余额
func (s *WalletService) subWithTx(
	ctx context.Context,
	tx *gorm.DB,
	req *provider.ChangeBalanceReq,
) (int64, error) {
	walletRepo := s.walletRepo.WithTx(tx)
	walletLogRepo := s.walletLogRepo.WithTx(tx)

	// 查询并锁定钱包
	wallet, err := walletRepo.GetByUIDForUpdate(ctx, req.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, pkg.ErrWalletNotFound
		}

		return 0, err
	}

	if wallet.Balance < req.Amount {
		return 0, pkg.ErrInsufficientBalance
	}

	before := wallet.Balance
	after := before - req.Amount

	if err := walletRepo.SubBalance(
		ctx,
		req.UID,
		req.Amount,
	); err != nil {
		if errors.Is(err, pkg.ErrInsufficientBalance) {
			return 0, pkg.ErrInsufficientBalance
		}

		return 0, err
	}

	if err := walletLogRepo.Create(ctx, &model.WalletLog{
		UID:           req.UID,
		AgentID:       req.AgentID,
		GameID:        req.GameID,
		Type:          req.Type,
		Amount:        req.Amount,
		BalanceBefore: before,
		BalanceAfter:  after,
		RefOrderNo:    req.RefOrderNo,
	}); err != nil {
		return 0, err
	}

	return after, nil
}

// Rollback 已结算或未结算订单回滚
func (s *WalletService) Rollback(
	ctx context.Context,
	req *provider.RollbackReq,
) (int64, error) {
	var (
		balance int64
		uid     uint64
	)

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		rollbackLogRepo := s.rollbackLogRepo.WithTx(tx)

		// 查询并锁定订单，防止并发重复回滚
		order, err := orderRepo.GetByOrderNoForUpdate(
			ctx,
			req.OrderNo,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pkg.ErrOrderNotFound
			}

			return err
		}

		uid = order.UID

		// 已回滚，直接返回当前余额
		if order.Status == model.OrderStatusRolledBack {
			walletRepo := s.walletRepo.WithTx(tx)

			wallet, err := walletRepo.GetByUIDForUpdate(
				ctx,
				order.UID,
			)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return pkg.ErrWalletNotFound
				}

				return err
			}

			balance = wallet.Balance
			return nil
		}

		// 只有待结算、已结算订单可以回滚
		if order.Status != model.OrderStatusPending &&
			order.Status != model.OrderStatusSettled {
			return pkg.ErrOrderStatusInvalid
		}

		var rollbackAmount int64

		switch order.Status {
		case model.OrderStatusPending:
			// 未结算：
			// Spin 时只扣除了 bet，因此回滚时退回 bet。
			rollbackAmount = order.BetAmount

			balance, err = s.addWithTx(
				ctx,
				tx,
				&provider.ChangeBalanceReq{
					UID:        order.UID,
					AgentID:    order.AgentID,
					GameID:     order.GameID,
					Amount:     rollbackAmount,
					Type:    	model.WalletLogTypeRollback,
					RefOrderNo: order.OrderNo,
				},
			)
			if err != nil {
				return err
			}

		case model.OrderStatusSettled:
			// 已结算：
			//
			// 当前余额变化为：
			//     -bet + win
			//
			// 要恢复到 Spin 前余额：
			//     rollback = bet - win
			rollbackAmount = order.BetAmount - order.WinAmount

			if rollbackAmount > 0 {
				// 输局或中奖金额小于下注金额，需要加钱
				balance, err = s.addWithTx(
					ctx,
					tx,
					&provider.ChangeBalanceReq{
						UID:        order.UID,
						AgentID:    order.AgentID,
						GameID:     order.GameID,
						Amount:     rollbackAmount,
						Type:    	model.WalletLogTypeRollback,
						RefOrderNo: order.OrderNo,
					},
				)
				if err != nil {
					return err
				}
			} else if rollbackAmount < 0 {
				// 中奖金额大于下注金额，需要扣除多出的奖金
				balance, err = s.subWithTx(
					ctx,
					tx,
					&provider.ChangeBalanceReq{
						UID:        order.UID,
						AgentID:    order.AgentID,
						GameID:     order.GameID,
						Amount:     -rollbackAmount,
						Type:    	model.WalletLogTypeRollback,
						RefOrderNo: order.OrderNo,
					},
				)
				if err != nil {
					return err
				}
			} else {
				// bet == win，余额不需要变化
				walletRepo := s.walletRepo.WithTx(tx)

				wallet, err := walletRepo.GetByUIDForUpdate(
					ctx,
					order.UID,
				)
				if err != nil {
					return err
				}

				balance = wallet.Balance
			}
		}

		// 更新订单状态
		if err := orderRepo.Rollback(
			ctx,
			order.OrderNo,
			req.Reason,
		); err != nil {
			return err
		}

		// 写入回滚记录
		if err := rollbackLogRepo.Create(
			ctx,
			&model.RollbackLog{
				RollbackType: req.RollbackType,
				RollbackNo:   pkg.GenRollbackNo(),
				OrderNo:     order.OrderNo,
				RoundID:     order.RoundID,
				RequestID:   req.RequestID,
				AgentID:     order.AgentID,
				UID:         order.UID,
				GameID:      order.GameID,

				// 保留正负值：
				// 正数代表加回钱包；
				// 负数代表从钱包扣除。
				Amount: rollbackAmount,

				Reason: req.Reason,
				Status: model.RollbackStatusSuccess,
			},
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	// 必须等事务成功后再删除钱包缓存
	if uid > 0 {
		s.deleteCache(ctx, uid)
	}

	return balance, nil
}

// Balance 查询余额
// func (s *WalletService) Balance(
// 	ctx context.Context,
// 	uid uint64,
// ) (*provider.BalanceResp, error) {
// 	wallet, err := s.getWallet(ctx, uid)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, pkg.ErrWalletNotFound
// 		}

// 		return nil, err
// 	}

// 	return &provider.BalanceResp{
// 		Balance:  pkg.ToAmount(wallet.Balance),
// 		//Currency: wallet.Currency,
// 	}, nil
// }

// Info 查询钱包完整信息
func (s *WalletService) Info(
	ctx context.Context,
	uid uint64,
) (*model.Wallet, error) {
	wallet, err := s.getWallet(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrWalletNotFound
		}

		return nil, fmt.Errorf("get wallet: %w", err)
	}

	return wallet, nil
}

// 读取缓存
func (s *WalletService) getWalletCache(
	ctx context.Context,
	uid uint64,
) (*model.Wallet, error) {

	key := walletCacheKey(uid)

	val, err := s.redisClient.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var wallet model.Wallet
	if err := json.Unmarshal([]byte(val), &wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}

// 写入缓存
func (s *WalletService) saveWalletCache(
	ctx context.Context,
	wallet *model.Wallet,
) error {

	key := walletCacheKey(wallet.UID)

	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	return s.redisClient.Set(
		key,
		data,
		walletCacheTTL,
	).Err()
}

// 删除缓存
func (s *WalletService) deleteWalletCache(
	ctx context.Context,
	uid uint64,
) error {

	return s.redisClient.
		Del(walletCacheKey(uid)).
		Err()
}

// replaySpin 幂等回放
//
// 当相同 request_id 再次请求时，不重新生成 order_no，
// 不重新调用 Operator Bet，也不重新执行 Skynet Spin。
//
// 根据原订单状态返回原结果，或者提示订单仍在处理中。
func (s *WalletService) replaySpin(
	order *model.GameOrder,
) (*provider.SpinResp, error) {

	if order == nil {
		return nil, pkg.ErrOrderNotFound
	}

	log.Printf("========== Spin Idempotent Replay ==========")
	log.Printf("RequestID    : %s", order.RequestID)
	log.Printf("OrderNo      : %s", order.OrderNo)
	log.Printf("RoundID      : %s", order.RoundID)
	log.Printf("UID          : %d", order.UID)
	log.Printf("AgentID      : %d", order.AgentID)
	log.Printf("GameID       : %d", order.GameID)
	log.Printf("Status       : %d", order.Status)
	log.Printf("BetAmount    : %d", order.BetAmount)
	log.Printf("WinAmount    : %d", order.WinAmount)
	log.Printf("BalanceAfter : %d", order.BalanceAfter)
	log.Printf("Currency     : %s", order.Currency)

	switch order.Status {

	// Provider 已创建订单，但 Operator 下注结果尚未确认。
	case model.OrderStatusPending:
		log.Printf(
			"Spin replay: order is pending, request_id=%s order_no=%s",
			order.RequestID,
			order.OrderNo,
		)

		return nil, pkg.NewError(
			pkg.ORDER_PROCESSING,
			"order is processing",
		)

	// Operator 已扣款，等待 Skynet 产生游戏结果。
	//
	// 不能重新调用 Bet，也不能重新生成 order_no。
	// 应由 Worker 或恢复流程查询 Skynet 游戏结果。
	case model.OrderStatusBetSuccess:
		log.Printf(
			"Spin replay: bet succeeded, waiting game result, request_id=%s order_no=%s",
			order.RequestID,
			order.OrderNo,
		)

		return nil, pkg.NewError(
			pkg.ORDER_PROCESSING,
			"bet succeeded, waiting for game result",
		)

	// Skynet 已经产生游戏结果，Operator 派奖尚未完成。
	//
	// 不能重新执行 Spin。
	// Worker 应使用原 order_no 和 win_amount 重试 Settle。
	case model.OrderStatusWaitSettle:
		log.Printf(
			"Spin replay: waiting settle, request_id=%s order_no=%s round_id=%s win_amount=%d",
			order.RequestID,
			order.OrderNo,
			order.RoundID,
			order.WinAmount,
		)

		return nil, pkg.NewError(
			pkg.ORDER_PROCESSING,
			"order is waiting for settlement",
		)

	// 整笔注单已完成。
	//
	// 直接返回第一次请求保存的结果，不能再次调用 Operator。
	case model.OrderStatusSettled:
		log.Printf(
			"Spin replay success: request_id=%s order_no=%s balance=%d currency=%s",
			order.RequestID,
			order.OrderNo,
			order.BalanceAfter,
			order.Currency,
		)

		return &provider.SpinResp{
			Balance:  pkg.ToAmount(order.BalanceAfter),
			Currency: order.Currency,
		}, nil

	// Skynet 明确失败，但 Operator 回滚尚未成功。
	//
	// Worker 应使用原 order_no 重试 Rollback。
	case model.OrderStatusWaitRollback:
		log.Printf(
			"Spin replay: waiting rollback, request_id=%s order_no=%s reason=%s",
			order.RequestID,
			order.OrderNo,
			order.RollbackReason,
		)

		return nil, pkg.NewError(
			pkg.ORDER_PROCESSING,
			"order is waiting for rollback",
		)

	// Operator 已完成回滚。
	//
	// 不能重新下注，因为同一个 request_id 对应的原请求已经终止。
	case model.OrderStatusRolledBack:
		log.Printf(
			"Spin replay: order already rolled back, request_id=%s order_no=%s reason=%s",
			order.RequestID,
			order.OrderNo,
			order.RollbackReason,
		)

		return nil, pkg.NewError(
			pkg.ORDER_ROLLED_BACK,
			"order has been rolled back",
		)

	// 明确失败且不再补偿。
	case model.OrderStatusFailed:
		log.Printf(
			"Spin replay: order failed, request_id=%s order_no=%s",
			order.RequestID,
			order.OrderNo,
		)

		return nil, pkg.NewError(
			pkg.ORDER_FAILED,
			"order failed",
		)

	default:
		log.Printf(
			"Spin replay invalid status: request_id=%s order_no=%s status=%d",
			order.RequestID,
			order.OrderNo,
			order.Status,
		)

		return nil, pkg.NewError(
			pkg.ORDER_STATUS_ERROR,
			"invalid order status",
		)
	}
}