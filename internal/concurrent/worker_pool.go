package concurrent

import (
	"context"
	"go-banking-api/internal/model"
	"log"
	"sync"
	"sync/atomic"
)

type TransactionProcessor interface {
	Process(ctx context.Context, transaction *model.Transaction) error
}

type WorkerPool struct {
	numWorkers     int
	jobQueue       chan *model.Transaction
	processor      TransactionProcessor
	wg             sync.WaitGroup
	processedCount atomic.Int64
	successCount   atomic.Int64
	failedCount    atomic.Int64
	isRunning      bool
	mu             sync.Mutex
}

func NewWorkerPool(numWorkers int, queueSize int, processor TransactionProcessor) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan *model.Transaction, queueSize),
		processor:  processor,
		isRunning:  false,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.isRunning {
		return
	}

	wp.isRunning = true
	wp.wg.Add(wp.numWorkers)

	for i := 0; i < wp.numWorkers; i++ {
		go wp.worker(ctx, i)
	}

	log.Printf("Started worker pool with %d workers", wp.numWorkers)
}

func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.isRunning {
		return
	}

	close(wp.jobQueue)
	wp.wg.Wait()
	wp.isRunning = false

	log.Printf("Stopped worker pool. Processed: %d, Success: %d, Failed: %d",
		wp.processedCount.Load(), wp.successCount.Load(), wp.failedCount.Load())
}

func (wp *WorkerPool) EnqueueTransaction(transaction *model.Transaction) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.isRunning {
		return false
	}

	select {
	case wp.jobQueue <- transaction:
		return true
	default:
		// Queue is full
		return false
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	log.Printf("Worker %d started", id)

	for transaction := range wp.jobQueue {
		wp.processedCount.Add(1)

		err := wp.processor.Process(ctx, transaction)
		if err != nil {
			wp.failedCount.Add(1)
			log.Printf("Worker %d failed to process transaction %d: %v", id, transaction.ID, err)
		} else {
			wp.successCount.Add(1)
			log.Printf("Worker %d successfully processed transaction %d", id, transaction.ID)
		}
	}

	log.Printf("Worker %d stopped", id)
}

func (wp *WorkerPool) GetStats() (processed, success, failed int64) {
	return wp.processedCount.Load(), wp.successCount.Load(), wp.failedCount.Load()
}

func (wp *WorkerPool) IsRunning() bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.isRunning
}
