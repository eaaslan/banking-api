package worker

import (
	"context"
	"log/slog"
	"sync/atomic"

	"backend/internal/models"
)

type ProcessorFunc func(ctx context.Context, tx *models.Transaction) error

type Pool struct {
	queue          chan *models.Transaction
	workerCount    int
	processor      ProcessorFunc
	processedCount int64
	errorCount     int64
}

func NewPool(workerCount int, queueSize int, processor ProcessorFunc) *Pool {
	return &Pool{
		queue:       make(chan *models.Transaction, queueSize),
		workerCount: workerCount,
		processor:   processor,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		go p.worker(ctx, i)
	}
	slog.Info("Worker pool started", "workers", p.workerCount)
}

func (p *Pool) worker(ctx context.Context, id int) {
	for {
		select {
		case tx := <-p.queue:
			if err := p.processor(ctx, tx); err != nil {
				atomic.AddInt64(&p.errorCount, 1)
				slog.Error("Worker failed to process transaction",
					"worker_id", id,
					"tx_id", tx.ID,
					"type", tx.Type,
					"amount", tx.Amount,
					"from_user", tx.FromUserID,
					"to_user", tx.ToUserID,
					"error", err,
				)
			} else {
				atomic.AddInt64(&p.processedCount, 1)
				slog.Info("Worker processed transaction", "worker_id", id, "tx_id", tx.ID)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *Pool) Submit(tx *models.Transaction) {
	p.queue <- tx
}
func (p *Pool) Stats() (processed int64, errors int64) {
	return atomic.LoadInt64(&p.processedCount), atomic.LoadInt64(&p.errorCount)
}
