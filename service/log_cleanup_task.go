package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"

	"github.com/bytedance/gopkg/util/gopool"
)

const (
	logCleanupTickInterval = 1 * time.Hour
	logCleanupBatchSize    = 1000
)

var (
	logCleanupOnce sync.Once
	logCleanupRunning atomic.Bool
)

func StartLogCleanupTask() {
	logCleanupOnce.Do(func() {
		if !common.IsMasterNode {
			return
		}
		gopool.Go(func() {
			logger.LogInfo(context.Background(), fmt.Sprintf("log cleanup task started: tick=%s", logCleanupTickInterval))
			ticker := time.NewTicker(logCleanupTickInterval)
			defer ticker.Stop()

			runLogCleanupOnce()
			for range ticker.C {
				runLogCleanupOnce()
			}
		})
	})
}

func runLogCleanupOnce() {
	if !logCleanupRunning.CompareAndSwap(false, true) {
		return
	}
	defer logCleanupRunning.Store(false)

	ctx := context.Background()

	// 计算今天零点的时间戳
	now := time.Now()
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	targetTimestamp := todayMidnight.Unix()

	deleted, err := model.DeleteOldLog(ctx, targetTimestamp, logCleanupBatchSize)
	if err != nil {
		logger.LogWarn(ctx, fmt.Sprintf("log cleanup task failed: %v", err))
		return
	}

	if common.DebugEnabled && deleted > 0 {
		logger.LogDebug(ctx, fmt.Sprintf("log cleanup: deleted=%d, before=%s", deleted, todayMidnight.Format("2006-01-02")))
	}
}