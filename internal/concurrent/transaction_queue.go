package concurrent

import (
	"context"
	"go-banking-api/internal/model"
	"sync"
	"sync/atomic"
	"time"
)

type TransactionQueue struct {
	queue         chan *model.Transaction
	capacity      int
	enqueueCount  atomic.Int64
	dequeueCount  atomic.Int64
	rejectedCount atomic.Int64
	timeoutCount  atomic.Int64
	mu            sync.RWMutex
	isOpen        bool
}

func NewTransactionQueue(capacity int) *TransactionQueue {
	return &TransactionQueue{
		queue:    make(chan *model.Transaction, capacity),
		capacity: capacity,
		isOpen:   true,
	}
}

func (q *TransactionQueue) Enqueue(transaction *model.Transaction) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if !q.isOpen {
		q.rejectedCount.Add(1)
		return false
	}

	select {
	case q.queue <- transaction:
		q.enqueueCount.Add(1)
		return true
	default:
		// Queue is full
		q.rejectedCount.Add(1)
		return false
	}
}

func (q *TransactionQueue) EnqueueWithTimeout(transaction *model.Transaction, timeout time.Duration) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if !q.isOpen {
		q.rejectedCount.Add(1)
		return false
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case q.queue <- transaction:
		q.enqueueCount.Add(1)
		return true
	case <-timer.C:
		q.timeoutCount.Add(1)
		return false
	}
}

func (q *TransactionQueue) Dequeue(ctx context.Context) *model.Transaction {
	select {
	case transaction, ok := <-q.queue:
		if ok {
			q.dequeueCount.Add(1)
			return transaction
		}
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (q *TransactionQueue) DequeueWithTimeout(timeout time.Duration) *model.Transaction {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case transaction, ok := <-q.queue:
		if ok {
			q.dequeueCount.Add(1)
			return transaction
		}
		return nil
	case <-timer.C:
		return nil
	}
}

func (q *TransactionQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isOpen {
		q.isOpen = false
		close(q.queue)
	}
}

func (q *TransactionQueue) Size() int {
	return len(q.queue)
}

func (q *TransactionQueue) Capacity() int {
	return q.capacity
}

func (q *TransactionQueue) IsOpen() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.isOpen
}

func (q *TransactionQueue) GetStats() (enqueued, dequeued, rejected, timedOut int64) {
	return q.enqueueCount.Load(), q.dequeueCount.Load(), q.rejectedCount.Load(), q.timeoutCount.Load()
}
