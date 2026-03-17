package base

import (
	"fmt"
	"testing"
)

/*
map 本质是一个哈希表，采用链地址法解决冲突，渐进式扩容来保证性能。
底层数据结构实际上是一个名为 hmap 的结构体，它管理着一组被称为 bmap 的桶。
type hmap struct {
	count     int    // 当前 map 中的元素个数（调用 len() 直接返回这个值）
	flags     uint8  // 状态标志（正在写入、正在扩容等） 写操作会将其置为 hashWriting。如果另一个协程发现此位已置，直接 Panic。
	B         uint8  // 桶的数量的对数，即桶的个数 = 2^B
	noverflow uint16 // 溢出桶的大致数量
	hash0     uint32 // 哈希种子，在 make 时随机生成，用于防止哈希碰撞攻击
	buckets    unsafe.Pointer // 指向当前的桶数组（2^B 个桶）
	oldbuckets unsafe.Pointer // 扩容时指向旧的桶数组
	nevacuate  uintptr        // 扩容进度计数器，表示已搬迁到新桶的数量
	extra *mapextra // 存储溢出桶的指针，优化 GC
}

桶结构 bmap (Bucket)：
type bmap struct {
    // tophash 存储 8 个 key 哈希值的高 8 位
    // 查找时先比对高 8 位，不一致直接跳过，提高速度
	tophash [8]uint8
	keys     [8]keytype   // 连续存储 8 个键
	values   [8]valuetype // 连续存储 8 个值
	overflow uintptr  // 指向下一个溢出桶的指针（拉链法）。
}

面试题：
1 map 底层数据结构是怎样的？
 -底层数据结构实际上是一个名为 hmap 的结构体，它管理着一组 bmap 结构的桶。

2 桶内部结构是怎样的？
 -桶里最多存 8 个键值对
 -Toppash：桶里有个数组存了 8 个哈希值的高 8 位，用来快速匹配 Key，提高查找效率。
 -内存对齐：桶内部是把 8 个 Key 放一起，8 个 Value 放一起，分开放是为了节省内存

3 map 遍历是有序还是无序的？为什么要这么设计？
 -map 的遍历顺序是随机的。即便桶里的数据没变，Go 每次遍历都会从一个随机的桶和随机的偏移量开始。防止开发者依赖遍历顺序，因为扩容后顺序一定会变。

4 map 是线程安全的吗？
 -非线程安全：并发读写会直接 panic: fatal error: concurrent map read and map write。
 -多个协程同时读写同一个，写操作会将其flags置为 hashWriting。如果另一个协程发现此位已置，直接 Panic。

5 delete 会释放内存吗？
 -delete 操作只是把位置标记为可用，即把bmap的tophash标记为 empty，并不会释放内存给操作系统。如果 map 曾存过百万数据现在只剩几条，内存依然占用很大。唯一的办法是置为 nil 重新 make。

6 可以对 key 取地址吗？
 -不可取址
 -&m["key"] 是不合法的。因为 map 会扩容搬迁，Key 的地址是不固定的。

7 map 什么时候扩容？怎么扩容？
 -扩容机制：扩容不是瞬间完成的，而是 渐进式 的（每次增删改时搬迁 1~2 个桶），避免瞬间卡顿。
 -翻倍扩容（增量扩容）：当 负载因子 > 6.5（平均每个桶装了超过 6.5 个元素）时，桶数量翻倍。
 -等量扩容（整理扩容）：当负载因子不大，但溢出桶太多时。这通常是因为大量删除后，虽然元素少了，但数据存储很稀疏，扩容是为了让数据更紧凑。

8 为什么用存储桶数量的对数 B ？包含溢出桶吗？
 -B 不包含溢出桶，只表示主桶的数量（2^B 个）。
  -1 提高效率：当我们要确定一个 Key 落在哪个桶时，标准做法是 hash % 桶数。因为桶的数量总是 2 的幂次方，使用对数可以快速计算桶索引：hash(key) & (2^B - 1)，避免了除法运算，提高效率。
	- hash key之后用“低 B 位”找桶 (Bucket Index)，用“高 8 位”找槽位 (TopHash)
  -2 节省空间：使用对数 B 可以节省存储空间。直接存储桶的数量需要更多的位数，而存储对数只需要较少的位数，节省了内存开销，翻倍直接加 B+1，桶数量扩容两倍，节省空间。

*/

// TestMapConcurrentPanic 陷阱题：并发读写
func TestMapConcurrentPanic(t *testing.T) {
	m := make(map[int]int)

	// 面试官问：这段代码会发生什么？
	// 回答：fatal error: concurrent map writes。
	// 注意：这不是普通的 error，是无法被 recover 的崩溃。
	go func() {
		for {
			m[1] = 1
		}
	}()
	go func() {
		for {
			_ = m[1]
		}
	}()
}

// TestMapDeleteMemory 性能题：删除元素会释放内存吗？
func TestMapDeleteMemory(t *testing.T) {
	m := make(map[int][100]int)
	// 往里塞 10000 个大对象
	for i := 0; i < 10000; i++ {
		m[i] = [100]int{}
	}

	// 删除所有 key
	for i := 0; i < 10000; i++ {
		delete(m, i)
	}

	// 面试官问：此时内存释放了吗？
	// 回答：没有。Map 的桶占用的内存不会缩小。
	// 解决：如果要释放内存，只能把原来的 map 指向 nil，等 GC 回收。
	fmt.Println("Map 元素删完了，但底层桶依然占用内存")
}
