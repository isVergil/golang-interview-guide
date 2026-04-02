# Redis 核心理论

---

## 1. Redis 是什么

Redis（Remote Dictionary Server）是基于内存的键值存储数据库，常用作缓存、消息队列、分布式锁。

核心特点：
- **纯内存操作**：读写性能极高，单机 10 万+ QPS
- **单线程模型**：避免锁竞争（6.0 后 I/O 线程多线程，命令执行仍单线程）
- **丰富的数据结构**：不只是 key-value，支持 String、Hash、List、Set、Sorted Set 等
- **持久化**：支持 RDB 和 AOF，数据不会因重启丢失
- **原子操作**：所有命令都是原子的，天然适合并发场景

### 1.1 为什么 Redis 这么快？

1. **纯内存操作**：数据在内存中读写，不需要磁盘 IO
2. **单线程**：无锁竞争、无上下文切换开销
3. **IO 多路复用**：用 epoll 处理大量并发连接，单线程也能扛住高并发
4. **高效数据结构**：底层用 SDS、ziplist、skiplist 等优化过的数据结构

### 1.2 单线程为什么不是瓶颈？

Redis 的瓶颈不在 CPU，而在网络和内存。单个命令执行时间是纳秒~微秒级别，CPU 根本不是瓶颈。网络 IO 才是真正的瓶颈，所以 Redis 6.0 引入了多线程 IO（读写网络数据并行化），但命令执行依然是单线程。

```
Redis 6.0 之前：                  Redis 6.0 之后：
  网络读 → 命令执行 → 网络写        网络读（多线程） → 命令执行（单线程） → 网络写（多线程）
  ↑ 全部单线程                     ↑ IO 多线程，执行单线程
```

---

## 2. 底层数据结构

Redis 对外暴露的是 String、Hash、List、Set、ZSet 这些**数据类型**，但底层实际存储用的是另一套**数据结构**。同一个数据类型在不同条件下会切换不同的底层结构，这是 Redis 省内存又快的关键。

先建立直觉——每个底层结构的本质是什么：

```
数据结构          本质（一句话）                              你可以想象成
──────────────────────────────────────────────────────────────────────────
SDS             带长度字段的字节数组                         Go 的 []byte（切片）
ziplist         一整块连续内存，元素紧挨着排列               一个紧凑的 byte[]，手动编码了多个元素
listpack        ziplist 的改良版，同样是连续内存              修复了 ziplist 的 bug 的 byte[]
hashtable       数组 + 单链表（拉链法）                     Java 的 HashMap（数组桶 + 链表解决冲突）
intset          有序整数数组                                一个 int[]，二分查找
skiplist        多层单向链表 + 最底层双向链表                 给有序链表加了"快速通道"的索引层
quicklist       双向链表，但每个节点是一个 listpack            LinkedList<listpack>
```

```
数据类型            底层数据结构（根据数据量自动切换）
─────────────────────────────────────────────────
String          →  int / embstr / raw（SDS）
Hash            →  listpack（小） / hashtable（大）
List            →  listpack（小） / quicklist（大）
Set             →  intset（纯整数） / hashtable（其他）
Sorted Set      →  listpack（小） / skiplist + hashtable（大）
```

### 2.1 SDS（Simple Dynamic String）

> **本质：带头部元信息的字节数组**，类似 Go 的 `[]byte` 切片——有 len（已用长度）、cap（总容量）、底层 buf 数组。不是链表，就是一段连续内存。

Redis 没有直接用 C 的 `char*`，而是自己实现了 SDS。这是 Redis 最基础的数据结构，几乎所有 key 和 String 类型的 value 都用它。

**C 字符串的问题**：

```c
char *s = "hello";
```

- 获取长度要遍历到 `\0`，O(n)
- 修改前不知道空间够不够，容易缓冲区溢出
- 不能存二进制数据（遇到 `\0` 就截断）

**SDS 的结构**：

```c
struct sdshdr {
    uint32_t len;      // 已使用长度
    uint32_t alloc;    // 分配的总长度（不含头和 \0）
    unsigned char flags; // SDS 类型（sdshdr5/8/16/32/64）
    char buf[];        // 实际数据
};
```

```
示例：存 "hello"

+------+-------+-------+-------------------+
| len=5 | alloc=8 | flags | h e l l o \0 _ _ _ |
+------+-------+-------+-------------------+
                          ↑ buf[]    ↑ 预留空间
```

**SDS 相比 C 字符串的优势**：

| 特性 | C 字符串 | SDS |
|------|---------|-----|
| 获取长度 | O(n)，遍历到 \0 | O(1)，直接读 len 字段 |
| 缓冲区安全 | 手动管理，易溢出 | 自动扩容，API 安全 |
| 二进制安全 | 不支持（\0 截断） | 支持，用 len 判断结束 |
| 内存分配 | 每次修改都分配 | 预分配 + 惰性释放 |

**空间预分配策略**：

修改 SDS 时，不是只分配刚好需要的空间，而是多分配一些：
- `len < 1MB`：额外分配 `len` 大小的空间（翻倍）
- `len >= 1MB`：额外分配 1MB

这样连续追加数据时减少内存重新分配次数。

**惰性空间释放**：

缩短 SDS 时不立刻回收内存（alloc 不变，len 减小），后续追加数据时可以直接使用。需要手动调用 `sdsRemoveFreeSpace` 才真正释放。

**SDS 的编码类型**：

根据字符串长度选择不同的头大小，进一步省内存：

| 类型 | 头大小 | 最大长度 |
|------|--------|---------|
| sdshdr5 | 1 字节 | 31 字节（不再使用） |
| sdshdr8 | 3 字节 | 255 字节 |
| sdshdr16 | 5 字节 | 64KB |
| sdshdr32 | 9 字节 | 4GB |
| sdshdr64 | 17 字节 | 理论无限 |

### 2.2 ziplist（压缩列表）

> **本质：一整块连续内存（字节数组），所有元素紧挨着排列**。不是链表、没有指针，更像是一个手动序列化的 `byte[]`，靠偏移量在里面定位每个元素。可以类比一个紧凑的"数组"，但每个元素长度不固定。

> Redis 7.0 之后 ziplist 已被 listpack 替代，但面试仍然常考，需要掌握。

ziplist 是一块**连续内存**，用于在数据量小的时候替代链表和哈希表，极致省内存。

**内存布局**：

```
<zlbytes> <zltail> <zllen> <entry1> <entry2> ... <entryN> <zlend>
  4字节    4字节    2字节                                   1字节(0xFF)
```

- `zlbytes`：整个 ziplist 占用的字节数
- `zltail`：最后一个 entry 的偏移量（支持从尾部反向遍历）
- `zllen`：entry 数量（< 65535 时有效，超过需遍历计算）
- `zlend`：结束标记 `0xFF`

**entry 结构**：

```
<prevlen> <encoding> <data>
```

- `prevlen`：前一个 entry 的长度（1 或 5 字节）
  - 前一个 entry < 254 字节 → 1 字节存
  - 前一个 entry >= 254 字节 → 5 字节存（1 字节标记 0xFE + 4 字节长度）
- `encoding`：当前 entry 的类型和长度
- `data`：实际数据

**优点**：
- 连续内存，CPU 缓存友好，没有指针开销
- 小数据量时比链表和哈希表省内存得多

**缺点——连锁更新（cascade update）**：

这是 ziplist 最大的问题。假设多个连续的 entry 长度都在 250~253 字节之间（prevlen 用 1 字节存）：

```
[entry1: 250B] [entry2: 252B] [entry3: 253B] ...
  prevlen=1B     prevlen=1B     prevlen=1B
```

如果在最前面插入一个 260 字节的 entry，entry1 的 prevlen 要从 1 字节变成 5 字节（因为前面的 entry 是 260 > 254），entry1 总长度增加了 4 字节变成 254 字节，导致 entry2 的 prevlen 也得从 1 字节变成 5 字节……像多米诺骨牌一样连锁更新。

最坏情况 O(n²)，这就是 Redis 7.0 引入 listpack 替代 ziplist 的原因。

### 2.3 listpack（紧凑列表）

> **本质：和 ziplist 一样是一块连续内存（字节数组）**，只是改了 entry 的编码方式。你可以理解为 ziplist v2——同样的设计思路，修掉了连锁更新的 bug。

Redis 7.0 引入，替代 ziplist，解决了**连锁更新**问题。

**和 ziplist 的关键区别**：

```
ziplist entry:   <prevlen> <encoding> <data>     ← 存前一个的长度（连锁更新根源）
listpack entry:  <encoding> <data> <backlen>     ← 存自己的长度
```

listpack 的 entry 不再记录前一个 entry 的长度，而是在末尾记录**自己的总长度**。反向遍历时读当前 entry 末尾的 backlen 就知道前一个 entry 在哪，不依赖前一个 entry 的大小。

这样不管你在哪里插入/修改数据，都不会影响后面的 entry，彻底解决了连锁更新。

**内存布局**：

```
<total_bytes> <num_elements> <entry1> <entry2> ... <entryN> <end>
   4字节         2字节                                       1字节(0xFF)
```

### 2.4 hashtable（哈希表）

> **本质：数组 + 单链表（拉链法）**。想象一排桶（数组），key 经过哈希函数算出桶编号，丢进对应的桶里。多个 key 落到同一个桶就用**单链表**串起来（头插法）。和 Java 的 `HashMap`（JDK7 之前）结构一样。Redis 额外维护了**两个**这样的哈希表，用于渐进式搬迁（rehash）。

```
直觉图：

  table（数组）
  ┌───┐
  │ 0 │ → entry → entry → NULL      ← 单链表（拉链法解决冲突）
  ├───┤
  │ 1 │ → NULL
  ├───┤
  │ 2 │ → entry → NULL
  ├───┤
  │ 3 │ → entry → entry → entry → NULL
  └───┘
```

Redis 的哈希表实现，用于 Hash 类型大数据量和 Set 类型。

**结构**：

```c
typedef struct dictht {
    dictEntry **table;      // 哈希桶数组
    unsigned long size;     // 桶数量（2 的幂次）
    unsigned long sizemask; // size - 1，用于取模
    unsigned long used;     // 已有元素数量
} dictht;

typedef struct dict {
    dictht ht[2];           // 两个哈希表，用于渐进式 rehash
    long rehashidx;         // rehash 进度，-1 表示没在 rehash
} dict;
```

**哈希冲突**：用链地址法，冲突的元素用链表串起来（头插法）。

**渐进式 rehash**：

当负载因子（used / size）达到阈值时需要扩容。Redis 不会一次性搬迁所有数据（那样会阻塞很久），而是**渐进式**搬迁：

```
1. 分配新的 ht[1]，大小是 ht[0] 的 2 倍
2. rehashidx 设为 0
3. 每次 CRUD 操作时，顺便把 ht[0].table[rehashidx] 这个桶的所有元素迁移到 ht[1]
4. rehashidx++
5. 期间查询两个表都查，写入只写 ht[1]
6. 全部迁移完成后，ht[0] 指向 ht[1]，ht[1] 置空，rehashidx = -1
```

这样把一次大的 rehash 分摊到每次操作里，不会有长时间的阻塞。

**扩容触发条件**：
- 没有在执行 bgsave/bgrewriteaof 时：负载因子 >= 1 就扩容
- 正在执行 bgsave/bgrewriteaof 时：负载因子 >= 5 才扩容（避免 COW 时扩容增加内存）

**缩容触发条件**：负载因子 < 0.1

### 2.5 intset（整数集合）

> **本质：有序整数数组**。就是一个 `int[]`，元素从小到大排列。查找用二分查找 O(log n)，插入要移动元素保持有序 O(n)。没有链表、没有哈希，纯数组。

当 Set 的元素全是整数且数量不多时，用 intset 存储，比 hashtable 省内存。

**结构**：

```c
typedef struct intset {
    uint32_t encoding;   // 编码方式：int16、int32、int64
    uint32_t length;     // 元素数量
    int8_t contents[];   // 有序整数数组
} intset;
```

**内存布局**（以 int16 为例）：

```
encoding=INT16, length=4
contents: [1, 5, 10, 20]   ← 有序排列，二分查找 O(log n)
```

**升级机制**：

假设当前是 int16 编码（范围 -32768 ~ 32767），如果插入一个 65535（超出 int16），整个数组升级为 int32：

```
升级前 (int16): [1, 5, 10, 20]       每个元素 2 字节，共 8 字节
升级后 (int32): [1, 5, 10, 20, 65535] 每个元素 4 字节，共 20 字节
```

升级过程需要重新分配内存并扩展每个元素的位宽。**只能升级不能降级**。

**优点**：紧凑数组、CPU 缓存友好、内存占用小
**缺点**：查找 O(log n)，插入 O(n)（需要移动元素保持有序），只适合小规模

### 2.6 skiplist（跳表）

> **本质：多层有序单向链表，最底层额外有反向指针（可双向遍历）**。想象一条从左到右排好序的链表，然后在上面再叠几层"快速通道"——高层跨度大、跳得远（像坐高铁），底层跨度小、一步一步走（像走路）。查找时从顶层往下找，越来越精确，实现 O(log n) 的性能。不是数组、不是树，就是多层链表。

```
直觉图：

Level 3（高铁）:  head ──────────────────→ 25 ────────────→ NULL
Level 2（快车）:  head ──────→ 10 ──────→ 25 ────→ 40 ──→ NULL
Level 1（步行）:  head → 5 → 10 → 15 → 20 → 25 → 30 → 40 → 50 → NULL
                                                ↑ ← backward 反向指针（只有最底层有）
```

Sorted Set 底层的核心数据结构。本质是多层有序链表，通过"跳跃"实现 O(log n) 的查找、插入、删除。

**为什么用跳表不用平衡树（AVL/红黑树）？**

- 实现简单：跳表比红黑树代码量少得多，容易理解和维护
- 范围查询友好：跳表最底层就是有序链表，`ZRANGEBYSCORE` 找到起点后直接遍历
- 插入性能稳定：不需要旋转操作

**结构**：

```
Level 3:  head ──────────────────────→ 25 ────────────────→ NULL
Level 2:  head ──────→ 10 ──────────→ 25 ──────→ 40 ────→ NULL
Level 1:  head → 5 → 10 → 15 → 20 → 25 → 30 → 40 → 50 → NULL
```

查找 30 的过程：
1. 从最高层开始，head → 25（25 < 30，继续）→ NULL（超了，降层）
2. Level 2：25 → 40（40 > 30，降层）
3. Level 1：25 → 30（找到！）

只走了 4 步，而不是从头遍历 6 步。

**Redis 跳表的实现细节**：

```c
typedef struct zskiplistNode {
    sds ele;                        // 成员值
    double score;                   // 分数
    struct zskiplistNode *backward; // 后退指针（只有 Level 1 有）
    struct zskiplistLevel {
        struct zskiplistNode *forward; // 前进指针
        unsigned long span;            // 跨度（用于计算排名）
    } level[];                      // 层级数组，随机高度
} zskiplistNode;

typedef struct zskiplist {
    struct zskiplistNode *header, *tail;
    unsigned long length;           // 节点数量
    int level;                      // 最高层级
} zskiplist;
```

**随机层高**：每个节点的层高是随机决定的，概率 p=0.25：
- Level 1：100%
- Level 2：25%
- Level 3：6.25%
- Level 4：1.5625%
- 最高 32 层

平均每个节点 1.33 个指针，内存开销比平衡树小。

**span（跨度）的作用**：

```
Level 2:  head ─(span=2)─→ 10 ─(span=3)─→ 40
Level 1:  head ─(1)→ 5 ─(1)→ 10 ─(1)→ 15 ─(1)→ 20 ─(1)→ 25 ─(1)→ 30 ─(1)→ 40
```

`ZRANK key member` 计算排名时，沿查找路径把经过的 span 加起来就是排名，O(log n)。

**backward（后退指针）的作用**：

每个节点只有**一个** backward 指针，指向最底层（Level 1）的前一个节点。注意不是每层都有，只有最底层有。

```
forward 方向（每层都有，可以跳跃前进）：
Level 2:  head ──────→ 10 ──────────→ 25 ──────→ 40
Level 1:  head → 5 → 10 → 15 → 20 → 25 → 30 → 40 → 50

backward 方向（只有最底层，每次只能退一步）：
Level 1:  head ← 5 ← 10 ← 15 ← 20 ← 25 ← 30 ← 40 ← 50
                                                       ↑ tail
```

forward 可以跳着走（高层一次跨好几个节点），但 backward **只能一步一步往回走**，因为每个节点只存了一个 backward 指针。

**backward 在实际 Redis 命令中的作用**：

所有"倒序"操作都依赖 backward 指针。Redis 先通过 forward（从高层到底层快速定位），找到终点后，沿着 backward 一步步往回走：

| 命令 | 作用 | 怎么用 backward |
|------|------|----------------|
| `ZREVRANGE key 0 9` | 按 score 倒序取 Top 10 | 从 tail 开始，backward 走 10 步 |
| `ZREVRANGEBYSCORE key 100 60` | score 100→60 之间的成员 | forward 定位到 score≤100 的位置，backward 走到 score<60 为止 |
| `ZREVRANK key member` | 倒序排名 | forward 定位到 member，从 tail 算到该节点的距离 |
| `ZPOPMAX key` | 弹出分数最高的 | 直接从 tail 取，backward 更新 tail |

**具体例子——排行榜 Top 10**：

```bash
# 排行榜：score 越高排名越靠前
ZADD ranking 100 alice 90 bob 95 charlie 80 dave 85 eve

# 取 Top 3（倒序）
ZREVRANGE ranking 0 2 WITHSCORES
# 结果：alice(100), charlie(95), bob(90)
```

底层执行过程：
```
1. 跳表的 tail 指向 score 最大的节点 alice(100)
2. alice.backward → charlie(95)    ← 第 1 步
3. charlie.backward → bob(90)      ← 第 2 步
4. 取够 3 个，返回
```

如果没有 backward 指针，倒序遍历就得每次都从 head 用 forward 走到末尾再往前数，性能会从 O(n) 的简单遍历变成 O(n²)。

**为什么 backward 只有一层不设多层？**

forward 需要多层是为了"跳着找"（查找 O(log n)），但 backward 的使用场景是**已经定位好了，往回遍历**，不需要跳跃，一步步走就行。多层 backward 浪费内存，收益几乎没有。

**Sorted Set 为什么同时用 skiplist + hashtable？**

```c
typedef struct zset {
    dict *dict;       // hashtable：member → score 的映射
    zskiplist *zsl;   // skiplist：按 score 排序
} zset;
```

- `ZSCORE key member`：需要 O(1) 查 score → 走 hashtable
- `ZRANGE key 0 10`：需要按 score 排序遍历 → 走 skiplist
- `ZRANK key member`：先 hashtable 查 score，再 skiplist 查排名

两个数据结构通过指针共享同一份数据（member 和 score），不会重复占内存。

### 2.7 quicklist（快速列表）

> **本质：双向链表，但每个节点不是单个元素，而是一个 listpack（小数组）**。可以类比 `LinkedList<listpack>`——大结构是双向链表（前后指针，可以从头尾快速操作），每个链表节点里装的是一块连续内存的 listpack。综合了链表的增删优势和连续内存的缓存友好。

```
直觉图：

 ┌──────────┐    ┌──────────┐    ┌──────────┐
 │ listpack │ ⇄  │ listpack │ ⇄  │ listpack │     ← 双向链表
 │[a, b, c] │    │[d, e, f] │    │[g, h]    │     ← 每个节点里是连续内存
 └──────────┘    └──────────┘    └──────────┘
```

Redis 3.2 引入，List 类型的底层实现。本质是一个**双向链表，但每个节点不是单个元素，而是一个 listpack（旧版是 ziplist）**。

```
quicklist: node1 ⇄ node2 ⇄ node3 ⇄ node4

node1: [listpack: elem1, elem2, elem3]
node2: [listpack: elem4, elem5, elem6]
node3: [listpack: elem7, elem8]
node4: [listpack: elem9]
```

**为什么不直接用链表或 ziplist/listpack？**

- 纯链表：每个元素一个节点 + 前后指针，内存浪费严重
- 纯 ziplist/listpack：连续内存，数据量大了增删效率差（需要移动大量数据）
- quicklist 折中：链表的节点是 listpack，小规模连续内存（缓存友好），大规模链表串联（增删快）

**配置参数**：

```
list-max-listpack-size -2   # 每个 listpack 节点最大 8KB（默认）
list-compress-depth 0       # 两端不压缩的节点数（0=不压缩）
```

`list-compress-depth`：对中间不常访问的节点用 LZF 压缩，进一步省内存。头尾节点不压缩（LPUSH/RPOP 频繁访问）。

### 2.8 redisObject

> **本质：一个 16 字节的包装头**，Redis 所有的 value 在内存里都先套一层这个结构，通过里面的 `type`（数据类型）和 `encoding`（底层编码）两个字段告诉 Redis 该用什么方式操作底层数据。类比面向对象里的"基类"，所有具体数据结构都通过它统一管理。

Redis 所有的 key-value 在内部都包装成 `redisObject`：

```c
typedef struct redisObject {
    unsigned type:4;       // 数据类型（STRING/LIST/HASH/SET/ZSET）
    unsigned encoding:4;   // 底层编码（int/embstr/raw/listpack/skiplist 等）
    unsigned lru:24;       // LRU/LFU 时间戳（淘汰策略用）
    int refcount;          // 引用计数
    void *ptr;             // 指向实际数据结构的指针
} robj;
```

一个 redisObject 占 16 字节。通过 `type` 和 `encoding` 的组合，Redis 可以在运行时根据数据特征选择最优的底层结构。

**查看某个 key 的编码**：

```bash
SET num 123
OBJECT ENCODING num      # "int"

SET name "hello"
OBJECT ENCODING name     # "embstr"

SET big_str "aaa...（超过44字节）"
OBJECT ENCODING big_str  # "raw"

HSET user:1 name alice
OBJECT ENCODING user:1   # "listpack"（数据少时）

# 往 Hash 里塞超过 128 个 field
OBJECT ENCODING user:1   # "hashtable"（自动升级）
```

### 2.9 编码转换阈值汇总

| 数据类型 | 小数据编码 | 大数据编码 | 转换条件 |
|---------|-----------|-----------|---------|
| String | int | - | 值是整数 |
| String | embstr | raw | 长度 > 44 字节 |
| Hash | listpack | hashtable | field 数 > 128 或单个值 > 64 字节 |
| List | listpack | quicklist | 元素数 > 128 或单个值 > 64 字节 |
| Set | intset | hashtable | 元素不全是整数，或数量 > 512 |
| ZSet | listpack | skiplist + hashtable | 元素数 > 128 或单个值 > 64 字节 |

以上阈值可通过 `redis.conf` 调整（如 `hash-max-listpack-entries`）。

**关于 ZSet 大数据编码的"共生"关系**：

表里写的 `skiplist + hashtable` 不是"二选一"，而是**同时存在、共生**的。当 ZSet 数据量超过阈值后，Redis 会同时维护一个 skiplist 和一个 hashtable，指向同一份数据：

```
ZSet 小数据量：               ZSet 大数据量：
  listpack（单一结构）          skiplist（按 score 排序）
                              +
                              hashtable（member → score 映射）
                              ↑ 两个结构共生，共享同一份数据
```

为什么要同时维护两个？因为单一结构无法同时满足所有操作的性能要求：

```
操作                       走哪个结构         为什么
─────────────────────────────────────────────────────
ZSCORE key member         hashtable          O(1) 直接查 member 对应的 score
ZRANGE key 0 10           skiplist           按 score 顺序遍历
ZRANGEBYSCORE key 60 100  skiplist           按 score 范围遍历
ZRANK key member          skiplist           沿查找路径累加 span 得排名
ZADD key 95 alice         两个都写           skiplist 插入排序位置 + hashtable 插 member→score
ZREM key alice            两个都删           skiplist 删节点 + hashtable 删映射
```

如果只有 skiplist，`ZSCORE` 就要 O(log n) 去跳表里找；如果只有 hashtable，`ZRANGE` 就没法按 score 排序遍历。两个共生才能让所有操作都高效。

注意：小数据量时用 listpack 就够了（遍历也很快，内存小），不需要这种双结构。

**业务场景选型指南**：

| 业务需求 | 推荐类型 | 为什么 |
|---------|---------|-------|
| 缓存单个值（JSON、计数器、锁） | **String** | 最简单，支持原子 INCR，可设过期时间 |
| 缓存对象，需要读写单个字段 | **Hash** | 不用整体序列化/反序列化，单字段 HGET/HSET |
| 缓存对象，整体读写为主 | **String**（存 JSON） | 一次 GET 拿全部数据，简单直接 |
| 消息队列、任务队列 | **List** / **Stream** | LPUSH+BRPOP 简单队列；Stream 支持消费者组、ACK |
| 去重、集合运算（交并差） | **Set** | SINTER 共同好友，SRANDMEMBER 抽奖 |
| 排行榜、按分数排序 | **ZSet** | 天然排序，ZREVRANGE 取 Top N |
| 延迟队列 | **ZSet** | score 存执行时间戳，ZRANGEBYSCORE 取到期任务 |
| 签到、在线状态（布尔值大量用户） | **Bitmap** | 1 bit/用户，百万用户只占 125KB |
| UV 统计（大规模去重计数） | **HyperLogLog** | 固定 12KB，误差 0.81%，亿级去重 |
| 附近的人、地理围栏 | **GEO** | 底层 ZSet + GeoHash，GEORADIUS 范围查 |
| 轻量消息队列（需要 ACK） | **Stream** | 消费者组 + 消息确认，比 List 更可靠 |

**容易选错的场景**：

```
❌ 用 String 存对象，频繁改单个字段
   → 每次要 GET 整个 JSON → 改字段 → SET 回去（读-改-写，不原子）
   ✅ 用 Hash，直接 HSET user:1 age 26

❌ 用 List 做排行榜
   → List 没有按分数排序的能力，排序要自己做
   ✅ 用 ZSet，天然按 score 排序

❌ 用 Set 做排行榜
   → Set 无序，没有分数概念
   ✅ 用 ZSet

❌ 用 ZSet 做简单去重
   → ZSet 有 score 开销，内存多占
   ✅ 用 Set，纯去重不需要排序

❌ 用 String 做签到记录（一天一个 key）
   → 百万用户 × 365 天 = 3.65 亿个 key
   ✅ 用 Bitmap，一个用户一年只占 46 字节
```

---

## 3. 数据类型与命令

### 3.1 五大基础类型

#### String（字符串）

最基础的类型，可以存字符串、整数、浮点数、二进制（如图片的 base64）。

底层编码：
- `int`：值是整数且 <= 2^63-1 时，直接用 long 存，最省内存
- `embstr`：字符串长度 <= 44 字节，SDS 和 redisObject 一次内存分配
- `raw`：字符串长度 > 44 字节，SDS 和 redisObject 两次分配

```bash
SET name "alice"
GET name                    # "alice"
SET counter 100
INCR counter                # 101（原子自增）
INCRBY counter 50           # 151
SETNX lock:order "1"        # 不存在才设置（分布式锁基础）
SET token "abc" EX 3600     # 设置并指定过期时间 1 小时
MSET k1 v1 k2 v2           # 批量设置（减少网络往返）
MGET k1 k2                 # 批量获取
```

典型场景：缓存、计数器、分布式锁、Session 存储、限流（INCR + EXPIRE）

#### Hash（哈希）

field-value 映射，适合存储对象。

底层编码：
- `listpack`（Redis 7.0+，旧版用 `ziplist`）：field 数量 <= 128 且每个值 <= 64 字节时使用，内存紧凑
- `hashtable`：超过阈值后转为标准哈希表

```bash
HSET user:1 name alice age 25 email alice@test.com
HGET user:1 name              # "alice"
HGETALL user:1                # 获取所有 field-value
HMSET user:1 name bob age 30  # 批量设置
HINCRBY user:1 age 1          # 某个 field 自增
HDEL user:1 email             # 删除某个 field
HLEN user:1                   # field 数量
```

典型场景：用户信息缓存、购物车（user:cart → {商品ID: 数量}）、对象缓存

**Hash vs String 存对象**：

```
方案一：String 存 JSON
  SET user:1 '{"name":"alice","age":25}'
  → 读写整个对象，修改单个字段需要读-改-写

方案二：Hash 存字段
  HSET user:1 name alice age 25
  → 可以单独读写某个字段，无需序列化/反序列化
```

需要频繁读写单个字段 → Hash。整体读写为主 → String JSON。

#### List（列表）

有序的字符串列表，支持头尾插入和弹出。

底层编码：
- `listpack`：元素少且小时
- `quicklist`（Redis 7.0+）：ziplist 节点组成的双向链表，兼顾内存和性能

```bash
LPUSH queue task1 task2      # 左侧插入（头部）
RPUSH queue task3            # 右侧插入（尾部）
LPOP queue                   # 左侧弹出
RPOP queue                   # 右侧弹出
LRANGE queue 0 -1            # 获取所有元素
LLEN queue                   # 长度
BRPOP queue 30               # 阻塞弹出，最多等 30 秒（消息队列基础）
```

典型场景：消息队列（LPUSH + BRPOP）、最新动态列表、评论列表

#### Set（集合）

无序、不重复的字符串集合。

底层编码：
- `intset`：元素全是整数且数量 <= 512 时，紧凑数组
- `hashtable`：其他情况

```bash
SADD tags:article:1 go redis mysql
SMEMBERS tags:article:1       # 获取所有成员
SISMEMBER tags:article:1 go   # 是否存在
SCARD tags:article:1          # 成员数量
SINTER tags:article:1 tags:article:2  # 交集
SUNION tags:article:1 tags:article:2  # 并集
SDIFF tags:article:1 tags:article:2   # 差集
SRANDMEMBER tags:article:1 2  # 随机取 2 个（抽奖）
SPOP tags:article:1 1         # 随机弹出 1 个（不可重复抽奖）
```

典型场景：标签系统、共同好友（交集）、抽奖、去重

#### Sorted Set / ZSet（有序集合）

每个成员关联一个 score（分数），按 score 排序，成员唯一。

底层编码：
- `listpack`：元素少且小时
- `skiplist + hashtable`：跳表做排序，哈希表做 O(1) 查找

```bash
ZADD ranking 100 alice 90 bob 95 charlie
ZRANGE ranking 0 -1 WITHSCORES    # 按 score 正序
ZREVRANGE ranking 0 2 WITHSCORES  # 按 score 倒序 Top 3
ZSCORE ranking alice              # 查某个成员的分数
ZRANK ranking alice               # 查排名（从 0 开始）
ZINCRBY ranking 10 alice          # 加分
ZRANGEBYSCORE ranking 90 100      # 按分数范围查
ZCARD ranking                     # 成员总数
```

典型场景：排行榜、延迟队列（score 存时间戳）、滑动窗口限流

### 3.2 特殊类型

#### BitMap（位图）

本质是 String，按 bit 位操作。

```bash
SETBIT sign:user:1:202401 0 1    # 1月1日签到
SETBIT sign:user:1:202401 1 1    # 1月2日签到
GETBIT sign:user:1:202401 0      # 查 1月1日是否签到
BITCOUNT sign:user:1:202401      # 统计签到天数
```

典型场景：签到系统、在线状态、布隆过滤器基础

#### HyperLogLog

基数估算，用固定 12KB 内存估算集合中不重复元素数量，误差约 0.81%。

```bash
PFADD page:uv:20240101 user1 user2 user3
PFADD page:uv:20240101 user1 user4        # user1 重复不计
PFCOUNT page:uv:20240101                   # 4（去重计数）
PFMERGE page:uv:total page:uv:20240101 page:uv:20240102  # 合并多天
```

典型场景：UV 统计、大规模去重计数

#### GEO（地理位置）

底层用 Sorted Set 实现，score 是 GeoHash 编码。

```bash
GEOADD city 116.40 39.90 beijing 121.47 31.23 shanghai
GEODIST city beijing shanghai km         # 两地距离
GEORADIUS city 116.40 39.90 500 km       # 500km 内的城市
```

典型场景：附近的人、门店搜索

#### Stream（Redis 5.0+）

专门为消息队列设计的数据结构，支持消费者组、消息确认、持久化。

```bash
XADD mystream * name alice age 25        # 发消息
XREAD COUNT 10 STREAMS mystream 0        # 读消息
XGROUP CREATE mystream mygroup 0         # 创建消费者组
XREADGROUP GROUP mygroup consumer1 COUNT 1 STREAMS mystream >  # 消费
XACK mystream mygroup 1234567890-0       # 确认
```

典型场景：轻量级消息队列（替代 List 方案）

---

## 4. 持久化

### 4.1 RDB（Redis Database）

将某一时刻的全量数据以二进制快照保存到磁盘。

**触发方式**：
- `save`：主线程执行，阻塞所有请求（生产禁用）
- `bgsave`：fork 子进程执行，利用 COW（Copy-On-Write）不阻塞主进程
- 自动触发：配置 `save 900 1`（900 秒内有 1 次写入就触发）

**bgsave 流程**：

```
1. 主进程 fork 出子进程（fork 本身会短暂阻塞，内存越大越慢）
2. 子进程拿到内存页表的副本
3. 子进程遍历内存数据写入 RDB 文件
4. 期间主进程继续处理请求
5. 主进程修改数据时触发 COW，复制对应的内存页
6. 子进程完成后通知主进程，替换旧的 RDB 文件
```

**优点**：
- 文件紧凑，恢复速度快
- 适合全量备份、灾难恢复

**缺点**：
- 两次快照之间的数据可能丢失
- fork 在数据量大时有短暂阻塞（与内存大小正相关）

### 4.2 AOF（Append Only File）

将每条写命令追加到文件末尾，重启时重放命令恢复数据。

**写回策略**（`appendfsync`）：

| 策略 | 含义 | 性能 | 安全 |
|------|------|------|------|
| `always` | 每条命令都 fsync | 最慢 | 最安全，最多丢 1 条 |
| `everysec` | 每秒 fsync 一次 | 折中（推荐） | 最多丢 1 秒 |
| `no` | 由 OS 决定何时 fsync | 最快 | 可能丢较多 |

**AOF 重写**：AOF 文件会越来越大，重写机制把冗余命令合并。比如对同一个 key 做了 100 次 INCR，重写后变成一条 SET。

```
重写前 AOF：
  SET counter 0
  INCR counter    ← 重复 100 次
  
重写后 AOF：
  SET counter 100  ← 合并为一条
```

重写也是 fork 子进程执行，期间主进程的新命令写入**重写缓冲区**，子进程完成后追加到新 AOF 文件。

### 4.3 RDB + AOF 混合持久化（Redis 4.0+）

`aof-use-rdb-preamble yes`（默认开启）

AOF 重写时，前半部分写 RDB 格式（全量快照），后半部分追加 AOF 格式（增量命令）。兼顾恢复速度和数据安全性。

**恢复优先级**：AOF 优先于 RDB（AOF 数据更完整）

### 4.4 怎么选？

| 场景 | 推荐 |
|------|------|
| 纯缓存，丢了无所谓 | 不开持久化 |
| 允许丢几分钟数据 | RDB |
| 数据尽量不丢 | AOF everysec |
| 生产环境推荐 | **RDB + AOF 混合**（默认配置就是） |

---

## 5. 内存管理与淘汰策略

### 5.1 过期删除策略

Redis 不会在 key 过期的瞬间立刻删除，而是用两种策略配合：

**惰性删除**：访问 key 的时候检查是否过期，过期了就删。省 CPU 但可能有大量过期 key 占内存。

**定期删除**：每 100ms 随机抽查一批设了过期时间的 key，删除已过期的。如果过期比例超过 25%，继续抽查。

两者配合：定期删除处理大部分，惰性删除兜底。但还是可能有漏网之鱼，所以需要内存淘汰策略。

### 5.2 内存淘汰策略

当内存达到 `maxmemory` 时，新写入命令会触发淘汰。

| 策略 | 说明 |
|------|------|
| `noeviction` | 不淘汰，写入直接报错（默认） |
| `allkeys-lru` | 从所有 key 中淘汰最久没访问的（**推荐**） |
| `allkeys-lfu` | 从所有 key 中淘汰访问频率最低的 |
| `allkeys-random` | 从所有 key 中随机淘汰 |
| `volatile-lru` | 从设了过期时间的 key 中淘汰最久没访问的 |
| `volatile-lfu` | 从设了过期时间的 key 中淘汰访问频率最低的 |
| `volatile-random` | 从设了过期时间的 key 中随机淘汰 |
| `volatile-ttl` | 从设了过期时间的 key 中淘汰 TTL 最短的 |

**生产推荐**：
- 纯缓存场景 → `allkeys-lru`（我们的 docker-compose 配的就是这个）
- 缓存 + 持久化混合 → `volatile-lru`（只淘汰设了过期时间的，不动持久化数据）

### 5.3 Redis 的 LRU 是近似 LRU

Redis 不是严格 LRU（维护链表开销太大），而是随机采样近似 LRU：每次需要淘汰时随机采样 N 个 key（默认 5 个），淘汰其中最久没访问的。

Redis 4.0 引入的 LFU（Least Frequently Used）更智能，基于访问频率而非时间，能识别出"偶尔被访问一次但不常用"的冷数据。

---

## 6. 缓存问题

### 6.1 缓存穿透

**问题**：查询一个**数据库里根本不存在**的数据，缓存自然也没有，每次请求都穿透缓存直接打到数据库。

```
正常请求：  请求 id=1 → 缓存 hit → 返回          ← 缓存挡住了
穿透请求：  请求 id=-1 → 缓存 miss → DB 查不到 → 不写缓存 → 返回空
           下次 id=-1 → 缓存还是 miss → DB 又查 → 又查不到 → ...
           ↑ 每次都打到 DB，缓存形同虚设
```

如果是恶意攻击（大量不存在的 id），会直接把数据库打挂。

**方案一：缓存空值**

核心思路：DB 查不到也往缓存写一个空值，下次直接从缓存返回空。

```
步骤：
1. 请求 id=-1
2. GET user:-1 → miss
3. 查 DB → 查不到
4. SET user:-1 "" EX 60        ← 缓存空值，设短过期（60秒）
5. 返回空

下次请求 id=-1：
1. GET user:-1 → hit（值是 ""）
2. 判断是空值 → 直接返回空，不查 DB  ← 缓存挡住了
```

缺点：如果攻击方每次用不同的随机 id（id=-1, -2, -3, ...），每个 id 都会写一个空值 key，缓存里全是垃圾数据。

**方案二：布隆过滤器（推荐）**

核心思路：在缓存前面加一层"门卫"，数据库里有的 id 提前注册到布隆过滤器，请求进来先问门卫"这个 id 存在吗"，不存在的直接拦截。

```
步骤：
1. 启动时 / 数据变更时，把所有合法 id 加入布隆过滤器
2. 请求 id=-1 进来
3. 布隆过滤器判断 → "一定不存在" → 直接返回空，不查缓存、不查 DB
4. 请求 id=1 进来
5. 布隆过滤器判断 → "可能存在" → 查缓存 → miss → 查 DB → 写缓存 → 返回
```

```
请求 → 布隆过滤器 → 不存在 → 直接返回空（拦截在最外层）
                  → 可能存在 → 查缓存 → 查 DB
```

布隆过滤器特点：说"不存在"就**一定不存在**，说"存在"可能**误判**（概率极低，可控）。误判顶多多查一次 DB，不会有问题。

**方案三：两者结合（生产推荐）**

布隆过滤器拦截大部分非法请求 + 缓存空值兜底漏网的。

### 6.2 缓存击穿

**问题**：某个**热点 key** 过期的瞬间，大量并发请求同时发现缓存失效，全部涌向数据库查同一条数据。

```
时间线：
  t=0    热点 key "product:1" 过期
  t=0.01 请求A GET → miss → 查 DB
  t=0.02 请求B GET → miss → 查 DB     ← 缓存还没回填
  t=0.03 请求C GET → miss → 查 DB     ← 大量请求同时打 DB
  ...
  t=0.05 请求A 查到结果，SET 回缓存
  ↑ 在 A 回填之前，所有请求都打到了 DB
```

**方案一：互斥锁（牺牲少量可用性，保证数据最新）**

核心思路：只放一个请求去查 DB，其他请求等它查完后从缓存取。

```
步骤（以请求 A、B、C 同时到达为例）：

请求 A：
  1. GET product:1 → miss
  2. SETNX lock:product:1 "uuid-A" EX 10  → 成功，拿到锁
  3. 查 DB → 得到数据
  4. SET product:1 {data} EX 3600          → 写缓存
  5. DEL lock:product:1                     → 释放锁
  6. 返回数据

请求 B（几乎同时）：
  1. GET product:1 → miss
  2. SETNX lock:product:1 → 失败，没拿到锁
  3. sleep 50ms
  4. 重试 GET product:1 → hit（A 已经回填了）  ← 从缓存拿
  5. 返回数据

请求 C（同上）：
  1~4. 同 B，等 A 回填后从缓存拿
```

只有请求 A 查了 DB，B 和 C 都从缓存拿。代价是 B、C 多等了几十毫秒。

**方案二：逻辑过期（牺牲短暂一致性，保证高可用，推荐）**

核心思路：Redis 的 key **永不过期**（不设 TTL），在 value 里自己存一个过期时间。发现逻辑过期后，抢锁异步刷新，抢不到锁的直接返回旧数据。**没有任何请求会阻塞等待**。

```
缓存的 value 格式：
  {"data": {"name": "iPhone", "price": 5999}, "expire_at": 1700000000}
   ↑ 实际业务数据                              ↑ 逻辑过期时间戳
```

```
步骤（以请求 A、B、C 同时到达为例）：

请求 A：
  1. GET product:1 → hit（key 永不过期，一定能 hit）
  2. 读 expire_at，发现 < 当前时间 → 逻辑过期了
  3. SETNX lock:refresh:product:1 "uuid-A" EX 10  → 成功，拿到锁
  4. 开一个异步 goroutine：
     ├── 查 DB → 得到最新数据
     ├── SET product:1 {"data": {新数据}, "expire_at": now+3600}
     └── DEL lock:refresh:product:1
  5. 当前请求不等异步结果，直接返回旧数据 ← 立即返回，不阻塞
  
请求 B（几乎同时）：
  1. GET product:1 → hit
  2. 读 expire_at → 逻辑过期了
  3. SETNX lock:refresh:product:1 → 失败，A 已经拿了锁
  4. 直接返回旧数据 ← 不等待、不重试、不查 DB

请求 C（同上）：
  1~4. 同 B，直接返回旧数据

几百毫秒后，A 的异步 goroutine 完成：
  缓存已刷新为最新数据
  后续所有请求拿到的都是新数据
```

```
整体流程图：

  GET key → hit → 检查 expire_at
                   ├── 没过期 → 直接返回数据
                   └── 过期了 → SETNX 抢锁
                                ├── 拿到锁 → 异步刷新 + 返回旧数据
                                └── 没拿到 → 直接返回旧数据（别人在刷新了）
```

**两种方案对比**：

| 维度 | 互斥锁 | 逻辑过期 |
|------|--------|---------|
| 请求是否阻塞 | 是，抢不到锁的要 sleep 重试 | 否，所有请求立即返回 |
| 数据一致性 | 强，等到新数据才返回 | 弱，过期后有短暂窗口返回旧数据 |
| DB 压力 | 只有 1 个请求查 DB | 只有 1 个请求查 DB |
| 复杂度 | 低 | 中（需要维护 expire_at 字段） |
| 适用场景 | 数据一致性要求高 | 高可用优先、允许短暂不一致（推荐） |

**方案三：热点 key 永不过期 + 后台定时刷新**

后台定时任务每隔 N 分钟主动刷新缓存，key 永不过期。最简单但刷新不够及时。

### 6.3 缓存雪崩

**问题**：不是一个 key 的问题，而是**大面积**缓存同时失效（集中过期或 Redis 宕机），海量请求同时打到数据库，DB 直接被压垮。

```
场景一：大量 key 同时过期
  t=0    批量导入数据时统一设了 expire = 3600（1小时后过期）
  t=3600 几万个 key 同时过期
         → 几万个请求同时 miss → 全部打到 DB → DB 崩溃

场景二：Redis 宕机
  Redis 挂了 → 所有请求都 miss → 全部打到 DB → DB 崩溃
```

和缓存击穿的区别：击穿是**一个热点 key**，雪崩是**大面积 key**。

**方案一：过期时间加随机值（防集中过期）**

核心思路：让 key 的过期时间分散开，不要在同一时刻集中过期。

```
步骤：
1. 设置缓存时，不要用固定过期时间
2. expire = baseTime + random(0, 300)

示例：基础过期 1 小时
  key1: expire = 3600 + 120 = 3720 秒
  key2: expire = 3600 + 45  = 3645 秒
  key3: expire = 3600 + 280 = 3880 秒
  ↑ 过期时间被打散到 5 分钟范围内，不会同时过期
```

**方案二：多级缓存（防 Redis 故障）**

核心思路：在 Redis 前面再加一层本地缓存，即使 Redis 挂了，本地缓存还能扛一阵。

```
请求 → 本地缓存（进程内，如 Go bigcache/ristretto）
       ├── hit → 直接返回（最快，不走网络）
       └── miss → Redis
                  ├── hit → 返回 + 回填本地缓存
                  └── miss → DB
                             ├── 查到 → 返回 + 写 Redis + 写本地缓存
                             └── 查不到 → 返回空
```

本地缓存容量小、过期时间短（如 30 秒），只用来扛突发流量。

**方案三：限流降级（兜底）**

核心思路：DB 前面加限流，超过阈值的请求直接返回降级结果（如默认值、错误提示），不让 DB 被打挂。

```
步骤：
1. Redis 全部 miss，请求涌向 DB
2. 限流器判断：当前 QPS > 阈值（如 1000）
3. 超过的请求 → 直接返回降级结果（"系统繁忙"/ 默认数据 / 上次缓存的快照）
4. 放行的请求 → 查 DB → 回填缓存
5. 缓存逐步恢复后，限流器放开
```

**方案四：Redis 高可用（防宕机）**

```
单机 Redis          → 挂了就全挂       ← 不要用
Sentinel 哨兵模式   → 自动故障转移     ← 中小规模
Cluster 集群模式    → 分片 + 高可用    ← 大规模
```

**生产环境通常组合使用**：

```
过期时间加随机（防集中过期）
+ 多级缓存（防 Redis 故障扛流量）
+ 限流降级（兜底保护 DB）
+ Redis Sentinel/Cluster（防单点宕机）
```

### 6.4 缓存与数据库一致性

**经典方案：Cache Aside（旁路缓存）**

Cache Aside 的核心是两条规则：

- **读路径——懒加载（Lazy Loading）**：缓存里没有才去查 DB，查到后回填缓存。不提前预热，按需加载。
- **写路径——先更新 DB，再删除缓存**：写操作不更新缓存，而是删掉缓存，让下次读来触发懒加载。

**读路径（懒加载）完整流程**：

```
1. GET product:1 → hit → 直接返回            ← 缓存命中，最快路径
2. GET product:1 → miss                       ← 缓存没有
3. 查 DB → 得到数据
4. SET product:1 {data} EX 3600               ← 回填缓存
5. 返回数据
```

```
请求 → 查缓存 → hit → 返回
              → miss → 查 DB → 写缓存 → 返回
                       ↑ 懒加载：第一次访问才加载到缓存
```

懒加载的好处：
- 只缓存真正被访问的数据，不浪费内存
- 新数据不需要额外操作，第一次读时自动加载
- 缓存挂了重启后，随着请求自然恢复，不需要手动预热（虽然生产环境通常还是会预热热点数据）

**写路径完整流程**：

```
1. UPDATE product SET price=6999 WHERE id=1   ← 先更新 DB
2. DEL product:1                               ← 再删除缓存
3. 下次读 product:1 → miss → 查 DB → 回填缓存  ← 懒加载拿到最新数据
```

**为什么写操作是删除缓存，而不是更新缓存？——防止脏数据**

如果写操作直接更新缓存，并发写入时会产生脏数据：

```
错误方案：先更新 DB，再更新缓存

  线程 A（改价格为 6999）    线程 B（改价格为 7999）
  ─────────────────────    ─────────────────────
  1. UPDATE DB price=6999
                            2. UPDATE DB price=7999
                            3. SET cache price=7999
  4. SET cache price=6999   ← A 比 B 晚写缓存，把 B 的新值覆盖了

  结果：DB 是 7999（正确），缓存是 6999（脏数据）
  ↑ 缓存和 DB 不一致，后续所有读请求拿到的都是错的
```

改成删除缓存就没这个问题：

```
正确方案：先更新 DB，再删除缓存

  线程 A（改价格为 6999）    线程 B（改价格为 7999）
  ─────────────────────    ─────────────────────
  1. UPDATE DB price=6999
                            2. UPDATE DB price=7999
                            3. DEL cache
  4. DEL cache              ← 删两次也没关系，反正都是删

  下次读：miss → 查 DB → 拿到 7999（最新值） → 写缓存
  ↑ 一定是最新的，因为懒加载从 DB 读
```

删除缓存还有一个好处：缓存可能是经过复杂计算的（比如聚合多张表、排序、过滤），更新缓存要重复这些计算，不如直接删掉让下次读时重新计算。

**为什么是先更新 DB 再删缓存，不是反过来？——也是防止脏数据**

```
错误方案：先删缓存，再更新 DB

  线程 A（写操作）           线程 B（读操作）
  ─────────────────────    ─────────────────────
  1. DEL cache              ← A 先删了缓存
                            2. GET cache → miss
                            3. 查 DB → 拿到旧值（A 还没更新 DB）
                            4. SET cache = 旧值   ← 把旧值写回缓存了！
  5. UPDATE DB = 新值

  结果：DB 是新值，缓存是旧值（脏数据）
  ↑ 在 A 删缓存到更新 DB 的窗口期，B 读到了旧数据并回填
```

```
正确方案：先更新 DB，再删缓存

  线程 A（写操作）           线程 B（读操作）
  ─────────────────────    ─────────────────────
  1. UPDATE DB = 新值
                            2. GET cache → hit（旧值）← 短暂读到旧值
  3. DEL cache
                            4. GET cache → miss → 查 DB → 新值 ← 懒加载拿到最新

  极端情况下 B 在 step 2 拿到了旧值，但窗口极短（A 更新 DB 到删缓存之间）
  ↑ 概率极低，且影响范围只有这一次请求
```

**进一步保证一致性**：

**延迟双删**：防止"先更新 DB 再删缓存"方案中极端情况下的短暂不一致。

```
步骤：
1. UPDATE DB = 新值
2. DEL cache                  ← 第一次删
3. sleep 500ms                ← 等可能的并发读完成回填
4. DEL cache                  ← 第二次删，清掉可能被回填的旧值

为什么等 500ms？
  估算一次"读 DB + 写缓存"的耗时，确保并发读操作已经完成回填
  第二次删除就能把回填的旧值清掉
```

**监听 Binlog（最可靠）**：

```
步骤：
1. 业务代码只负责更新 DB + 删缓存
2. Canal 监听 MySQL Binlog
3. Binlog 有变更 → Canal 异步发消息 → 消费者删除/更新对应的缓存

好处：
  - 即使第一次 DEL cache 失败了（网络抖动），Binlog 兜底会再删一次
  - 解耦，业务代码不需要关心缓存一致性的细节
```

---

## 7. 分布式锁

### 7.1 基本实现

```bash
# 加锁：SET key value NX EX seconds（原子操作）
SET lock:order:123 "uuid-xxx" NX EX 30

# 释放锁：Lua 脚本保证原子性（先判断是不是自己的锁再删）
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
```

关键点：
- **NX**：不存在才设置，保证互斥
- **EX**：设过期时间，防止持锁进程崩溃导致死锁
- **value 用 UUID**：释放时校验是不是自己加的锁，防止误删别人的锁
- **Lua 脚本释放**：GET + DEL 要原子操作，不然中间可能过期被别人拿到锁

### 7.2 存在的问题

**锁过期但业务没执行完**：

加锁时设了 30 秒，但业务执行了 35 秒。第 30 秒锁自动过期，别人拿到锁，两个进程同时执行。

解决：**看门狗机制**（Redisson 的 watchdog）。后台线程定期续期，只要持锁进程还活着就续。

**Redis 主从切换锁丢失**：

主节点加锁成功后还没同步到从节点就挂了，从节点升主后没有这个锁，别人又能加锁成功。

解决：**RedLock 算法**（有争议）。向 N 个独立 Redis 节点加锁，多数成功才算成功。但 Martin Kleppmann 指出 RedLock 在网络分区和时钟漂移场景下不可靠。

**生产建议**：
- 一般场景：单节点 Redis 锁 + 看门狗续期就够了
- 强一致要求：用 etcd 或 ZooKeeper 做分布式锁

---

## 8. 高可用架构

### 8.1 主从复制

最基础的高可用方案：一个 Master 负责写，多个 Slave 负责读。

```
               ┌──→ Slave1（读）
Master（读+写）─┤
               └──→ Slave2（读）
```

**读写分离**：

```
写操作（SET/DEL/INCR 等）  → 只能发到 Master
读操作（GET/HGET/ZRANGE 等）→ 可以发到 Slave（分担 Master 压力）
```

客户端需要自己做读写分离（或通过代理），写操作发 Master，读操作发 Slave。Slave 默认是**只读**的，写操作发到 Slave 会报错。

**Slave 的作用**：

| 作用 | 说明 |
|------|------|
| 读扩展 | 读请求分散到多个 Slave，Master 只处理写，提升整体读吞吐 |
| 数据备份 | Slave 有完整数据副本，Master 磁盘坏了还有 Slave 的数据 |
| 故障恢复 | Master 挂了可以手动把 Slave 提升为新 Master（手动，不自动） |

**同步机制**：

```
首次连接——全量同步：
  1. Slave 发送 PSYNC 命令给 Master
  2. Master 执行 bgsave 生成 RDB 快照
  3. Master 把 RDB 发送给 Slave
  4. Slave 加载 RDB，拥有了 Master 某一时刻的全量数据
  5. 期间 Master 新产生的写命令暂存到 repl_backlog 缓冲区
  6. RDB 传完后，Master 把 repl_backlog 里积压的命令也发给 Slave

后续——增量同步：
  Master 每执行一条写命令，异步发给所有 Slave 重放
  ↑ 注意是异步的，Slave 的数据有短暂延迟（毫秒级）
```

**主从复制的局限**：

- **Master 单点**：Master 挂了，写操作全部不可用，需要手动把 Slave 提升为 Master
- **不能自动故障转移**：运维必须人工介入
- **写不能扩展**：写操作只能走一个 Master，写瓶颈无法解决

### 8.2 Sentinel（哨兵）

在主从复制的基础上，加了一组 Sentinel 进程来**自动监控和故障转移**。解决了主从复制"Master 挂了要手动切换"的问题。

```
Sentinel1 ─── Sentinel2 ─── Sentinel3    ← 哨兵集群（奇数个，至少3个）
    ↓             ↓             ↓
  监控 Master + Slave1 + Slave2           ← 数据节点还是主从架构
```

**读写分离（和主从复制一样）**：

```
写操作 → Master（只有一个 Master）
读操作 → Slave（多个 Slave 分担读压力）
```

Sentinel 本身不转发读写请求，它只负责**监控和切换**。客户端连接 Sentinel 获取当前 Master 地址，然后直连 Master/Slave 进行读写。

**Slave 的作用（比主从复制多了一个）**：

| 作用 | 说明 |
|------|------|
| 读扩展 | 同主从复制，分担读压力 |
| 数据备份 | 同主从复制，数据冗余 |
| **自动顶替 Master** | Master 挂了，Sentinel 自动选一个 Slave 提升为新 Master |

**故障转移完整流程**：

```
1. Sentinel 每秒 PING Master
2. Master 无响应 → Sentinel1 标记"主观下线"（一个哨兵觉得挂了，不一定真挂）
3. Sentinel1 问其他 Sentinel："你们觉得 Master 挂了吗？"
4. 多数 Sentinel（≥ quorum）确认 → "客观下线"（确定真挂了）
5. Sentinel 之间选举一个 leader（Raft 协议）
6. leader 从 Slave 中选一个最优的提升为新 Master
   选择优先级：slave-priority 配置 > 复制偏移量最大的（数据最新）> runid 最小的
7. 其他 Slave 执行 SLAVEOF 新Master，从新 Master 同步数据
8. 通知客户端新 Master 的地址，客户端自动切换连接

整个过程对业务无感知（短暂不可用几秒）
```

**Sentinel 的局限**：

- **写不能扩展**：和主从一样，只有一个 Master 处理写操作
- **存储不能扩展**：所有数据都在一个 Master 里，单机内存是上限
- 适合中小规模（数据量 < 单机内存上限）

### 8.3 Cluster（集群）

Redis 官方分布式方案。和前两种架构最大的区别：**有多个 Master，数据分片存储**。

**数据分片**：16384 个 slot（哈希槽），每个 Master 负责一部分 slot。

```
key → CRC16(key) % 16384 → 得到 slot 编号 → 找到负责这个 slot 的 Master
```

**架构**：

```
Master1 (slot 0-5460)      ⇄  Slave1     ← Master1 的备份
Master2 (slot 5461-10922)  ⇄  Slave2     ← Master2 的备份
Master3 (slot 10923-16383) ⇄  Slave3     ← Master3 的备份
```

注意：这里有 **3 个 Master**（不是 1 个），每个 Master 只存一部分数据。

**读写都在 Master 上进行**：

```
写操作 → 根据 key 算出 slot → 发到对应的 Master
读操作 → 根据 key 算出 slot → 默认也发到对应的 Master

例如：
  SET user:1 alice → CRC16("user:1") % 16384 = 5678 → slot 5678 → Master2 处理
  GET user:1       → 同样算出 slot 5678 → Master2 处理
  SET order:99 ... → CRC16("order:99") % 16384 = 1234 → slot 1234 → Master1 处理
```

Cluster 模式下**读写默认都走 Master**，不像主从/哨兵那样读走 Slave。Slave 默认只做热备，不参与读请求。

**Slave 在 Cluster 中的作用**：

| 作用 | 说明 |
|------|------|
| **故障自动顶替** | Master1 挂了 → Slave1 自动提升为新 Master1，接管 slot 0-5460 |
| 数据备份 | 全量复制 Master 的数据，防数据丢失 |
| 读扩展（可选） | 客户端发 `READONLY` 命令后可以从 Slave 读（需要容忍数据延迟） |

```
正常情况：
  Master1 处理 slot 0-5460 的所有读写
  Slave1  实时同步 Master1 数据，待命

Master1 宕机：
  Cluster 自动把 Slave1 提升为新 Master1
  Slave1（现在是 Master1）接管 slot 0-5460 的读写
  整个过程自动完成，其他 Master 不受影响
```

**MOVED 重定向**：

客户端可以连任意一个 Master，如果 key 不归它管，会告诉你去找谁：

```
客户端 → Master1: GET user:1
Master1: MOVED 5678 192.168.1.2:6379    ← "这个 key 不归我管，去找 Master2"
客户端 → Master2: GET user:1
Master2: "alice"                         ← 正确的节点处理了
```

智能客户端（如 go-redis）会缓存 slot 映射表，后续直接发到正确节点，避免重定向。

**数据路由模式**：

客户端发送命令时，需要知道 key 对应的 slot 在哪个 Master 上。有三种路由方式：

```
方式一：哑客户端（Dummy Client）
  客户端不缓存 slot 映射，随便连一个节点发命令
  如果发错了，收到 MOVED 后再重新发 → 每次可能需要两次请求
  简单但性能差，基本没人用

方式二：智能客户端（Smart Client）← go-redis、Jedis 等都是这种
  客户端启动时拉取一次完整的 slot 映射表（CLUSTER SLOTS 命令）
  后续直接根据本地映射表发到正确节点 → 一次请求搞定
  收到 MOVED 时更新本地映射表

方式三：代理模式（Proxy）
  客户端连代理（如 Twemproxy、Codis），代理负责路由
  客户端不感知集群，像连单机一样用
  多一跳网络延迟，适合不想改客户端的老项目
```

**智能客户端的完整路由流程**（以 go-redis 为例）：

```
第 1 步：启动时获取映射
  客户端 → 任意节点：CLUSTER SLOTS
  节点返回：
    slot 0-5460     → Master1 (192.168.1.1:6379)
    slot 5461-10922 → Master2 (192.168.1.2:6379)
    slot 10923-16383 → Master3 (192.168.1.3:6379)
  客户端缓存这张映射表

第 2 步：发送命令
  SET user:1001 "alice"
  → CRC16("user:1001") % 16384 = 5649
  → 查本地映射：slot 5649 → Master2
  → 直接发给 Master2
  → Master2 返回 OK（一次请求完成）

第 3 步：节点变化时（扩容/缩容/故障转移）
  客户端 → Master2：GET order:500
  Master2：MOVED 12000 192.168.1.3:6379
  → 客户端更新本地映射：slot 12000 → Master3
  → 重新发给 Master3
  → 后续 slot 12000 的请求直接发 Master3

第 4 步：slot 迁移中
  客户端 → Master2：GET migrating:key
  Master2：ASK 8888 192.168.1.3:6379
  → 客户端临时重定向到 Master3（不更新本地映射，因为迁移还没完成）
  → 发 ASKING + GET migrating:key 给 Master3
  
  MOVED vs ASK：
    MOVED = 永久转移，更新本地映射
    ASK   = 临时重定向（迁移中），不更新映射，下次还发旧节点
```

**Cluster 的特殊限制**：

```
❌ 不支持跨 slot 的多 key 操作：
  MGET user:1 order:99          → 报错（两个 key 在不同 slot）

✅ 用 Hash Tag 强制同 slot：
  MGET {user}:1 {user}:99       → 只对 {} 内的部分算 hash，保证同 slot
```

**Gossip 协议与节点数量上限**：

Cluster 里没有中心节点，所有节点通过 **Gossip 协议**互相通信来同步状态信息（谁负责哪些 slot、谁挂了、谁是谁的 Slave）。

Gossip 的工作方式：每个节点每秒随机选几个节点发消息交换信息，像"八卦传播"一样最终所有节点都知道整个集群的状态。

```
每条 Gossip 消息的内容：
  ┌──────────────────────────────────────────────────┐
  │ 消息头（发送者自身信息）                            │
  │   - 节点 ID、IP、端口、负责的 slot                  │
  │   - slot 位图：16384 bit = 2KB                     │
  │                                                    │
  │ 消息体（携带其他节点的信息）                         │
  │   - 随机选 N 个节点的状态信息（每个约 104 字节）      │
  │   - N = 集群节点总数的 1/10（至少 3 个）             │
  └──────────────────────────────────────────────────┘
```

**问题来了**：节点越多，每条 Gossip 消息越大，网络带宽开销越高。

```
假设 100 个节点：
  消息头：约 2KB（主要是 slot 位图）
  消息体：100 × 1/10 = 10 个节点 × 104B ≈ 1KB
  每条消息 ≈ 3KB
  每个节点每秒发几条 → 100 个节点 → 每秒几百 KB 的 Gossip 流量

假设 1000 个节点：
  消息体：1000 × 1/10 = 100 个节点 × 104B ≈ 10KB
  每条消息 ≈ 12KB
  每个节点每秒发几条 → 1000 个节点 → 每秒几十 MB 的 Gossip 流量
  ↑ 带宽被 Gossip 吃掉了，还挤占处理业务请求的网络带宽
```

所以 Redis 官方建议 **Cluster 最多不超过 1000 个节点**。实际生产中一般控制在几十到几百个节点。

```
Gossip 的代价：
  - 节点少（几十个）：消息小、收敛快，几乎无影响
  - 节点多（几百个）：消息变大、带宽增加，但还能接受
  - 节点太多（上千个）：Gossip 风暴，带宽和 CPU 被大量消耗
```

**为什么哈希槽是 16384 而不是更多？**

这个是 Redis 作者 antirez 在 GitHub issue 中解释过的，核心原因就是 **Gossip 消息体积**。

每条 Gossip 消息的头部要携带一个 **slot 位图（bitmap）**，表示"我负责哪些 slot"——每个 slot 占 1 bit：

```
16384 个 slot → 16384 bit = 2KB   ← 每条 Gossip 消息头 2KB，可以接受
65536 个 slot → 65536 bit = 8KB   ← 每条消息头 8KB，节点一多就爆了
```

```
为什么不更多？
  slot 越多 → 位图越大 → 每条 Gossip 消息越大 → 带宽浪费
  16384 (2KB) 是在消息体积和 slot 粒度之间的平衡点

为什么不更少（比如 1024）？
  slot 太少 → 每个 Master 分到的 slot 太少 → 数据分布不够均匀
  假设 100 个 Master，1024 slot → 每个 Master 只有 10 个 slot
  ↑ 某些 slot 数据特别多时，负载不均，没法精细调整

16384 是刚好合适的数量：
  - 3 个 Master：每个约 5461 个 slot → 分布足够均匀
  - 100 个 Master：每个约 164 个 slot → 迁移粒度足够细
  - Gossip 位图：2KB → 消息体积可控
```

总结一下就是这个公式：

```
slot 数量 = 位图大小
  ↑ 太大 → Gossip 消息膨胀，浪费带宽
  ↑ 太小 → 分片粒度太粗，数据分布不均
  16384 (2KB) → 平衡点
```

### 8.4 三种架构对比

| 特性 | 主从复制 | Sentinel | Cluster |
|------|---------|----------|---------|
| Master 数量 | 1 个 | 1 个 | **多个**（如 3/6/9） |
| 写操作去哪 | Master | Master | key 所在 slot 的 Master |
| 读操作去哪 | Slave | Slave | 默认 Master（可选 Slave） |
| Slave 角色 | 读扩展 + 备份 | 读扩展 + 自动顶替 | 热备 + 自动顶替 |
| 自动故障转移 | 否（手动） | 是 | 是 |
| 数据分片 | 否 | 否 | 是（16384 slot） |
| 写扩展 | 否 | 否 | 是（多 Master 分担） |
| 存储扩展 | 否 | 否 | 是（数据分散在多 Master） |
| 复杂度 | 低 | 中 | 高 |
| 适用场景 | 读多写少，数据量小 | 中小规模高可用 | 大规模、大数据量、高并发 |

---

## 9. 常用场景与最佳实践

### 9.1 分布式 Session

```
用户登录 → 生成 token → SET session:token {user_info} EX 7200
请求带 token → GET session:token → 得到用户信息
```

### 9.2 限流

**固定窗口**：
```bash
INCR rate:user:1:202401011200    # 当前分钟计数
EXPIRE rate:user:1:202401011200 60
# 判断是否超过阈值
```

**滑动窗口**（用 Sorted Set）：
```bash
ZADD rate:user:1 {timestamp} {uuid}        # 记录每次请求
ZREMRANGEBYSCORE rate:user:1 0 {1分钟前}    # 删除窗口外的
ZCARD rate:user:1                           # 当前窗口请求数
```

### 9.3 延迟队列

用 Sorted Set，score 是执行时间戳：

```bash
# 生产者：30 秒后执行
ZADD delay_queue {now + 30} task_data

# 消费者：轮询到期任务
ZRANGEBYSCORE delay_queue 0 {now} LIMIT 0 1
ZREM delay_queue task_data  # 取出后删除
```

### 9.4 排行榜

```bash
ZINCRBY ranking 1 article:123     # 文章阅读量 +1
ZREVRANGE ranking 0 9 WITHSCORES  # Top 10
ZREVRANK ranking article:123      # 某篇文章的排名
```

### 9.5 Key 设计规范

```
业务:对象类型:ID:属性

user:info:1001            # 用户信息
user:session:token-xxx    # 用户 Session
order:detail:2001         # 订单详情
product:stock:3001        # 商品库存
rate:limit:user:1001      # 限流计数
lock:order:2001           # 分布式锁
```

- 用冒号分隔层级
- 避免超长 key（浪费内存）
- 避免 bigkey（String > 10KB，集合 > 5000 元素）
- 设置合理的过期时间，避免内存无限增长

---

## 10. Pipeline 与 Lua 脚本

### 10.1 Pipeline（管道）

正常情况下每条命令一个网络往返（RTT）。Pipeline 把多条命令打包一次发送，减少网络开销。

```
普通模式：                      Pipeline：
  cmd1 → 响应 → cmd2 → 响应      cmd1 + cmd2 + cmd3 → 响应1 + 响应2 + 响应3
  3 次 RTT                       1 次 RTT
```

注意：Pipeline 不是原子的，中间某条命令失败不影响其他命令。

### 10.2 Lua 脚本

Redis 内置 Lua 解释器，脚本在服务端原子执行。

```lua
-- 限流脚本：固定窗口
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = tonumber(redis.call("GET", key) or "0")
if current + 1 > limit then
    return 0  -- 超限
end
redis.call("INCR", key)
if current == 0 then
    redis.call("EXPIRE", key, window)
end
return 1  -- 放行
```

场景：分布式锁释放、限流、库存扣减等需要多步原子操作的地方。

---

## 11. Redis 与 Memcached 对比

| 维度 | Redis | Memcached |
|------|-------|-----------|
| 数据结构 | 丰富（String/Hash/List/Set/ZSet...） | 只有 String |
| 持久化 | 支持（RDB/AOF） | 不支持 |
| 集群 | 原生 Cluster | 客户端一致性哈希 |
| 线程模型 | 单线程（6.0 IO 多线程） | 多线程 |
| 内存效率 | 有额外数据结构开销 | 更高（纯 KV） |
| 适用场景 | 缓存 + 数据结构 + 消息 + 锁 | 纯缓存 |

现在新项目基本都选 Redis，Memcached 只在纯 KV 缓存且对内存效率极致要求时才考虑。
