package master

import (
	"bytes"
	"encoding/json"
	"go-test/scheduler/task"
	"log"
	"net/http"
	"sync"
	"time"
)

// Master 表示调度器主节点
type Master struct {
	Workers   map[string]*WorkerInfo
	TaskQueue task.TaskQueue
	mu        sync.RWMutex
	stopChan  chan struct{}
}

// WorkerInfo 存储 worker 节点信息
type WorkerInfo struct {
	ID        string
	Address   string
	Status    string
	LastPing  time.Time
	TaskCount int
}

// NewMaster 创建新的调度器主节点
func NewMaster() *Master {
	return &Master{
		Workers:   make(map[string]*WorkerInfo),
		TaskQueue: task.NewInMemoryTaskQueue(),
		stopChan:  make(chan struct{}),
	}
}

// Start 启动调度器
func (m *Master) Start() {
	// 启动 worker 健康检查
	go m.healthCheck()

	// 启动任务调度
	go m.scheduleTasks()

	// 启动 HTTP 服务
	http.HandleFunc("/register", m.handleRegister)
	http.HandleFunc("/tasks", m.handleTasks)
	http.HandleFunc("/workers", m.handleWorkers)

	log.Println("Master started")
}

// Stop 停止调度器
func (m *Master) Stop() {
	close(m.stopChan)
}

// healthCheck 定期检查 worker 健康状态
func (m *Master) healthCheck() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkWorkers()
		case <-m.stopChan:
			return
		}
	}
}

// checkWorkers 检查所有 worker 的健康状态
func (m *Master) checkWorkers() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, worker := range m.Workers {
		// 如果超过 10 秒没有心跳，认为 worker 已离线
		if now.Sub(worker.LastPing) > 10*time.Second {
			log.Printf("Worker %s is offline", id)
			delete(m.Workers, id)
		}
	}
}

// scheduleTasks 调度任务到 worker
func (m *Master) scheduleTasks() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.dispatchTasks()
		case <-m.stopChan:
			return
		}
	}
}

// dispatchTasks 分发任务到 worker
func (m *Master) dispatchTasks() {
	m.mu.RLock()
	if len(m.Workers) == 0 {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	// 获取待处理任务
	curTask, err := m.TaskQueue.Pop()
	if err != nil || curTask == nil {
		return
	}

	// 选择负载最低的 worker
	worker := m.selectWorker()
	if worker == nil {
		// 如果没有可用 worker，将任务重新放回队列
		curTask.Status = task.TaskStatusPending
		m.TaskQueue.Push(curTask)
		return
	}

	// 发送任务到 worker
	if err := m.sendTaskToWorker(worker, curTask); err != nil {
		log.Printf("Failed to send task to worker %s: %v", worker.ID, err)
		curTask.Status = task.TaskStatusPending
		m.TaskQueue.Push(curTask)
	}
}

// selectWorker 选择负载最低的 worker
func (m *Master) selectWorker() *WorkerInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var selected *WorkerInfo
	minTasks := -1

	for _, worker := range m.Workers {
		if minTasks == -1 || worker.TaskCount < minTasks {
			selected = worker
			minTasks = worker.TaskCount
		}
	}

	return selected
}

// sendTaskToWorker 发送任务到 worker
func (m *Master) sendTaskToWorker(worker *WorkerInfo, t *task.Task) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(t); err != nil {
		return err
	}

	resp, err := http.Post("http://"+worker.Address+"/run", "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	// 更新 worker 任务计数
	m.mu.Lock()
	worker.TaskCount++
	m.mu.Unlock()

	return nil
}

// handleRegister 处理 worker 注册
func (m *Master) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var info struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	m.Workers[info.ID] = &WorkerInfo{
		ID:       info.ID,
		Address:  info.Address,
		Status:   "ready",
		LastPing: time.Now(),
	}
	m.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

// handleTasks 处理任务相关请求
func (m *Master) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// 创建新任务
		var t task.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
		t.Status = task.TaskStatusPending

		if err := m.TaskQueue.Push(&t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodGet:
		// 获取任务列表
		m.mu.RLock()
		workers := make([]*WorkerInfo, 0, len(m.Workers))
		for _, w := range m.Workers {
			workers = append(workers, w)
		}
		m.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workers)
	}
}

// handleWorkers 处理 worker 相关请求
func (m *Master) handleWorkers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	m.mu.RLock()
	workers := make([]*WorkerInfo, 0, len(m.Workers))
	for _, w := range m.Workers {
		workers = append(workers, w)
	}
	m.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
}
