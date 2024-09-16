package manager

import (
	"cube/pkg/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

//Manager represents a manager that can manage tasks
//High Level goals of a Manager:
//- Accept requests from users to start and stop tasks
//- Schedule tasks onto workers
//- Keep track of tasks, their state and which machines they are running on

type Manager struct {
	PendingTasks  queue.Queue
	TaskDB        map[string]*task.Task
	EventDB       map[string]*task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID // Map of tasks assigned to each worker
	TaskWorkerMap map[uuid.UUID]string   // Map of worker assigned to each task
}

// SelectWorker selects the worker best suited to run a task based on -
// - Task requirements
// - Worker statistics and availability
func (m *Manager) SelectWorker() {
}

func (m *Manager) UpdateTasks() {
}

func (m *Manager) SendWork() {
}
