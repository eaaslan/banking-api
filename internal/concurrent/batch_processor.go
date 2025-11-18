package concurrent

import (
	"context"
	"go-banking-api/internal/model"
	"sync"
	"sync/atomic"
	"time"
)

type BatchResult struct {
	TotalProcessed int
	Successful     int
	Failed         int
	FailedItems    []*model.Transaction
	Duration       time.Duration
}

type BatchProcessor struct {
	concurrency    int
	processor      TransactionProcessor
	processedCount atomic.Int64
	successCount   atomic.Int64
	failedCount    atomic.Int64
}

func NewBatchProcessor(concurrency int, processor TransactionProcessor) *BatchProcessor {
	if concurrency <= 0 {
		concurrency = 1
	}

	return &BatchProcessor{
		concurrency: concurrency,
		processor:   processor,
	}
}

func (bp *BatchProcessor) ProcessBatch(ctx context.Context, transactions []*model.Transaction) *BatchResult {
	startTime := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex
	failedItems := make([]*model.Transaction, 0)

	bp.processedCount.Store(0)
	bp.successCount.Store(0)
	bp.failedCount.Store(0)

	jobs := make(chan *model.Transaction, len(transactions))

	for i := 0; i < bp.concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for tx := range jobs {
				bp.processedCount.Add(1)

				err := bp.processor.Process(ctx, tx)
				if err != nil {
					bp.failedCount.Add(1)

					mu.Lock()
					failedItems = append(failedItems, tx)
					mu.Unlock()
				} else {
					bp.successCount.Add(1)
				}
			}
		}(i)
	}

	for _, tx := range transactions {
		select {
		case jobs <- tx:

		case <-ctx.Done():

			close(jobs)
			wg.Wait()

			return &BatchResult{
				TotalProcessed: int(bp.processedCount.Load()),
				Successful:     int(bp.successCount.Load()),
				Failed:         int(bp.failedCount.Load()),
				FailedItems:    failedItems,
				Duration:       time.Since(startTime),
			}
		}
	}

	close(jobs)
	wg.Wait()

	return &BatchResult{
		TotalProcessed: int(bp.processedCount.Load()),
		Successful:     int(bp.successCount.Load()),
		Failed:         int(bp.failedCount.Load()),
		FailedItems:    failedItems,
		Duration:       time.Since(startTime),
	}
}

// GetStats returns the current statistics of the batch processor
func (bp *BatchProcessor) GetStats() (processed, success, failed int64) {
	return bp.processedCount.Load(), bp.successCount.Load(), bp.failedCount.Load()
}
