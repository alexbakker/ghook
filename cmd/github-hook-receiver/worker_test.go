package main

import (
	"context"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker := Worker{Queue: new(TaskQueue)}
	go worker.Run(ctx)

	worker.Queue.Push(&Task{"sleep 2 && echo task 1 is done"})
	worker.Queue.Push(&Task{"sleep 2 && echo task 2 is done"})
	worker.Queue.Push(&Task{"sleep 2 && echo task 3 is done"})

	for worker.Queue.Size() > 0 || worker.Busy() {
		time.Sleep(1 * time.Second)
	}
}
