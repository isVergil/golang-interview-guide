package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// ---------------------------------------------------------------
// Twitter Snowflake 算法 Go 实现
//
// 64 位 ID 结构（从高位到低位）:
//   1 bit  - 符号位（始终为 0）
//  41 bits - 毫秒级时间戳（相对于 epoch，可用 ~69 年）
//  10 bits - 机器 ID（0 ~ 1023，支持 1024 个节点）
//  12 bits - 序列号（同一毫秒内自增，0 ~ 4095）
//
// 单节点每毫秒可生成 4096 个 ID，即 QPS 上限约 409.6 万
// ---------------------------------------------------------------

const (
	epoch          int64 = 1700000000000 // 自定义起始时间 2023-11-14T22:13:20Z
	nodeBits       uint8 = 10           // 机器 ID 位数
	sequenceBits   uint8 = 12           // 序列号位数
	nodeMax        int64 = -1 ^ (-1 << nodeBits)      // 1023
	sequenceMask   int64 = -1 ^ (-1 << sequenceBits)  // 4095
	nodeShift            = sequenceBits                // 左移 12 位
	timestampShift       = nodeBits + sequenceBits     // 左移 22 位
)

// Node 代表一个雪花 ID 生成节点
type Node struct {
	mu        sync.Mutex
	nodeID    int64
	timestamp int64
	sequence  int64
}

// NewNode 创建一个雪花节点，nodeID 范围 [0, 1023]
func NewNode(nodeID int64) (*Node, error) {
	if nodeID < 0 || nodeID > nodeMax {
		return nil, fmt.Errorf("node ID must be between 0 and %d", nodeMax)
	}
	return &Node{nodeID: nodeID}, nil
}

// Generate 生成一个全局唯一的雪花 ID（线程安全）
func (n *Node) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()

	// 时钟回拨检测
	if now < n.timestamp {
		// 生产环境可选择阻塞等待或告警，这里选择等待恢复
		for now < n.timestamp {
			time.Sleep(time.Millisecond)
			now = time.Now().UnixMilli()
		}
	}

	if now == n.timestamp {
		// 同一毫秒内，序列号自增
		n.sequence = (n.sequence + 1) & sequenceMask
		if n.sequence == 0 {
			// 序列号溢出，阻塞到下一毫秒
			for now <= n.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 新的毫秒，序列号重置
		n.sequence = 0
	}

	n.timestamp = now

	return (now-epoch)<<timestampShift | n.nodeID<<int64(nodeShift) | n.sequence
}

// ParseID 反解雪花 ID 各字段（调试用）
func ParseID(id int64) (timestamp time.Time, nodeID int64, sequence int64) {
	ms := (id >> timestampShift) + epoch
	nodeID = (id >> int64(nodeShift)) & nodeMax
	sequence = id & sequenceMask
	timestamp = time.UnixMilli(ms)
	return
}
