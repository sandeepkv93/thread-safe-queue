// Package threadsafequeue provides a thread-safe queue implementation.
package threadsafequeue

import (
	"sync"
)

// ThreadSafeQueue represents a FIFO (first-in-first-out) data structure that
// supports safe concurrent access. It uses a slice to store the items
// and a condition variable to synchronize access.
type ThreadSafeQueue struct {
	queue []interface{} // Internal slice to hold the queue items.
	mu    sync.Mutex    // Mutex to protect concurrent access to the queue slice.
	cond  *sync.Cond    // Condition variable to coordinate enqueue and dequeue operations.
}

// NewThreadSafeQueue initializes and returns a new instance of ThreadSafeQueue.
// It is safe to be used concurrently.
func NewThreadSafeQueue() *ThreadSafeQueue {
	q := &ThreadSafeQueue{}
	q.cond = sync.NewCond(&q.mu) // Create a condition variable with the queue's mutex.
	return q
}

// Enqueue adds an item to the end of the queue. The provided item can be of any type.
// If there are any waiting Dequeue calls, it signals one of them that an item is available.
// This method is safe for concurrent use.
func (q *ThreadSafeQueue) Enqueue(item interface{}) {
	q.mu.Lock() // Lock the mutex to protect concurrent access.
	q.queue = append(q.queue, item)
	q.cond.Signal() // Signal any waiting Dequeue operations that a new item is available.
	q.mu.Unlock()
}

// Dequeue removes and returns the item from the front of the queue.
// If the queue is empty, this call will block until an item is enqueued.
// The return value is the dequeued item and a boolean indicating success.
// If the queue is empty, the boolean value will be false.
// This method is safe for concurrent use.
func (q *ThreadSafeQueue) Dequeue() (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.queue) == 0 {
		q.cond.Wait() // Wait until an item is available.
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item, true
}

// IsEmpty returns true if the queue has no items, and false otherwise.
// This method is safe for concurrent use.
func (q *ThreadSafeQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0
}

// Size returns the number of items currently in the queue.
// This method is safe for concurrent use.
func (q *ThreadSafeQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}
