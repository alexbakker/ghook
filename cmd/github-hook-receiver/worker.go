package main

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Worker struct {
	Queue *TaskQueue
	m     sync.RWMutex
	busy  bool
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-time.After(1 * time.Second):
			if task := w.Queue.Pop(); task != nil {
				w.setBusy(true)

				cmd := exec.CommandContext(ctx, "/bin/sh", "-c", task.Command)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Printf("task failed: %s\n", err)
				}

				w.setBusy(false)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) Busy() bool {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.busy
}

func (w *Worker) setBusy(busy bool) {
	w.m.Lock()
	defer w.m.Unlock()
	w.busy = busy
}
