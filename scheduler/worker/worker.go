package worker

import (
	"encoding/json"
	"go-test/scheduler/task"
	"log"
	"net/http"
	"sync"
	"time"
)

// Worker 表示一个工作节点
type Worker struct {
	ID        string
	Port      string
	Status    string
	LastPing  time.Time
	mu        sync.RWMutex
	taskQueue chan *task.Task
}

// NewWorker 创建新的工作节点
func NewWorker(id, port string) *Worker {
	return &Worker{
		ID:        id,
		Port:      port,
		Status:    "ready",
		LastPing:  time.Now(),
		taskQueue: make(chan *task.Task, 100),
	}
}

// StartWorker 启动工作节点
func StartWorker(port string) {
	w := NewWorker("worker-"+port, port)

	// 启动任务处理协程
	go w.processTasks()

	// 注册 HTTP 处理函数
	http.HandleFunc("/run", w.handleTask)
	http.HandleFunc("/ping", w.handlePing)
	http.HandleFunc("/status", w.handleStatus)

	log.Printf("Worker %s listening on port %s", w.ID, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// processTasks 处理任务队列
func (w *Worker) processTasks() {
	for t := range w.taskQueue {
		log.Printf("Worker %s processing task %s", w.ID, t.ID)

		// 更新任务状态
		t.Status = task.TaskStatusRunning
		t.WorkerID = w.ID
		t.UpdatedAt = time.Now()

		// 模拟任务处理
		time.Sleep(time.Second * 10)

		// 更新任务状态
		t.Status = task.TaskStatusCompleted
		t.UpdatedAt = time.Now()

		log.Printf("Worker %s completed task %s", w.ID, t.ID)
	}
}

// handleTask 处理任务请求
func (w *Worker) handleTask(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t task.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// 将任务加入队列
	w.taskQueue <- &t

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{
		"status": "accepted",
		"worker": w.ID,
	})
}

// handlePing 处理心跳请求
func (w *Worker) handlePing(rw http.ResponseWriter, r *http.Request) {
	w.mu.Lock()
	w.LastPing = time.Now()
	w.mu.Unlock()

	rw.WriteHeader(http.StatusOK)
}

// handleStatus 处理状态查询请求
func (w *Worker) handleStatus(rw http.ResponseWriter, r *http.Request) {
	w.mu.RLock()
	status := map[string]interface{}{
		"id":        w.ID,
		"status":    w.Status,
		"last_ping": w.LastPing,
		"queue_len": len(w.taskQueue),
	}
	w.mu.RUnlock()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(status)
}
