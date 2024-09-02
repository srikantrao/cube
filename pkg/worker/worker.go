package worker

import (
	"cube/pkg/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

//Worker represents a worker that can process tasks
//High Level goals of a Worker:
//- Run tasks as Docker containers
//- Accept tasks to run from a manager
//- Provide relevant statistics to the manager for the purpose of task scheduling
//- Kep track of tasks and their state

type Worker struct {
	name      string
	Queue     *queue.Queue
	DB        map[uuid.UUID]*task.Task
	TaskCount int
}

// NewWorker creates a new instance of a Worker
func NewWorker(name string) *Worker {
	return &Worker{
		name:      name,
		Queue:     queue.New(),
		DB:        make(map[uuid.UUID]*task.Task),
		TaskCount: 0,
	}
}

func (w *Worker) CollectStats() {
}

func (w *Worker) RunTask() {
}

func (w *Worker) StopTask() {
}

func (w *Worker) StartTask() {
}
