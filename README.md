# Thread-Safe Queue Library

This Go library provides a simple and efficient implementation of a thread-safe FIFO (first-in-first-out) queue. It's suitable for scenarios where you need to enqueue and dequeue items in a concurrent environment.

## Installation

To install the Thread-Safe Queue library, you can run:

```shell
go get github.com/sandeepkv93/thread-safe-queue
```

## Usage

Here's a simple guide on how to consume the Thread-Safe Queue library in your project.

### Importing the Library

Import the library using the path where the library is hosted:

```go
import "github.com/yourusername/queue"
```

### Creating a New Queue

To create a new instance of a thread-safe queue:

```go
q := queue.NewThreadSafeQueue()
```

### Enqueueing Items

To enqueue items into the queue:

```go
q.Enqueue(42)
q.Enqueue("Hello World!")
```

### Dequeueing Items

To dequeue items from the queue:

```go
item, ok := q.Dequeue()
if ok {
    // Use the dequeued item
}
```

If the queue is empty, the Dequeue method will block until an item is enqueued.

### Checking if the Queue is Empty

To check if the queue is empty:

```go
if q.IsEmpty() {
    // Queue is empty
}
```

### Checking the Size of the Queue

To check the size of the queue:

```go
size := q.Size()
```

## Examples

### Producer-Consumer Example

Here's an example of how you might use this queue in a producer-consumer scenario:

```go
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
```

## Contributing

Feel free to open issues or submit pull requests if you find any bugs or have suggestions for improvements.
