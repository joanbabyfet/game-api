package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game-api/internal/adapter"
	"game-api/internal/client/operator"
	"game-api/internal/client/skynet"
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
		return nil, pkg.ErrUnauthorized
	}

	// 2. 查询游戏
	game, err := s.gameRepo.GetByCode(req.GameCode)
	if err != nil {
		return nil, pkg.ErrGameNotFound
	}

	// 3. 查询代理信息
	agent, err := s.agentRepo.GetByID(claims.AgentID)
	if err != nil {
		return nil, err
	}

	//是否免费旋转
	isFreeSpin := req.FreeSpinID != ""
	spinType := model.SpinTypeNormal
	if isFreeSpin {
		spinType = model.SpinTypeFreeSpin
	}

	// 普通 Spin 必须有下注金额
	if !isFreeSpin && req.BetAmount <= 0 {
			return nil, pkg.NewError(
					pkg.INVALID_PARAM,
					"bet_amount must be greater than zero",
			)
	}
	// Free Spin 不接受正数或负数下注额
	if isFreeSpin && req.BetAmount != 0 {
			return nil, pkg.NewError(
					pkg.INVALID_PARAM,
					"bet_amount must be zero or omitted for free spin",
			)
	}

	// 4. request_id 幂等检查
	order, err := s.orderRepo.GetByRequestID(ctx, req.RequestID)
	if err == nil {
		// 已存在注单，直接回放，不重新下注
		return s.replaySpin(order)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 查询数据库异常
		return nil, err
	}

	//这些由入口 Provider API 生成
	requestID := req.RequestID //幂等(由 Provider API 传入)
	orderNo := pkg.GenOrderNo() //注單号

	// 普通 Spin 的下注金额来自请求
	// Free Spin 的基准下注金额暂时为 0，等 Skynet 校验 free_spin 后返回真实基准下注额
	betAmount := int64(0)
	orderStatus := model.OrderStatusPending

	if !isFreeSpin {
		betAmount = pkg.ToMoney(req.BetAmount)
	} else {
		// Free Spin 没有 Operator Bet，直接进入可执行游戏状态
		orderStatus = model.OrderStatusBetSuccess
	}

	//5. 第一次请求，开始创建注单
	order = &model.GameOrder{
		RequestID:  requestID,
		OrderNo:    orderNo,
		UID:        claims.UID,
		AgentID:    claims.AgentID,
		GameID:     game.ID,
		BetAmount:  betAmount,
		WinAmount:  0,
		Profit:     0,
		Currency:   claims.Currency,
		SpinType:   spinType,
		FreeSpinID:   req.FreeSpinID,
		FreeSpinIndex: 0,
		Status:     orderStatus, //Provider 已创建 game_order，Operator 下注结果尚未确认
		CreateTime: time.Now().Unix(),
	}
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	log.Printf(
		"[Spin] start request_id=%s order_no=%s uid=%d game=%s spin_type=%d",
		requestID,
		orderNo,
		claims.UID,
		game.GameCode,
		spinType,
	)

	// 6. 普通 Spin 才调用 Operator 扣款
	//
	// 单一钱包模式下：
	// Provider 不直接修改本地钱包余额，
	// 而是调用 Operator 的下注接口扣除玩家余额。
	if !isFreeSpin {
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
				"[Bet] failed request_id=%s order_no=%s err=%v",
				requestID,
				orderNo,
				err,
			)

			// 明确业务失败，例如余额不足
			if errors.Is(err, pkg.ErrInsufficientBalance) {
				if updateErr := s.orderRepo.UpdateStatus(
					ctx,
					order.ID,
					model.OrderStatusPending,
					model.OrderStatusFailed,
				); updateErr != nil {
					log.Printf(
						"[Order] update failed status error request_id=%s order_no=%s err=%v",
						requestID,
						orderNo,
						updateErr,
					)

					return nil, pkg.ErrOrderUpdateFailed
				}

				order.Status = model.OrderStatusFailed
			}

			// 网络超时、连接中断等结果未知，先保持 Pending
			return nil, err
		}

		if betResp == nil {
			return nil, pkg.NewError(
				pkg.OPERATOR_BET_FAILED,
				"operator bet returned nil response",
			)
		}

		// Operator Bet 返回的是扣款后的余额。
		// 因此：扣款前余额 = 扣款后余额 + 下注金额。
		balanceBefore := pkg.ToMoney(betResp.Balance) + order.BetAmount

		// 更新注单为 扣款成功
		if err := s.orderRepo.UpdateBetSuccess(
			ctx,
			order.ID,
			balanceBefore,
		); err != nil {
			log.Printf(
				"[Order] update bet success failed request_id=%s order_no=%s err=%v",
				requestID,
				orderNo,
				err,
			)

			return nil, pkg.ErrOrderUpdateFailed
		}

		order.BalanceBefore = balanceBefore
		order.Status = model.OrderStatusBetSuccess

		log.Printf(
			"[Bet] success order_no=%s balance=%v",
			orderNo,
			betResp.Balance,
		)
	}

	pbBetAmount := order.BetAmount
	if isFreeSpin {
		pbBetAmount = 0
	}

	// 7. 调用 Skynet 执行游戏逻辑
	pbReq := &slotpb.SpinReq{
		RequestId: requestID,
		OrderNo:   orderNo,
		Uid:       claims.UID,
		AgentId:   claims.AgentID,
		GameId:    game.ID,
		BetAmount: pbBetAmount,
		Currency:  claims.Currency,
		SpinType:  uint32(spinType),
		FreeSpinId: req.FreeSpinID,
		DebugFail: req.DebugFail, //测试取消下注用
	}
	var (
		spinResp *slotpb.SpinResp
		spinErr  error
	)
	if isFreeSpin {
		spinResp, spinErr = s.slotAdapter.FreeSpin(ctx, pbReq)
	} else {
		spinResp, spinErr = s.slotAdapter.Spin(ctx, pbReq)
	}
	if spinErr != nil {
		log.Printf(
			"[Skynet] failed request_id=%s order_no=%s err=%v",
			requestID,
			orderNo,
			spinErr,
		)

		// Free Spin 没有扣钱，不执行 Operator Rollback。
		if isFreeSpin {
			var rpcErr *skynet.Error

			// 非业务错误（例如网络超时），保持 BetSuccess，方便后续幂等恢复。
			if !errors.As(spinErr, &rpcErr) {
				return nil, spinErr
			}

			switch rpcErr.Code {
			case pkg.FREE_SPIN_NOT_FOUND,
				pkg.FREE_SPIN_FINISHED,
				pkg.FREE_SPIN_OWNER_ERROR,
				pkg.FREE_SPIN_GAME_ERROR:

				if updateErr := s.orderRepo.UpdateStatus(
					ctx,
					order.ID,
					model.OrderStatusBetSuccess,
					model.OrderStatusFailed,
				); updateErr != nil {

					log.Printf(
						"[Order] update free spin failed request_id=%s order_no=%s order_id=%d err=%v",
						requestID,
						orderNo,
						order.ID,
						updateErr,
					)

					return nil, pkg.ErrOrderUpdateFailed
				}

				order.Status = model.OrderStatusFailed
			}

			return nil, spinErr
		}

		// 普通 Spin：这里仅适用于 Skynet 明确业务失败。
		// TCP 超时等结果未知错误，不应立即回滚。
		if err := s.orderRepo.UpdateStatus(
			ctx,
			order.ID,
			model.OrderStatusBetSuccess,
			model.OrderStatusWaitRollback,
		); err != nil {
			log.Printf(
				"[Order] update wait rollback failed request_id=%s order_no=%s err=%v",
				requestID,
				orderNo,
				err,
			)

			return nil, pkg.ErrOrderUpdateFailed
		}

		order.Status = model.OrderStatusWaitRollback

		// Skynet 失败，回滚扣款
		rollbackResp, rollbackErr := s.operatorClient.Rollback(
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
				"[Rollback] failed request_id=%s order_no=%s err=%v",
				requestID,
				orderNo,
				rollbackErr,
			)

			// 状态保持 WAIT_ROLLBACK，交给 Worker 重试
			return nil, spinErr
		}

		if rollbackResp == nil {
			// 状态保持 WAIT_ROLLBACK
			return nil, pkg.NewError(
				pkg.OPERATOR_ROLLBACK_FAILED,
				"operator rollback returned nil response",
			)
		}

		// Operator 回滚成功
		balanceAfter := pkg.ToMoney(rollbackResp.Balance)
		// 更新注单状态为 5=已回滚
		if err := s.orderRepo.UpdateRolledBack(
			ctx,
			order.ID,
			balanceAfter,
			rollbackResp.Currency,
		); err != nil {
			log.Printf(
				"[Order] update rolled back failed request_id=%s order_no=%s err=%v",
				requestID,
				orderNo,
				err,
			)

			return nil, pkg.ErrOrderUpdateFailed
		}

		order.BalanceAfter = balanceAfter
		order.Currency = rollbackResp.Currency
		order.Status = model.OrderStatusRolledBack

		return nil, spinErr
	}

	if spinResp == nil {
		return nil, pkg.NewError(
			pkg.SPIN_FAILED,
			"skynet spin returned nil response",
		)
	}

	if spinResp.WinAmount < 0 {
		log.Printf(
			"[Skynet] invalid win amount request_id=%s order_no=%s win_amount=%d",
			requestID,
			orderNo,
			spinResp.WinAmount,
		)
		return nil, pkg.ErrInvalidParam
	}

	log.Printf(
		"[Skynet] success order_no=%s round_id=%s win_amount=%d",
		orderNo,
		spinResp.RoundId,
		spinResp.WinAmount,
	)

	if isFreeSpin {
		if spinResp.FreeSpinId != req.FreeSpinID {
			return nil, pkg.NewError(
				pkg.SPIN_FAILED,
				"free_spin_id mismatch",
			)
		}

		if spinResp.BetAmount <= 0 {
			return nil, pkg.NewError(
				pkg.SPIN_FAILED,
				"invalid free spin bet amount",
			)
		}

		// Free Spin 基准下注金额由 Skynet 决定
		order.BetAmount = spinResp.BetAmount
		order.FreeSpinIndex = spinResp.FreeSpinIndex
	}

	// 普通 Spin 净输赢 = 中奖 - 实际下注
	// Free Spin 没有扣款，所以净输赢 = 中奖金额
	profit := spinResp.WinAmount
	if !isFreeSpin {
		profit = spinResp.WinAmount - order.BetAmount
	}

	// 保存 Skynet 游戏结果，并进入待派奖状态
	if err := s.orderRepo.UpdateGameResult(
		ctx,
		order.ID,
		spinResp.RoundId,
		order.BetAmount,
		spinResp.WinAmount,
		profit,
		spinType,
		spinResp.FreeSpinId,
		spinResp.FreeSpinIndex,
	); err != nil {
		log.Printf(
			"[Order] update game result failed request_id=%s order_no=%s round_id=%s err=%v",
			requestID,
			orderNo,
			spinResp.RoundId,
			err,
		)
		// 可以使用相同 request_id 重调 Skynet，由 Skynet 回放。
		return nil, pkg.ErrOrderUpdateFailed
	}

	order.RoundID = spinResp.RoundId
	order.WinAmount = spinResp.WinAmount
	order.Profit = profit
	order.Status = model.OrderStatusWaitSettle

	// 普通 Spin：派奖本局 win_amount
	// Free Spin：同样派奖本次 Free Spin win_amount
	winAmount := pkg.ToAmount(spinResp.WinAmount)

	// 8. Operator 派奖
	//
	// 单一钱包模式下，派奖也由 Operator 修改玩家余额。
	// orderNo 继续使用原下注注单号，作为该笔注单的关联标识。
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
			"[Settle] failed request_id=%s order_no=%s win_amount=%v err=%v",
			requestID,
			orderNo,
			winAmount,
			err,
		)

		// 注意：
		// Skynet 游戏结果已经产生，不能再取消下注。
		//
		// 此时应将注单标记为 WAIT_SETTLE，
		// Worker 使用同一个 orderNo 和 winAmount 重试派彩。
		//
		// 不要重新执行 Spin，也不要重新生成注单号。
		return nil, err
	}

	if settleResp == nil {
		// 同样应该标记 WAIT_SETTLE，避免丢失派彩
		return nil, pkg.NewError(
			pkg.OPERATOR_SETTLE_FAILED,
			"operator settle returned nil response",
		)
	}

	log.Printf(
		"[Settle] success order_no=%s balance=%v",
		orderNo,
		settleResp.Balance,
	)

	// 结算后余额
	balanceAfter := pkg.ToMoney(settleResp.Balance)

	// 普通 Spin：balance_before 已经在 Operator Bet 成功后写入。
	// Free Spin：没有 Operator Bet，需要根据结算后余额反推执行前余额。
	balanceBefore := order.BalanceBefore
	if isFreeSpin {
		balanceBefore = balanceAfter - spinResp.WinAmount
	}

	// 防御性检查，正常情况下不应该小于 0。
	if balanceBefore < 0 {
		log.Printf(
			"[Settle] invalid balance before request_id=%s order_no=%s balance_after=%d win_amount=%d",
			requestID,
			orderNo,
			balanceAfter,
			spinResp.WinAmount,
		)

		return nil, pkg.NewError(
			pkg.OPERATOR_SETTLE_FAILED,
			"invalid balance before",
		)
	}

	// 更新注单为 Settled
	if err := s.orderRepo.UpdateSettled(
		ctx,
		order.ID,
		balanceBefore,
		balanceAfter,
		settleResp.Currency,
	); err != nil {
		log.Printf(
			"[Order] update settled failed request_id=%s order_no=%s err=%v",
			requestID,
			orderNo,
			err,
		)

		return nil, pkg.ErrOrderUpdateFailed
	}

	order.BalanceBefore = balanceBefore
	order.BalanceAfter = balanceAfter
	order.Currency = settleResp.Currency
	order.Status = model.OrderStatusSettled

	resp := &provider.SpinResp{
		// 单一钱包模式下，最终余额和币种以 Operator 返回为准
		Balance:  settleResp.Balance,
		Currency: settleResp.Currency,
		RoundID:             spinResp.RoundId,
		WinAmount:           pkg.ToAmount(spinResp.WinAmount),
		SpinType:            uint8(spinResp.SpinType),
		FreeSpinID:          spinResp.FreeSpinId,
		FreeSpinIndex:       spinResp.FreeSpinIndex,
		FreeSpinTotalCount:  spinResp.FreeSpinTotalCount,
		FreeSpinRemainCount: spinResp.FreeSpinRemainCount,
	}

	log.Printf(
		"[Spin] success request_id=%s order_no=%s round_id=%s spin_type=%d",
		requestID,
		orderNo,
		spinResp.RoundId,
		spinType,
	)

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

// Rollback 已结算或未结算注单回滚
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

		// 查询并锁定注单，防止并发重复回滚
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

		// 只有待结算、已结算注单可以回滚
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

		// 更新注单状态
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
// 根据原注单状态返回原结果，或者提示注单仍在处理中。
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

	// Provider 已创建注单，但 Operator 下注结果尚未确认。
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