package main

/*
RWMutex 特点
1 读写分离（读读并行）：RWMutex 区分读锁和写锁。多个协程可以同时持有读锁（RLock），但同一时间只能有一个协程持有写锁（Lock）。
2 不可重入：和 Mutex 一样，RWMutex 也不可重入。
3 公平性：写优先原则：为了防止写操作被饿死，一旦有协程申请了写锁，后续新来的读锁都会被阻塞，直到写锁完成。
4 适用场景：RWMutex 适用于读多写少的场景，可以显著提升读操作的并发性能。

RWMutex 的底层结构:互斥锁+读写队列+读写协程数量+读等待
它的源码在 sync/rwmutex.go 中，结构如下：
type RWMutex struct {
    w           Mutex  // 互斥锁：用于保护写操作之间的互斥
    writerSem   uint32 // 写等待队列：用于读操作结束后唤醒写操作
    readerSem   uint32 // 读等待队列：用于写操作结束后唤醒读操作
    readerCount int32  // 当前正在读的协程数量（负数表示有写操作在排队）
    readerWait  int32  // 写操作开始前，还需要等待多少个读操作完成
}

底层原理：利用了 readerCount 这个变量：
1 正常读：readerCount 是正数，表示有多少个读协程。
2 写锁进入时：它会把 readerCount 减去一个巨大的常量 2^30，让它变成一个很大的负数。
3 写锁释放时：它会把 readerCount 加回去那个巨大的常量 2^30，恢复成正常的读协程数量，唤醒等待中的读锁操作
3 新读锁判断：当新来的读锁看到 readerCount 是负数，就知道“哦，有个写大佬在排队或写呢”，于是乖乖去 readerSem 队列睡觉。

RWMutex核心方法：
- Lock方法用于在写操作时获取写锁，会阻塞等待当前未释放的写锁。当处于写锁状态时，新的读操作将会阻塞等待。
- Unlock方法用于释放写锁。
- RLock方法用于在读操作时获取读锁，会阻塞等待当前写锁的释放。如果锁处于读锁状态，当前协程也能获取读锁。
- RUnlock方法用于释放读锁。
- RLocker方法用于获取一个Locker接口的对象，调用其Lock方法时会调用RLock方法，调用Unlock方法时会调用RUnlock方法。

四个核心规则
1 读读并行：多个协程可以同时 RLock()，不阻塞。
2 读写互斥：有人在读时，Lock()（写锁）会阻塞；有人在写时，RLock()（读锁）会阻塞。
3 写写互斥：同一时间只能有一个 Lock() 成功，其他 Lock() 必须排队。
4 写优先原则：为了防止写操作被饿死，一旦有协程申请了写锁，后续新来的读锁都会被阻塞，直到写锁完成。

panic 示例
lock 和 unlock 未配对使用，会引发 panic。
*/

import (
	"fmt"
	"sync"
	"time"
)

type Config struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Config) Get(key string) string {
	c.mu.RLock() // 1. 获取读锁：多个协程可以同时执行到这里
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Config) Set(key, value string) {
	c.mu.Lock() // 2. 获取写锁：此时所有读、写操作都会被排斥
	defer c.mu.Unlock()
	c.data[key] = value
}

func main() {
	cfg := &Config{data: make(map[string]string)}

	// 模拟高频读
	for i := 0; i < 5; i++ {
		go func(id int) {
			for {
				_ = cfg.Get("ip")
				fmt.Printf("Reader %d is reading...\n", id)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// 模拟偶尔写
	go func() {
		for {
			cfg.Set("ip", "192.168.1.1")
			fmt.Println("--- Writer: Config updated ---")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	select {} // 阻塞主协程
}
