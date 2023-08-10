package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	queue "github.com/sandeepkv93/threadsafequeue"
)

func main() {
	q := queue.NewThreadSafeQueue()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			// Sleep for a random amount of time between 100-200ms to simulate work
			time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)

			q.Enqueue(i)
			fmt.Println("Enqueued", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			// Sleep for a random amount of time between 100-200ms to simulate work
			time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)

			item, _ := q.Dequeue()
			fmt.Println("Dequeued", item)
		}
	}()

	wg.Wait()
}
