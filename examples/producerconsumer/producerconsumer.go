package main

import (
	"fmt"
	"time"

	queue "github.com/sandeepkv93/threadsafequeue"
)

func producer(q *queue.ThreadSafeQueue) {
	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}
}

func consumer(q *queue.ThreadSafeQueue) {
	for i := 0; i < 10; i++ {
		item, _ := q.Dequeue()
		fmt.Println("Consumed:", item)
	}
}

func main() {
	q := queue.NewThreadSafeQueue()

	go producer(q)
	go consumer(q)

	time.Sleep(time.Second) // Allow time for producer and consumer to complete
}
