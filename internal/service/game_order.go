package service

import (
	"context"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
)

type GameOrderService struct {
	repo *repository.GameOrderRepository
	gameRepo *repository.GameRepository
	authService *AuthService
}

func NewGameOrderService(
	repo *repository.GameOrderRepository,
	gameRepo *repository.GameRepository,
	authService *AuthService,
) *GameOrderService {
	return &GameOrderService{
		repo: repo,
		gameRepo: gameRepo,
		authService: authService,
	}
}

// Create 新增注单
func (s *GameOrderService) Create(order *model.GameOrder) error {
	return s.repo.Create(order)
}

// Update 更新注单
func (s *GameOrderService) Update(order *model.GameOrder) error {
	return s.repo.Update(order)
}

// Delete 删除注单
func (s *GameOrderService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *GameOrderService) GetByID(id uint64) (*model.GameOrder, error) {
	return s.repo.GetByID(id)
}

// GetByOrderNo 根据订单号查询
func (s *GameOrderService) GetByOrderNo(orderNo string) (*model.GameOrder, error) {
	return s.repo.GetByOrderNo(orderNo)
}

// GetByRequestID 根据RequestID查询
func (s *GameOrderService) GetByRequestID(requestID string) (*model.GameOrder, error) {
	return s.repo.GetByRequestID(requestID)
}

// GetByRoundID 根据RoundID查询
func (s *GameOrderService) GetByRoundID(roundID string) ([]model.GameOrder, error) {
	return s.repo.GetByRoundID(roundID)
}

// List 注单列表
func (s *GameOrderService) List(q repository.GameOrderQuery) ([]model.GameOrder, error) {
	return s.repo.List(q)
}

//获取注单记录
func (s *GameOrderService) GetOrderLog(ctx context.Context, req *provider.OrderLogReq) ([]provider.OrderLogResp, error) {
	
	// 驗證 Agent 簽名
	data := pkg.BuildSignData(map[string]any{
		"app_id":    	req.AppID,
		"start_time":	req.StartTime,
		"end_time":  	req.EndTime,
		"timestamp": 	req.Timestamp,
	})

	agent, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}

	// 时间转换
	startTime, err := pkg.DateTimeToUnix(req.StartTime)
	if err != nil {
		return nil, pkg.ErrInvalidParam
	}

	endTime, err := pkg.DateTimeToUnix(req.EndTime)
	if err != nil {
		return nil, pkg.ErrInvalidParam
	}

	if endTime < startTime {
		return nil, pkg.ErrInvalidParam
	}

	// 查询注单
	query := repository.GameOrderQuery{
		AgentID:   agent.ID,
		SettleStartTime: uint32(startTime),
		SettleEndTime:   uint32(endTime),
	}

	orders, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	// 查询游戏（去重）
	gameIDSet := make(map[uint32]struct{})

	for _, order := range orders {
		gameIDSet[order.GameID] = struct{}{}
	}

	gameIDs := make([]uint32, 0, len(gameIDSet))
	for id := range gameIDSet {
		gameIDs = append(gameIDs, id)
	}

	games, err := s.gameRepo.ListByIDs(gameIDs)
	if err != nil {
		return nil, err
	}

	gameMap := make(map[uint32]string, len(games))
	for _, game := range games {
		gameMap[uint32(game.ID)] = game.GameCode
	}

	// Response
	resp := make([]provider.OrderLogResp, 0, len(orders))

	for _, order := range orders {
		resp = append(resp, provider.OrderLogResp{
			OrderNo:    order.OrderNo,
			RoundID:    order.RoundID,
			UID:        order.UID,
			GameCode:   gameMap[order.GameID],
			BetAmount:  pkg.ToAmount(int64(order.BetAmount)),
			WinAmount:  pkg.ToAmount(int64(order.WinAmount)),
			Currency:   order.Currency,
			Status:     order.Status,
			SettleTime: order.SettleTime,
			CreateTime: order.CreateTime,
		})
	}

	return resp, nil
}