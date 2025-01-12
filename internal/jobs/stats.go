package jobs

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// StatsJob handles the periodic display of application stats
func StatsJob(ctx context.Context, startTime time.Time, count *uint64, runtimeGauge prometheus.Gauge, logger *logrus.Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	statsTicker := time.NewTicker(5 * time.Second)
	defer statsTicker.Stop()

	for {
		select {
		case <-statsTicker.C:
			elapsed := time.Since(startTime).Seconds()
			currentCount := atomic.LoadUint64(count)
			logger.Infof("[Stats] Elapsed time: %.2fs | Messages printed: %d", elapsed, currentCount)
			runtimeGauge.Set(elapsed)
		case <-ctx.Done():
			logger.Info("Shutting down Stats job...")
			return
		}
	}
}
