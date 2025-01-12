package jobs

import (
	"context"
	"sync"

	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// PrintJob handles the periodic printing of messages
func PrintJob(ctx context.Context, interval time.Duration, message string, count *uint64, msgCounter prometheus.Counter, logger *logrus.Logger, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			atomic.AddUint64(count, 1)
			logger.Infof("[Print] %s (Message #%d)", message, *count)
			msgCounter.Inc()
		case <-ctx.Done():
			logger.Info("Shutting down Print job...")
			return
		}
	}
}
