package services

import (
	"admin-dashboard/cmd/internal/config"
	"context"
	"testing"
	"time"
)

func TestWorkerGracefulShutdown(t *testing.T) {
	svc := &Service{
		Config: config.Config{
			QueueWorkerInterval:  1,
			QueueWorkerBatchSize: 10,
			QueueWorkerMaxTries:  3,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		svc.StartWorker(ctx, 10*time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
		// Succeeded: Worker stopped cleanly upon context expiry
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("worker failed to stop gracefully within timeout")
	}
}

func TestWorkerImmediateCancellation(t *testing.T) {
	svc := &Service{
		Config: config.Config{
			QueueWorkerInterval:  0,
			QueueWorkerBatchSize: 0,
			QueueWorkerMaxTries:  0,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	done := make(chan struct{})
	go func() {
		svc.StartWorker(ctx, 0)
		close(done)
	}()

	select {
	case <-done:
		// Succeeded: Worker exited immediately
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("worker hung on already-cancelled context")
	}
}
