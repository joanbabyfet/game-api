package cron

import (
	"context"
	"game-api/internal/config"
	"game-api/internal/repository"
	"game-api/internal/service"
	"log"
	"time"
)

type GameOrderRecoveryJob struct {
	repo    *repository.GameOrderRepository
	service *service.WalletService
}

func NewGameOrderRecoveryJob(repo *repository.GameOrderRepository, service *service.WalletService) *GameOrderRecoveryJob {
	return &GameOrderRecoveryJob{repo: repo, service: service}
}

//现扫描、抢占、执行、失败退避和释放
func (j *GameOrderRecoveryJob) Run() {
	cfg := config.Cfg.Cron
	if cfg.RecoverBatchSize <= 0 || cfg.RecoverMaxRetry <= 0 || cfg.RecoverLockTTL <= 0 {
		log.Printf("[Recovery] invalid worker configuration")
		return
	}

	now := time.Now().Unix()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	orders, err := j.repo.ListRecoverable(
		ctx, now, now-int64(cfg.RecoverMinAge), uint32(cfg.RecoverMaxRetry), cfg.RecoverBatchSize,
	)
	if err != nil {
		log.Printf("[Recovery] list orders failed: %v", err)
		return
	}

	for i := range orders {
		order := &orders[i]
		claimed, err := j.repo.ClaimRecovery(ctx, order.ID, now, now+int64(cfg.RecoverLockTTL))
		if err != nil || !claimed {
			continue
		}

		if err := j.service.RecoverOrder(ctx, order); err != nil {
			next := time.Now().Unix() + retryDelay(order.RetryCount+1)
			if markErr := j.repo.MarkRetryFailed(ctx, order.ID, err.Error(), next); markErr != nil {
				log.Printf("[Recovery] mark retry failed order_no=%s err=%v", order.OrderNo, markErr)
			}
			log.Printf("[Recovery] failed order_no=%s status=%d wallet_mode=%d err=%v", order.OrderNo, order.Status, order.WalletMode, err)
			continue
		}

		if err := j.repo.ReleaseRecovery(ctx, order.ID); err != nil {
			log.Printf("[Recovery] release order failed order_no=%s err=%v", order.OrderNo, err)
			continue
		}
		log.Printf("[Recovery] success order_no=%s", order.OrderNo)
	}
}

//重试间隔 (最多：20次)
func retryDelay(retryCount uint32) int64 {
	switch retryCount {
	case 1: //第1次：10秒
		return 10 
	case 2: //第2次：30秒
		return 30
	case 3: //第3次：60秒
		return 60
	case 4: //第4次：5分钟
		return 300
	default: //之后：10分钟
		return 600
	}
}
