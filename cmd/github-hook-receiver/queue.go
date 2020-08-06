package main

import (
	"sync"
)

type Task struct {
	Command string
}

type TaskQueue struct {
	tasks []*Task
	m     sync.Mutex
}

func (q *TaskQueue) Push(cmd *Task) {
	q.m.Lock()
	defer q.m.Unlock()

	q.tasks = append(q.tasks, cmd)
}

func (q *TaskQueue) Pop() (task *Task) {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.tasks) > 0 {
		task, q.tasks = q.tasks[0], q.tasks[1:]
		return
	}

	return nil
}

func (q *TaskQueue) Size() int {
	q.m.Lock()
	defer q.m.Unlock()

	return len(q.tasks)
}
