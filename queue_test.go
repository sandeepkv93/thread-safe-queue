package threadsafequeue

import (
	"sync/atomic"
	"testing"
	"time"
)

// Test that NewThreadSafeQueue returns a non-nil queue
func TestNewThreadSafeQueue(t *testing.T) {
	q := NewThreadSafeQueue()
	if q == nil {
		t.Error("Expected new queue to be non-nil")
	}

	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
}

// Test that Enqueue and Dequeue work as expected
func TestEnqueueDequeue(t *testing.T) {
	q := NewThreadSafeQueue()

	q.Enqueue(42)

	if q.IsEmpty() {
		t.Error("Queue should not be empty after enqueue")
	}

	if q.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", q.Size())
	}

	item, ok := q.Dequeue()

	if !ok || item != 42 {
		t.Errorf("Expected to dequeue 42, got %v", item)
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty after dequeue")
	}
}

func TestDequeueFromEmptyQueue(t *testing.T) {
	q := NewThreadSafeQueue()
	go func() {
		q.Enqueue(42)
	}()
	item, ok := q.Dequeue()
	if !ok || item != 42 {
		t.Errorf("Expected to dequeue 42, got %v", item)
	}
}

func TestSize(t *testing.T) {
	q := NewThreadSafeQueue()
	for i := 0; i < 5; i++ {
		q.Enqueue(i)
	}
	if q.Size() != 5 {
		t.Errorf("Expected size to be 5, got %d", q.Size())
	}
}

func TestIsEmpty(t *testing.T) {
	q := NewThreadSafeQueue()
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}

	q.Enqueue(42)

	if q.IsEmpty() {
		t.Error("Queue should not be empty after enqueue")
	}
}

// Test concurrent enqueues and dequeues
func TestConcurrentOperations(t *testing.T) {
	q := NewThreadSafeQueue()
	const count = 1000
	done := make(chan bool)

	go func() {
		for i := 0; i < count; i++ {
			q.Enqueue(i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < count; i++ {
			_, ok := q.Dequeue()
			if !ok {
				t.Errorf("Dequeue failed on iteration %d", i)
			}
		}
		done <- true
	}()

	<-done
	<-done

	if q.Size() != 0 {
		t.Errorf("Expected queue size to be 0, got %d", q.Size())
	}
}

// Test dequeuing from an empty queue and the synchronization of a subsequent enqueue
func TestDequeueWait(t *testing.T) {
	q := NewThreadSafeQueue()

	go func() {
		_, ok := q.Dequeue()
		if !ok {
			t.Error("Dequeue failed when it should have succeeded")
		}
	}()

	// Allow some time for the Dequeue goroutine to start and block
	time.Sleep(100 * time.Millisecond)
	q.Enqueue(42)
}

// Test that the queue maintains correct order of elements
func TestOrdering(t *testing.T) {
	q := NewThreadSafeQueue()
	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}
	for i := 0; i < 10; i++ {
		item, ok := q.Dequeue()
		if !ok || item != i {
			t.Errorf("Expected to dequeue %d, got %v", i, item)
		}
	}
}

// Test enqueuing and dequeuing multiple types of values
func TestEnqueueDequeueDifferentTypes(t *testing.T) {
	q := NewThreadSafeQueue()
	values := []interface{}{42, "hello", 3.14, true}

	for _, v := range values {
		q.Enqueue(v)
	}

	for _, expected := range values {
		item, ok := q.Dequeue()
		if !ok || item != expected {
			t.Errorf("Expected to dequeue %v, got %v", expected, item)
		}
	}
}

// Test that the queue size reflects concurrent enqueue and dequeue operations accurately
func TestSizeWithConcurrentOperations(t *testing.T) {
	q := NewThreadSafeQueue()
	const count = 1000
	var enqueued, dequeued int32

	enqueue := func() {
		for i := 0; i < count; i++ {
			q.Enqueue(i)
			atomic.AddInt32(&enqueued, 1)
		}
	}

	dequeue := func() {
		for i := 0; i < count; i++ {
			q.Dequeue()
			atomic.AddInt32(&dequeued, 1)
		}
	}

	go enqueue()
	go enqueue()
	go dequeue()
	go dequeue()

	time.Sleep(1 * time.Second) // Allow time for the goroutines to complete

	expectedSize := enqueued - dequeued
	if int(expectedSize) != q.Size() {
		t.Errorf("Expected size to be %d, got %d", expectedSize, q.Size())
	}
}

// Test that multiple Dequeue operations will all eventually succeed after enough Enqueue operations
func TestMultipleDequeueWaits(t *testing.T) {
	q := NewThreadSafeQueue()
	const count = 3
	done := make(chan bool, count)

	for i := 0; i < count; i++ {
		go func() {
			_, ok := q.Dequeue()
			if !ok {
				t.Error("Dequeue failed when it should have succeeded")
			}
			done <- true
		}()
	}

	// Allow some time for the Dequeue goroutines to start and block
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < count; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < count; i++ {
		<-done
	}
}

// Test that Dequeue will not proceed when the queue is consistently empty
func TestDequeueDoesNotProceedWhenEmpty(t *testing.T) {
	q := NewThreadSafeQueue()
	done := make(chan bool)

	go func() {
		_, ok := q.Dequeue()
		if ok {
			t.Error("Dequeue succeeded when it should have failed due to an empty queue")
		}
		done <- true
	}()

	// Allow some time for the Dequeue goroutine to start
	time.Sleep(100 * time.Millisecond)

	// Since we have not enqueued anything, the Dequeue should continue to block
	select {
	case <-done:
		t.Error("Dequeue should not have completed")
	case <-time.After(500 * time.Millisecond): // Allow time to ensure that Dequeue is still blocking
	}
}

// Test that a very large number of Enqueue and Dequeue operations succeed without error
func TestHighVolume(t *testing.T) {
	q := NewThreadSafeQueue()
	const count = 10000

	for i := 0; i < count; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < count; i++ {
		_, ok := q.Dequeue()
		if !ok {
			t.Errorf("Dequeue failed on iteration %d", i)
		}
	}
}

// Test that the correct number of items are enqueued and dequeued
func TestEnqueueDequeueCount(t *testing.T) {
	q := NewThreadSafeQueue()
	const count = 10

	for i := 0; i < count; i++ {
		q.Enqueue(i)
	}

	if q.Size() != count {
		t.Errorf("Expected size to be %d, got %d", count, q.Size())
	}

	for i := 0; i < count; i++ {
		q.Dequeue()
	}

	if q.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", q.Size())
	}
}
