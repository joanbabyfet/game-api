package service

import (
	"game-api/internal/model"
	"game-api/internal/repository"
)

type RequestLogService struct {
	repo *repository.RequestLogRepository
}

func NewRequestLogService(
	repo *repository.RequestLogRepository,
) *RequestLogService {
	return &RequestLogService{
		repo: repo,
	}
}

// Create 新增请求日志
func (s *RequestLogService) Create(log *model.RequestLog) error {
	return s.repo.Create(log)
}

// Update 更新请求日志
func (s *RequestLogService) Update(log *model.RequestLog) error {
	return s.repo.Update(log)
}

// Delete 删除请求日志
func (s *RequestLogService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *RequestLogService) GetByID(id uint64) (*model.RequestLog, error) {
	return s.repo.GetByID(id)
}

// GetByRequestID 根据RequestID查询
func (s *RequestLogService) GetByRequestID(requestID string) (*model.RequestLog, error) {
	return s.repo.GetByRequestID(requestID)
}

// GetByOrderNo 根据注单号查询
func (s *RequestLogService) GetByOrderNo(orderNo string) ([]model.RequestLog, error) {
	return s.repo.GetByOrderNo(orderNo)
}

// List 请求日志列表
func (s *RequestLogService) List(q repository.RequestLogQuery) ([]model.RequestLog, error) {
	return s.repo.List(q)
}