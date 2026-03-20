# MySQL 理论知识

## 一、整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      客户端                              │
└─────────────────────────┬───────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Server 层                                               │
│  ┌─────────┐ ┌────────┐ ┌──────────┐ ┌────────────────┐│
│  │ 连接器   │→│ 解析器  │→│ 优化器   │→│ 执行器         ││
│  │         │ │        │ │          │ │                ││
│  │认证/连接 │ │词法分析 │ │生成执行  │ │调用存储引擎    ││
│  │管理     │ │语法分析 │ │计划      │ │接口            ││
│  └─────────┘ └────────┘ └──────────┘ └────────────────┘│
└─────────────────────────┬───────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  存储引擎层（插件式）                                     │
│  InnoDB / MyISAM / Memory / ...                         │
└─────────────────────────┬───────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  文件系统：数据文件、日志文件                             │
└─────────────────────────────────────────────────────────┘
```

### 1.1 连接器
- 管理连接、身份认证、权限校验
- 长连接 vs 短连接：长连接复用，但内存占用高，需定期 `mysql_reset_connection`
- 参数：`max_connections`（默认 151）、`wait_timeout`（默认 8h）

### 1.2 解析器
- 词法分析：识别 SELECT、FROM、表名、列名等
- 语法分析：生成抽象语法树（AST）
- 语法错误在这一步报出

### 1.3 优化器
- 基于成本（Cost-Based Optimizer）选择执行计划
- 决定使用哪个索引、JOIN 顺序、是否使用临时表等
- `EXPLAIN` 查看优化器的选择

### 1.4 执行器
- 检查权限
- 调用存储引擎接口逐行执行
- 返回结果集

---

## 二、InnoDB 存储引擎
#### 一切高性能设计的本质，都是为了减少慢速设备（磁盘）的访问，并解决内存碎片和锁竞争。

### 2.1 内存结构

```
InnoDB 内存架构
├── Buffer Pool（缓冲池）
│   ├── 数据页（Data Page）
│   ├── 索引页（Index Page）
│   ├── Change Buffer（写缓冲）
│   ├── 自适应哈希索引（AHI）
│   └── Lock Info（锁信息）
├── Log Buffer（日志缓冲）
├── Adaptive Hash Index
└── Dictionary Cache
```

#### Buffer Pool
- 默认 128MB，生产建议物理内存 60%-80%
- 改良 LRU：分为 young 区（5/8）和 old 区（3/8）
  - 新页先进 old 区头部
  - 在 old 区停留超过 `innodb_old_blocks_time`（1s）后被访问才移入 young 区
  - 解决全表扫描污染热数据问题
- 多实例：`innodb_buffer_pool_instances`，减少并发锁竞争

#### Change Buffer
- 缓存对二级索引的 INSERT/UPDATE/DELETE
- 页不在 Buffer Pool 时，先记录 Change Buffer，后续读取时 merge
- 适用于写多读少场景，减少随机 IO
- `innodb_change_buffer_max_size`（默认 25%）

### 2.2 磁盘结构

```
磁盘文件
├── 系统表空间 ibdata1
│   ├── 数据字典
│   ├── Doublewrite Buffer
│   ├── Change Buffer
│   └── Undo Logs（可分离）
├── 独立表空间 *.ibd（每表一个）
├── Redo Log（ib_logfile0/1）
├── Undo Tablespace
├── 临时表空间 ibtmp1
└── Binlog（Server 层）
```

### 2.3 数据页结构（16KB）

```
┌───────────────────────────┐
│ File Header (38B)         │ 页号、校验和、前后指针、LSN
├───────────────────────────┤
│ Page Header (56B)         │ 记录数、槽数、页类型
├───────────────────────────┤
│ Infimum + Supremum        │ 虚拟最小/最大记录
├───────────────────────────┤
│ User Records              │ 行数据（单链表，按主键有序）
├───────────────────────────┤
│ Free Space                │ 空闲空间
├───────────────────────────┤
│ Page Directory            │ 页目录槽数组，二分查找定位记录
├───────────────────────────┤
│ File Trailer (8B)         │ 校验和（保证完整性）
└───────────────────────────┘
```

### 2.4 行格式

页里面具体的“行记录”是怎么摆放的。这就是 Row Format（行格式）。按引入版本排序：

#### REDUNDANT（MySQL 5.0 之前）
- **特点**：最早的行格式，字段长度偏移列表记录所有列（含定长），NULL 不压缩存储
- **存储开销**：大，每行额外存储信息多
- **场景**：仅用于兼容旧版本，不推荐使用

#### COMPACT（MySQL 5.0，5.0-5.6 默认）
- **特点**：相比 REDUNDANT 大幅优化存储效率
  - 变长字段长度列表：只记录变长列（VARCHAR/VARBINARY/BLOB/TEXT）的实际长度
  - NULL 标志位：用 bitmap 表示 NULL 列，节省空间
  - 大字段溢出：超过 768B 的部分存到溢出页，**页内保留 768B 前缀**
- **场景**：通用场景，旧项目中常见

#### DYNAMIC（MySQL 5.7，5.7+ 默认）
- **特点**：基于 COMPACT 改进，区别在于大字段处理方式
  - 大字段**完全溢出**，页内只存 20B 指针，不保留前缀
  - 数据页能存更多行，减少页分裂
- **场景**：**推荐使用**，适合绝大多数场景，尤其是包含 TEXT/BLOB 的表

#### COMPRESSED（MySQL 5.5）
- **特点**：在 DYNAMIC 基础上增加 zlib 压缩
  - 数据页和溢出页都会压缩
  - 需要设置 `KEY_BLOCK_SIZE`（压缩页大小，通常 8KB）
  - 读取时需要解压，CPU 换 IO
- **场景**：日志表、归档表等读少写少、数据量大的冷数据场景

```sql
-- 查看表的行格式
SHOW TABLE STATUS LIKE 'users'\G

-- 建表时指定
CREATE TABLE t (id INT) ROW_FORMAT=DYNAMIC;

-- 修改已有表
ALTER TABLE t ROW_FORMAT=DYNAMIC;
```

**总结**：新项目直接用 DYNAMIC（MySQL 5.7+ 默认就是），有大字段冷数据可以考虑 COMPRESSED

---

## 三、索引机制

### 3.1 B+Tree 索引

```
                    [非叶子节点：只存 key]
                   /          |          \
        [10,20]           [30,40]         [50,60]
       /   |   \         /   |   \       /   |   \
   [数据] [数据] [数据] [数据] [数据] [数据] [数据] [数据] [数据]
   ←────────────── 叶子节点双向链表 ──────────────→
```

**为什么用 B+Tree**：
1. 非叶子节点不存数据 → 单页能放更多 key → 树更矮 → IO 次数少
2. 叶子节点双向链表 → 范围查询和排序高效
3. 所有查询都到叶子 → 性能稳定

**高度估算**（主键 bigint）：
- 非叶子节点每页约存 1170 个 key
- 三层 B+Tree：1170 × 1170 × 16 ≈ **2000 万行**

### 3.2 聚簇索引 vs 二级索引

| 类型 | 叶子节点存储 | 数量 | 备注 |
|------|-------------|------|------|
| 聚簇索引 | 完整行数据 | 每表 1 个 | 主键 → 唯一非空索引 → 隐藏 row_id |
| 二级索引 | 主键值 | 多个 | 查询需回表 |

**回表**：通过二级索引找到主键 → 再到聚簇索引查完整行

**覆盖索引**：查询列全在索引中 → 无需回表 → EXPLAIN Extra: `Using index`

### 3.3 索引类型

- 主键索引、唯一索引、普通索引、前缀索引
- 联合索引（最左前缀原则）
- 全文索引、空间索引

| 维度 | 字段 | 关键值 (性能从高到低) | 含义与性能细节 | 优化建议 |
| :--- | :--- | :--- | :--- | :--- |
| **访问类型** | **type** | **`system / const`** | 命中主键或唯一索引，单行匹配。$O(1)$ 级别。 | **最优**，无需优化。 |
| | | **`eq_ref`** | 多表关联时，关联字段是主键或唯一索引。 | **优秀**。 |
| | | **`ref`** | 命中非唯一性二级索引。 | **良好**，通常已经很快。 |
| | | **`range`** | 命中索引的范围扫描（如 `BETWEEN`, `>`, `IN`）。 | **一般**，需注意扫描行数（rows）。 |
| | | **`index`** | **全索引扫描**。遍历了整颗索引树，但不扫数据页。 | **预警**，在大表下依然存在大量 I/O。 |
| | | **`ALL`** | **全表扫描**。完全没走索引，直接扫磁盘。 | **必须优化**。 |
| **索引命中** | **key** | `[Index Name]` | 优化器实际选择的索引名称。 | 为 `NULL` 时说明索引失效。 |
| | **key_len** | `[Number]` | 索引使用的字节长度。 | 用于判断联合索引是否触发**最左匹配**。 |
| **辅助信息** | **Extra** | **`Using index`** | **覆盖索引**：不回表，直接从索引页拿数据。 | **性能极佳**。 |
| | | **`Using index condition`** | **索引下推 (ICP)**：在引擎层过滤，减少回表。 | 5.6+ 默认开启。 |
| | | **`Using filesort`** | **文件排序**：无法利用索引排序，在内存/磁盘排序。 | **性能杀手**，需建排序索引。 |
| | | **`Using temporary`** | **临时表**：常出现在 `GROUP BY` 或 `DISTINCT`。 | **高危**，容易导致磁盘 I/O 飙升。 |

### 3.4 最左前缀原则

联合索引 `(a, b, c)`：

| 查询条件 | 是否走索引 |
|---------|-----------|
| `WHERE a = 1` | ✅ |
| `WHERE a = 1 AND b = 2` | ✅ |
| `WHERE a = 1 AND b = 2 AND c = 3` | ✅ |
| `WHERE a = 1 AND c = 3` | ✅ 只用 a |
| `WHERE b = 2` | ❌ 跳过 a |
| `WHERE a = 1 ORDER BY b` | ✅ |
| `WHERE a > 1 AND b = 2` | ✅ 只用 a（范围后失效）|

### 3.5 索引失效场景

```sql
-- 函数/计算
WHERE YEAR(date) = 2024        → WHERE date >= '2024-01-01'
WHERE age + 1 = 20             → WHERE age = 19

-- 隐式类型转换
WHERE varchar_col = 123        → WHERE varchar_col = '123'

-- 左模糊
WHERE name LIKE '%abc'         → 无法优化

-- OR 中有无索引列
WHERE a = 1 OR b = 2           → b 无索引则整体不走

-- != / NOT IN / IS NULL（视情况）
```

### 3.6 索引下推（ICP）

MySQL 5.6+，将 WHERE 条件下推到存储引擎在索引中过滤，减少回表。

```sql
-- 索引 (name, age)
SELECT * FROM users WHERE name LIKE 'A%' AND age = 25;
-- 无 ICP：返回所有 name LIKE 'A%'，Server 层过滤 age
-- 有 ICP：存储引擎同时过滤 name 和 age
```

---

## 四、事务与 MVCC

### 4.1 ACID

| 特性 | 含义 | 实现 | 
|------|------|------|
| 原子性 Atomicity | 要么全做，要么全不做。 | Undo Log |
| 一致性 Consistency | 事务前后数据逻辑一致。 | 其他三者共同保证 |
| 隔离性 Isolation | 多个并发事务之间不能互相干扰。 | MVCC + 锁 |
| 持久性 Durability | 一旦提交，永久保存。 | Redo Log |

### 4.2 隔离级别

| 级别 | 脏读(读到别人没提交的数据) | 不可重复读(两次读同一行结果不同) | 幻读(两次查询结果集行数不同) |
|------|------|-----------|------|
| READ UNCOMMITTED(读未提交) | ✅ (允许)| ✅ | ✅ |
| READ COMMITTED (RC)(读已提交) | ❌（避免） | ✅ | ✅ |
| REPEATABLE READ (RR)(可重复读) | ❌ | ❌ | InnoDB 解决 |
| SERIALIZABLE(串行化) | ❌ | ❌ | ❌ |

**MySQL 默认 RR**，InnoDB 通过 MVCC + Next-Key Lock 解决幻读

### 4.3 MVCC 实现（Multi-Version Concurrency Control，多版本并发控制）
#### 把它想象成一种 “写时复制（Copy-on-Write）” 的变种：写操作不覆盖旧数据，而是通过版本链条把旧数据藏起来，让读操作能看到“过去的幻觉”。
InnoDB 会在每一行数据后面自动加上 3 个隐藏列：
- DB_TRX_ID (6字节)：最近一次修改这行数据的 事务 ID。
- DB_ROLL_PTR (7字节)：回滚指针。它指向 Undo Log 里的上一个版本。
- DB_ROW_ID：如果没有主键，自动生成的隐藏 ID。

**核心组件**：
1. **隐藏列**：`trx_id`（修改事务 ID）、`roll_pointer`（回滚指针指向 Undo Log） （每行记录的“户口本”）
2. **Undo Log 版本链**：同一行的历史版本链 （数据的“时光机”）
3. **ReadView**：事务快照（事务的“快照相机”）

**ReadView 结构**：
- `m_ids`：当前系统里所有还没提交的活跃事务 ID 列表
- `min_trx_id`：m_ids 里最小活跃事务 ID
- `max_trx_id`：下一个待分配事务 ID（即当前最大 ID + 1）。
- `creator_trx_id`：创建 ReadView 的事务 ID

**可见性判断**：
```
trx_id == creator_trx_id          → 可见（自己改的）
trx_id < min_trx_id               → 可见（说明这个事务在你拍照前就提交了，已提交）
trx_id >= max_trx_id              → 不可见（说明这个事务是在你拍照后才开启的，之后开启）
trx_id 在 m_ids 中                 → 不可见（说明你拍照时它还没提交，不能看，未提交）
trx_id 不在 m_ids 中               → 可见（说明它已经提交了，已提交）
```
**如果不可见**：
- 顺着 DB_ROLL_PTR 指针去 Undo Log 找上一个版本，再重复上面的判断，直到找到能看的那个版本为止。


**RC vs RR**：
- RC：每次 SELECT 都会重新生成一个 ReadView。 所以别人提交了，你下一秒查就能看到最新 ReadView。
- RR：只有第一次 SELECT 时生成 ReadView，后面全事务共用。 这样哪怕别人提交了，你的快照没变，看到的还是旧数据。这就是"可重复读"的原理！

### 4.4 InnoDB 如何解决幻读（MVCC + Next-Key Lock）

InnoDB 中有两种读，解决幻读的机制完全不同：

#### 快照读（普通 SELECT）—— 由 MVCC 解决

RR 下第一次 SELECT 生成 ReadView，后续复用。其他事务插入的新行 `trx_id` 不满足可见性规则，"看不到"即"没有幻读"。

```
事务A: BEGIN;
事务A: SELECT * FROM users WHERE age > 20;    → 3 行（生成 ReadView）
事务B: INSERT INTO users (age) VALUES (25);   COMMIT;
事务A: SELECT * FROM users WHERE age > 20;    → 仍然 3 行（复用 ReadView，看不到事务B）
```

#### 当前读（SELECT FOR UPDATE / DML）—— 由 Next-Key Lock 解决

当前读不走 MVCC，读最新数据，必须用锁阻止其他事务插入。

```
假设 age 上有索引，表中现有 age 值：10, 20, 30

执行：SELECT * FROM users WHERE age > 20 FOR UPDATE;

加锁范围（Next-Key Lock）：
  Record Lock: 锁住 age=30 这条记录
  Gap Lock:    锁住 (20, 30) 和 (30, +∞) 的间隙

  被锁定的整体范围：(20, +∞)

此时事务 B 尝试：
  INSERT INTO users (age) VALUES (25);  → 阻塞！落在 (20, 30) 间隙
  INSERT INTO users (age) VALUES (35);  → 阻塞！落在 (30, +∞) 间隙
  UPDATE users SET age=25 WHERE age=10; → 成功，不在锁范围内
```

Gap Lock 锁住索引间隙，阻止其他事务往这些间隙中插入新行，从而防止幻读。

#### 不完美的边界情况：快照读与当前读混用

```
事务A: BEGIN;
事务A: SELECT * FROM users WHERE id = 5;            → 空（快照读，看不到）
事务B: INSERT INTO users (id, name) VALUES (5, 'x'); COMMIT;
事务A: UPDATE users SET name = 'y' WHERE id = 5;    → 影响 1 行（当前读，能看到）
事务A: SELECT * FROM users WHERE id = 5;            → 返回 (5, 'y')（快照被打破）
```

事务 A 快照读没看到 id=5，但 UPDATE（当前读）操作了事务 B 插入的行，之后快照读就能看到了——这就是"幻读"。

**解决**：如果业务要求强一致，第一次读就加锁：`SELECT ... FOR UPDATE`，用 Next-Key Lock 阻止其他事务插入。

#### 总结

| 读类型 | 机制 | 方式 |
|--------|------|------|
| 快照读（SELECT） | MVCC ReadView | "看不到"新行 |
| 当前读（FOR UPDATE / DML） | Next-Key Lock | "阻止"插入新行 |

两者配合在大多数场景下解决幻读，但混用快照读和当前读时仍有边界情况。要求强一致的场景应在第一次读时就加锁。

---

## 五、日志系统

### 5.1 Redo Log（重做日志）

- **所属**：InnoDB 存储引擎特有
- **作用**：崩溃恢复，保证持久性。防止数据库宕机时，Buffer Pool 里的脏页丢失
- **物理日志**：记录页的物理修改
- **WAL**：Write-Ahead Logging，先写日志再写数据，由于磁盘顺序写极快（远快于随机写数据页），MySQL 只要把 Redo Log 写成功，就可以认为事务提交了。
- **循环写入**：固定大小文件组（比如 4GB），写满了就从头覆盖，所以它只负责故障恢复（Crash-safe）。

```
 ib_logfile0              ib_logfile1
┌────────────┐          ┌────────────┐
│ ███░░░░░░░ │   →→→    │ ░░░░░░░░░░ │
└────────────┘          └────────────┘
      ↑ write pos             ↑ checkpoint
```

**刷盘策略** `innodb_flush_log_at_trx_commit`：
- `0`：每秒刷（可能丢 1 秒数据）
- `1`：每次提交刷（**最安全，默认**）
- `2`：写入 OS cache，由系统刷

### 5.2 Undo Log（回滚日志）

- **所属**：InnoDB 存储引擎特有
- **作用**：事务回滚 + MVCC 版本链
- **逻辑日志**：记录反向操作
- INSERT → DELETE，DELETE → INSERT，UPDATE → 旧值 UPDATE

### 5.3 Binlog（归档日志）

- **Server 层**，所有引擎共用，记录逻辑操作
- **作用**：主从复制 + 数据恢复（Point-in-Time Recovery）
- **追加写入**，文件可以有很多个，写满一个开下一个，不会覆盖。

**三种格式**：
- `STATEMENT`：记录 SQL（可能不一致）
- `ROW`：记录行变化（**推荐**）
- `MIXED`：混合

### 5.4 两阶段提交

保证 Redo Log 和 Binlog 一致：
- Redo Log 写了，Binlog 没写 ：数据库重启数据恢复了，但从库没同步到，主从不一致。
- Binlog 写了，Redo Log 没写 ：数据库重启数据丢了，但从库却多出一条，主从不一致。

解决办法：只有两边都“握手”成功，事务才算真正提交。
1. Prepare 阶段：InnoDB 写 Redo Log，标记为 prepare 状态。
2. Commit 阶段：Server 层写 Binlog，写成功后，InnoDB 把 Redo Log 改为 commit 状态。

---

## 六、锁机制

### 6.1 锁分类

```
├── 按粒度
│   ├── 表锁（MyISAM）
│   └── 行锁（InnoDB）
├── 按模式
│   ├── 共享锁 S：SELECT ... LOCK IN SHARE MODE / FOR SHARE
│   └── 排他锁 X：SELECT ... FOR UPDATE / DML
├── 按算法（行锁实现）
│   ├── Record Lock：锁单行
│   ├── Gap Lock：锁间隙，防止插入
│   └── Next-Key Lock：Record + Gap（左开右闭）
└── 意向锁
    ├── IS：意向共享
    └── IX：意向排他
```

### 6.2 加锁规则（RR 级别）

假设表 `users` 主键 `id`（唯一索引），现有数据 `id: 5, 10, 15, 20, 25`，InnoDB 在记录间形成间隙：

```
(-∞, 5]  (5, 10]  (10, 15]  (15, 20]  (20, 25]  (25, +∞)
          ↑ Gap     ↑ Gap     ↑ Gap     ↑ Gap     ↑ Gap
```

Next-Key Lock = 一条记录 + 它前面的间隙（左开右闭），如 `(10, 15]`。

#### 规则 1：加锁基本单位是 Next-Key Lock

所有加锁操作的起点都是 Next-Key Lock，然后根据条件可以"退化"为更小的锁。

#### 规则 2：唯一索引等值查询，命中 → 退化为 Record Lock

```sql
SELECT * FROM users WHERE id = 15 FOR UPDATE;
```

唯一索引保证 `id=15` 只有一条，不可能有其他事务插入同样的 `id=15`，间隙锁没必要。

```
(-∞, 5]  (5, 10]  (10, 15]  (15, 20]  (20, 25]  (25, +∞)
                         🔒
                     只锁 15 这条记录
```

其他事务 INSERT `id=12` 或 `id=17` 不会被阻塞。

#### 规则 3：唯一索引等值查询，未命中 → 退化为 Gap Lock

```sql
SELECT * FROM users WHERE id = 13 FOR UPDATE;
```

`id=13` 不存在，落在 `(10, 15)` 间隙里。记录不存在没有行需要锁，只需锁间隙防止插入。

```
(-∞, 5]  (5, 10]  (10, 15]  (15, 20]  (20, 25]  (25, +∞)
                    🔒🔒🔒
                  锁住 (10, 15) 间隙（开区间，不锁 10 和 15）
```

- `INSERT id=12` → 阻塞（在间隙内）
- `INSERT id=9` → 成功（不在间隙内）
- `UPDATE id=10` → 成功（10 本身没被锁）

#### 规则 4：范围查询锁住扫描到的所有记录和间隙

```sql
SELECT * FROM users WHERE id >= 15 AND id < 22 FOR UPDATE;
```

扫描经过 `id=15`、`id=20`，然后 `id=25` 不满足条件停止。扫描经过的每条记录都加 Next-Key Lock：

```
(-∞, 5]  (5, 10]  (10, 15]  (15, 20]  (20, 25]  (25, +∞)
                   🔒🔒🔒🔒  🔒🔒🔒🔒  🔒🔒🔒🔒
                   (10,15]   (15,20]   (20,25]
```

其中 `id=15` 命中等值且是唯一索引，`(10,15]` 退化为只锁 15。最终实际锁住：`15` + `(15, 25]`

- `INSERT id=18` → 阻塞
- `INSERT id=23` → 阻塞（在 (20,25] 内）
- `INSERT id=26` → 成功
- `INSERT id=12` → 成功

**为什么范围查询要锁这么多？** 范围查询的结果集可能因其他事务的插入而变化（幻读），锁住扫描经过的所有间隙，保证范围内不会有新行插入。

#### 一句话记忆

Next-Key Lock 是默认，能退化就退化——唯一索引等值命中退化为行锁（最小），未命中退化为间隙锁（中等），范围查询不退化（最大）。

### 6.3 死锁

**产生条件**：互斥、占有等待、不可抢占、循环等待

**处理**：
- 等待超时：`innodb_lock_wait_timeout`（默认 50s）
- 死锁检测：`innodb_deadlock_detect`（默认 ON），回滚代价小的事务

**排查**：`SHOW ENGINE INNODB STATUS` → LATEST DETECTED DEADLOCK

**预防**：
- 固定顺序访问表和行
- 缩短事务
- 合理索引，减少锁范围

---

## 七、高可用架构

### 7.1 主从复制

```
Master                            Slave
┌─────────┐                     ┌─────────┐
│ Binlog  │ ─── Dump Thread ──→ │Relay Log│
└─────────┘                     └────┬────┘
                                     │ SQL Thread
                                     ↓
                                ┌─────────┐
                                │ 数据文件 │
                                └─────────┘
```

**三个线程**：
1. Binlog Dump Thread（主）：发送 binlog
2. I/O Thread（从）：接收写入 Relay Log
3. SQL Thread（从）：回放 Relay Log

**复制模式**：
| 模式 | 说明 |
|------|------|
| 异步复制 | 主库不等从库，性能最好但可能丢数据 |
| 半同步复制 | 至少一个从库 ACK 才返回 |
| 组复制 MGR | Paxos 协议，强一致 |

### 7.2 GTID 复制

`GTID = server_uuid:transaction_id`

优点：故障转移自动定位，无需手动指定 binlog 位点

### 7.3 并行复制

| 版本 | 策略 |
|------|------|
| 5.6 | 按库并行 |
| 5.7 | 基于组提交 LOGICAL_CLOCK |
| 8.0 | 基于 writeset |

### 7.4 读写分离

```
App → 代理（ProxySQL/MySQL Router）→ Master(写) / Slave(读)
```

**主从延迟问题**：
- 强制读主
- `WAIT_FOR_EXECUTED_GTID_SET`
- 半同步复制

### 7.5 分库分表

#### 为什么需要

```
单表数据量 → 几千万行后，B+Tree 层级增加，查询变慢
单库连接数 → max_connections 有限，扛不住高并发
单机磁盘   → 数据量到 TB 级，IO 瓶颈
单机 CPU   → 所有读写打到一台机器
```

主从复制只解决读的压力，解决不了写和数据量的瓶颈。

#### 什么时候该分

```
1. 单表超过 2000 万行了吗？        → 没到就别分
2. 加索引、SQL 优化做了吗？         → 没做就别分
3. 读写分离做了吗？                 → 没做就别分
4. 缓存（Redis）加了吗？            → 没加就别分
5. 以上都做了还扛不住？              → 可以考虑分了
```

分库分表是最后手段，一旦拆了开发复杂度大幅增加。

#### 垂直拆分 —— 按业务切

**垂直分库**：按业务域独立成库，本质是微服务在数据库层的体现

```
拆分前                          拆分后
┌──────────────┐     ┌──────────┐ ┌──────────┐ ┌──────────┐
│   all_db     │     │ user_db  │ │ order_db │ │ pay_db   │
│ users        │ →   │ users    │ │ orders   │ │ payments │
│ orders       │     │ addresses│ │ items    │ │ refunds  │
│ payments     │     └──────────┘ └──────────┘ └──────────┘
└──────────────┘
```

**垂直分表**：一张宽表拆成多张窄表，冷热分离

```
拆分前                              拆分后
┌─────────────────────────┐     ┌───────────────┐ ┌──────────────┐
│ id | name | age | bio   │ →   │ id | name|age │ │ uid | bio    │
│         (bio 是大 TEXT)  │     │   热数据，快   │ │  冷数据，按需 │
└─────────────────────────┘     └───────────────┘ └──────────────┘
```

#### 水平拆分 —— 按数据行切

单表几千万行扛不住，把行数据分散到多张表或多个库。

```
水平分表：拆成 4 张表，还在同一个库
┌──────────────┐ ┌──────────────┐
│ orders_0     │ │ orders_1     │
│ id % 4 == 0  │ │ id % 4 == 1  │
└──────────────┘ └──────────────┘
┌──────────────┐ ┌──────────────┐
│ orders_2     │ │ orders_3     │
│ id % 4 == 2  │ │ id % 4 == 3  │
└──────────────┘ └──────────────┘

水平分库：分散到不同 MySQL 实例，连接数/磁盘/CPU 都分散
┌──────────────┐ ┌──────────────┐
│ db_0 (机器A)  │ │ db_1 (机器B)  │
│ orders_0     │ │ orders_1     │
│ orders_2     │ │ orders_3     │
└──────────────┘ └──────────────┘
```

#### 分片键（Sharding Key）

决定一行数据去哪个分片的字段。以电商订单为例：

| 候选分片键 | 效果 |
|-----------|------|
| `user_id`（推荐）| 同一用户的订单在同一分片，查"我的订单"不跨库 |
| `order_id` | 分布均匀，但查某用户所有订单要查所有分片 |
| `create_time` | 按时间 range 分，新数据集中在最新分片，热点严重 |

好的分片键 = 查询最频繁的条件字段，让大多数查询只落在一个分片。

#### 路由策略

| 策略 | 原理 | 优点 | 缺点 |
|------|------|------|------|
| Hash 取模 | `hash(user_id) % 分片数` | 数据均匀 | 扩容难，几乎所有数据要重分配 |
| Range 范围 | `id 1~1000万→db_0, 1001万~2000万→db_1` | 扩容简单 | 热点问题，新数据集中在最新分片 |
| 一致性哈希 | 数据和节点映射到哈希环 | 扩容只迁移一小部分 | 实现复杂 |

#### 带来的五大问题

**1. 跨库 JOIN**

```sql
-- 拆分前一条 SQL
SELECT o.*, u.name FROM orders o JOIN users u ON o.user_id = u.id;
-- 拆分后 orders 和 users 在不同库，无法 JOIN
```

解决：应用层两次查询聚合、冗余字段、同步到 ES

**2. 分布式事务**

```
下单 = 扣库存(product_db) + 建订单(order_db) + 扣余额(pay_db)
三个库，本地事务管不了
```

| 方案 | 特点 |
|------|------|
| XA 两阶段提交 | 性能差，很少用 |
| TCC（Try-Confirm-Cancel）| 业务侵入大但可靠 |
| Saga | 长事务拆成多个本地事务 + 补偿 |
| 最终一致性（消息队列）| 最常用 |

**3. 全局唯一 ID**

拆分后不能用 AUTO_INCREMENT，两个库都从 1 开始会冲突。

| 方案 | 特点 |
|------|------|
| Snowflake 雪花算法 | 64 位 = 时间戳+机器ID+序列号，最主流 |
| 数据库号段 | 从中心库批量取 ID 段（如 1-1000） |
| UUID | 简单但太长无序，不适合做主键 |

**4. 跨分片分页**

```sql
SELECT * FROM orders ORDER BY create_time LIMIT 90, 10;
-- 4 个分片各取 Top 100 → 应用层归并排序 → 取第 91-100 条
-- 页码越深，每个分片取的数据越多，性能越差
```

解决：游标分页 `WHERE create_time < 上一页最后时间 LIMIT 10`，禁止跳页

**5. 扩容迁移**

4 分片扩到 8 分片的双写方案：

```
1. 新建 8 个分片
2. 开启双写：新数据同时写新旧两套
3. 后台迁移历史数据到新分片
4. 校验一致后切读到新分片
5. 停掉旧分片写入
```

#### 常见中间件

| 中间件 | 类型 | 说明 |
|--------|------|------|
| ShardingSphere | 代理 + SDK | 功能最全，Java 生态主流 |
| Vitess | 代理 | YouTube 出品，Go 友好 |
| ProxySQL | 代理 | 轻量，主要做读写分离 |

---

## 八、性能调优

### 8.1 EXPLAIN 关键字段

| 字段 | 重点 |
|------|------|
| type | `const > eq_ref > ref > range > index > ALL` |
| key | NULL 表示全表扫描 |
| rows | 预估扫描行数 |
| Extra | `Using index`（好）、`Using filesort`/`Using temporary`（需优化）|

### 8.2 索引设计原则

1. 选择性高的列优先
2. 联合索引：等值 → 范围 → 排序
3. 覆盖索引减少回表
4. 前缀索引节省空间
5. 避免冗余索引
6. 单表索引不超过 5-6 个

### 8.3 SQL 优化

```sql
-- 深分页优化
-- ❌ LIMIT 1000000, 10
-- ✅ WHERE id > last_id LIMIT 10

-- COUNT 优化
-- 维护计数表 / Redis 缓存 / 估算

-- JOIN 优化
-- 小表驱动大表，被驱动表建索引，不超过 3 表

-- ORDER BY 优化
-- 建立 (where_col, order_col) 联合索引
```

### 8.4 关键参数

```ini
# Buffer Pool
innodb_buffer_pool_size = 物理内存 60%-80%
innodb_buffer_pool_instances = 8

# 日志
innodb_log_file_size = 1G
innodb_flush_log_at_trx_commit = 1
sync_binlog = 1

# 连接
max_connections = 500
wait_timeout = 28800
```

### 8.5 慢查询分析

```sql
SET GLOBAL slow_query_log = ON;
SET GLOBAL long_query_time = 1;
```

工具：`mysqldumpslow`、`pt-query-digest`

---

## 九、InnoDB vs MyISAM

| 特性 | InnoDB | MyISAM |
|------|--------|--------|
| 事务 | ✅ | ❌ |
| 行锁 | ✅ | 表锁 |
| 外键 | ✅ | ❌ |
| 崩溃恢复 | ✅ | ❌ |
| MVCC | ✅ | ❌ |
| 聚簇索引 | ✅ | ❌ |
| COUNT(*) | 全表扫描 | O(1) |

**结论**：几乎所有场景都应使用 InnoDB

---

## 十、备份与恢复

### 10.1 备份工具

| 工具 | 类型 | 特点 |
|------|------|------|
| mysqldump | 逻辑 | 简单，小数据量 |
| mysqlpump | 逻辑 | 多线程 |
| Xtrabackup | 物理 | 热备，不锁表，**生产首选** |

### 10.2 恢复策略

- 全量恢复：还原最近备份
- 增量恢复：全量 + binlog 回放（PITR）

```bash
# Xtrabackup
xtrabackup --backup --target-dir=/backup/full

# binlog 恢复
mysqlbinlog --start-datetime="2024-01-01 00:00:00" binlog.000001 | mysql -uroot -p
```
