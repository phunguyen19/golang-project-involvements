package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/pflag"
	"github.com/phunguyen19/golang-project-involvements/internal/config"
	"github.com/phunguyen19/golang-project-involvements/internal/health"
	"github.com/phunguyen19/golang-project-involvements/internal/jobs"
	"github.com/phunguyen19/golang-project-involvements/internal/logger"
	"github.com/phunguyen19/golang-project-involvements/internal/metrics"
)

func main() {
	// Setup logger
	log := logger.NewLogger()

	// Create a new FlagSet
	flags := pflag.NewFlagSet("app", pflag.ExitOnError)

	// Pass os.Args[1:] to capture all command-line arguments except the program name
	args := os.Args[1:]

	// Load configuration
	cfg, err := config.LoadConfig(log, flags, args)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup WaitGroup
	var wg sync.WaitGroup

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start metrics server
	wg.Add(1)
	go metrics.StartMetricsServer(ctx, cfg.MetricsPort, log, &wg)

	// Start health server
	wg.Add(1)
	go health.StartHealthServer(ctx, cfg.HealthPort, log, &wg)

	log.Info("App started. Press Ctrl+C to exit.")

	// Tracking variables
	var messageCount uint64
	startTime := time.Now()

	// metric message count
	messageCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "messages_printed_total",
		Help: "Total number of messages printed.",
	})
	prometheus.MustRegister(messageCounter)

	// metrics runtime
	elapsedGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_runtime_seconds",
		Help: "Elapsed time since application started in seconds.",
	})
	prometheus.MustRegister(elapsedGauge)

	// Start print job
	wg.Add(1)
	go jobs.PrintJob(ctx, cfg.TickInterval, cfg.Message, &messageCount, messageCounter, log, &wg)

	// Start stats job
	wg.Add(1)
	go jobs.StatsJob(ctx, startTime, &messageCount, elapsedGauge, log, &wg)

	// Wait for termination signal
	<-sigChan
	log.Info("Termination signal received. Shutting down...")

	// Cancel the context to signal goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()

	// Final shutdown message
	elapsed := time.Since(startTime).Seconds()
	totalMessages := messageCount
	log.Infof("Total messages printed: %d | Total runtime: %.2fs", totalMessages, elapsed)
	log.Info("Goodbye!")
}
