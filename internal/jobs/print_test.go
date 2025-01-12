package jobs_test

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/phunguyen19/golang-project-involvements/internal/jobs"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/sirupsen/logrus"
)

func TestPrintJob(t *testing.T) {
	var count uint64
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_messages_printed_total",
		Help: "Test counter.",
	})

	log := logrus.New()
	log.SetOutput(io.Discard) // Disable logging during tests

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go jobs.PrintJob(ctx, 1*time.Second, "Test Message", &count, counter, log, &wg)

	// Let it run for 350ms, expecting ~3 messages
	time.Sleep(3500 * time.Millisecond)
	cancel()
	wg.Wait()

	assert.Equal(t, count, uint64(3))

	// Check Prometheus counter
	actualCount := testutil.ToFloat64(counter)
	assert.Equal(t, float64(3), actualCount)
}
