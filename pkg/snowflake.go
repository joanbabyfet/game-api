package pkg

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	// 自定义开始时间：2024-01-01 00:00:00 UTC (ms)
	epoch int64 = 1704067200000

	workerBits   uint8 = 10
	sequenceBits uint8 = 12

	maxWorkerID int64 = -1 ^ (-1 << workerBits)
	maxSequence int64 = -1 ^ (-1 << sequenceBits)

	workerShift    = sequenceBits
	timestampShift = sequenceBits + workerBits
)

var snowflake *Snowflake

type Snowflake struct {
	mu        sync.Mutex
	workerID  int64
	lastTs    int64
	sequence  int64
}

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("invalid worker id")
	}

	return &Snowflake{
		workerID: workerID,
	}, nil
}

// 初始化（程序启动时调用一次）
func InitSnowflake(workerID int64) error {
	sf, err := NewSnowflake(workerID)
	if err != nil {
		return err
	}

	snowflake = sf
	return nil
}

func currentMs() int64 {
	return time.Now().UnixMilli()
}

func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	ts := currentMs()

	if ts < s.lastTs {
		panic("clock moved backwards")
	}

	if ts == s.lastTs {
		s.sequence++

		if s.sequence > maxSequence {
			for ts <= s.lastTs {
				ts = currentMs()
			}
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}

	s.lastTs = ts

	timestamp := ts - epoch

	id :=
		(timestamp << timestampShift) |
			(s.workerID << workerShift) |
			s.sequence

	return id
}

func (s *Snowflake) NextIDString() string {
	return strconv.FormatInt(s.NextID(), 10)
}