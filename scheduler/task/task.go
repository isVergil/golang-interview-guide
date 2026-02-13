package task

import (
	"time"
)

// TaskStatus 表示任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Task 表示一个任务
type Task struct {
	ID         string     `json:"id"`
	Data       string     `json:"data"`
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Retries    int        `json:"retries"`
	MaxRetries int        `json:"max_retries"`
	WorkerID   string     `json:"worker_id,omitempty"`
}

// TaskQueue 任务队列接口
type TaskQueue interface {
	// Push 添加任务到队列
	Push(task *Task) error
	// Pop 从队列中取出任务
	Pop() (*Task, error)
	// Update 更新任务状态
	Update(task *Task) error
	// Get 获取任务
	Get(id string) (*Task, error)
}

// InMemoryTaskQueue 内存中的任务队列实现
type InMemoryTaskQueue struct {
	tasks map[string]*Task
}

// NewInMemoryTaskQueue 创建新的内存任务队列
func NewInMemoryTaskQueue() *InMemoryTaskQueue {
	return &InMemoryTaskQueue{
		tasks: make(map[string]*Task),
	}
}

// Push 实现 TaskQueue 接口
func (q *InMemoryTaskQueue) Push(task *Task) error {
	q.tasks[task.ID] = task
	return nil
}

// Pop 实现 TaskQueue 接口
func (q *InMemoryTaskQueue) Pop() (*Task, error) {
	for _, task := range q.tasks {
		if task.Status == TaskStatusPending {
			task.Status = TaskStatusRunning
			task.UpdatedAt = time.Now()
			return task, nil
		}
	}
	return nil, nil
}

// Update 实现 TaskQueue 接口
func (q *InMemoryTaskQueue) Update(task *Task) error {
	q.tasks[task.ID] = task
	return nil
}

// Get 实现 TaskQueue 接口
func (q *InMemoryTaskQueue) Get(id string) (*Task, error) {
	if task, exists := q.tasks[id]; exists {
		return task, nil
	}
	return nil, nil
}
