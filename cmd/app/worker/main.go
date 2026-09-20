package main

import (
	"admin-dashboard/cmd/app/services"
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/storage"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	glog "github.com/labstack/gommon/log"
)

func main() {

	fmt.Println("==================================================")
	fmt.Println("   Admin Dashboard - Continuous Queue Worker      ")
	fmt.Println("   Checking & Dispatching Pending Email & SMS     ")
	fmt.Println("==================================================")

	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := databases.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	cache, err := cacheSystem.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	echoLogger := glog.New("worker")
	echoLogger.SetLevel(glog.INFO)

	service := services.Service{
		Db:          db,
		Logger:      echoLogger,
		CacheSystem: cache,
		Storage:     storage.New(cfg),
		Config:      *cfg,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n[WORKER] Received signal %v. Initiating graceful shutdown...\n", sig)
		cancel()
	}()

	interval := time.Duration(cfg.QueueWorkerInterval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	fmt.Printf("[WORKER] Starting continuous polling loop (interval: %v, batch size: %d, max retries: %d)...\n", interval, cfg.QueueWorkerBatchSize, cfg.QueueWorkerMaxTries)
	fmt.Println("[WORKER] Press Ctrl+C to stop.")

	service.StartWorker(ctx, interval)

	fmt.Println("[WORKER] Worker process terminated gracefully.")
}
