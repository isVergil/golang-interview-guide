# PostgreSQL 面试高频题

> 以下回答按照面试中口语化表达整理，重点在于思路清晰、逻辑连贯，面试官听得舒服。
> 部分题目会和 MySQL 做对比，帮助建立两者的知识映射。

---

## 一、架构与原理

### 1. PostgreSQL 的整体架构是什么样的？和 MySQL 有什么区别？

PG 是**多进程**架构，MySQL 是多线程架构，这是最本质的区别。

PG 启动后有一个 Postmaster 主进程，每来一个客户端连接就 fork 一个子进程（Backend Process）来处理。所有后端进程通过**共享内存**通信，共享内存里主要有三块东西：

- **Shared Buffer Pool**：数据页缓存，类似 MySQL 的 InnoDB Buffer Pool
- **WAL Buffer**：预写日志缓冲区，类似 MySQL 的 Redo Log Buffer
- **CLOG**：事务提交状态日志，记录每个事务是提交了还是回滚了

多进程 vs 多线程的影响：

- **隔离性好**：一个进程崩了不影响其他连接，MySQL 一个线程崩了可能拖垮整个服务
- **内存开销大**：每个进程独占一份内存，默认 100 个连接就不少了
- **连接数受限**：生产环境**必须**配连接池（PgBouncer / Pgpool-II），否则几百个连接就把服务器吃满了

所以你会发现 PG 的部署方案里连接池是标配，而 MySQL 连接池是可选的。

### 2. PostgreSQL 的 Shared Buffer Pool 和 MySQL 的 Buffer Pool 有什么区别？

功能类似，都是数据页的内存缓存，减少磁盘 IO。但有几个不同点：

**页大小不同**：PG 默认 8KB，MySQL InnoDB 默认 16KB。

**LRU 策略不同**：MySQL 的 Buffer Pool 用改进的 LRU（young/old 分区），PG 用的是 **clock-sweep** 算法（时钟扫描）。clock-sweep 比 LRU 开销更小，但在缓存淘汰策略上各有优劣。

**脏页刷盘机制不同**：PG 有 Background Writer 和 Checkpointer 两个后台进程负责刷脏页。MySQL 是 InnoDB 的后台线程负责。

配置建议：`shared_buffers` 一般设为系统内存的 25%。MySQL 的 `innodb_buffer_pool_size` 一般设 50%~75%，因为 PG 还依赖操作系统的页缓存（double buffering），不能设太大。

### 3. WAL 是什么？和 MySQL 的 Redo Log 有什么区别？

WAL（Write-Ahead Logging）就是预写日志，和 MySQL 的 Redo Log 本质一样：先写日志再写数据，保证崩溃后能恢复。

写数据的流程是这样的：
1. 修改 Shared Buffer 中的数据页（变成脏页）
2. 写 WAL 日志到 WAL Buffer
3. 事务提交时 WAL Buffer fsync 到磁盘
4. 脏页由后台进程异步刷盘

和 MySQL Redo Log 的区别：

**大小不同**：Redo Log 是固定大小循环写（写满了覆盖），WAL 是追加写不覆盖，旧的 WAL 文件可以归档用于复制和备份恢复。

**用途更广**：MySQL 的 Redo Log 只用于崩溃恢复。PG 的 WAL 同时用于**崩溃恢复**、**流复制**和 **PITR（基于时间点恢复）**。相当于 PG 把 MySQL 的 Redo Log + Binlog 的功能合二为一了。

**没有两阶段提交**：MySQL 需要 Redo Log 和 Binlog 之间做两阶段提交保证一致性。PG 只有 WAL 一套日志系统，天然一致，不需要两阶段提交，架构更简洁。

### 4. PG 有 Binlog 吗？主从复制用什么？

PG 没有 Binlog 这个概念。MySQL 的 Binlog 干两件事：主从复制和数据恢复。PG 用 **WAL** 一套机制搞定了这两件事。

**流复制（Streaming Replication）**就是基于 WAL 日志流做的。主库把 WAL 实时发给从库，从库回放 WAL 来同步数据。这是**物理复制**（直接复制数据页的变更），比 MySQL 的 Binlog **逻辑复制**（复制 SQL 语句级别的变更）延迟更低。

PG 10 之后也支持**逻辑复制**（类似 MySQL Binlog 的 row-based 模式），可以选择性复制部分表，支持跨版本。

---

## 二、存储结构

### 5. PG 的堆表和 MySQL 的聚簇索引表有什么区别？

这是 PG 和 MySQL 存储层面**最核心**的区别。

**MySQL InnoDB 是聚簇索引表**：数据按主键顺序存在 B+Tree 的叶子节点里，主键索引和数据是一体的。二级索引叶子节点存主键值，查到主键后要**回表**去主键索引取完整数据。

**PG 是堆表（Heap Table）**：数据按插入顺序无序地堆在文件里。所有索引（包括主键索引）叶子节点存的都是 **ctid**（行指针，格式是 `(page, offset)`），指向堆表中的物理位置。

对比影响：

| 场景 | PG 堆表 | MySQL 聚簇索引 |
|------|---------|---------------|
| 主键查询 | 索引 → ctid → 回堆表（多一次 IO） | 直接在叶子节点取数据 |
| 写入性能 | 堆表无序追加，写入快 | 主键顺序写快，乱序写有页分裂 |
| UPDATE | 写新版本到堆表（旧版本留在原地） | 原地更新（旧值写 undo log） |
| 全表扫描 | 顺序扫描堆表，连续 IO 效率高 | 扫主键索引叶子节点 |
| 表膨胀 | 有（旧版本残留在堆表） | 无（undo log 单独管理） |

面试官追问的话要提到：PG 的 UPDATE 是"删旧写新"（旧行标记 t_xmax，新行写入新位置），这导致了**表膨胀**问题，需要 VACUUM 清理。

### 6. ctid 是什么？

ctid 是 PG 中每行数据的**物理地址**，格式 `(page, offset)`，比如 `(0, 5)` 表示第 0 页第 5 个 item。

它是 PG 索引的核心——所有索引叶子节点存的都是 ctid，通过 ctid 跳到堆表取数据。

```sql
-- 可以直接查看
SELECT ctid, id, name FROM users;
-- (0,1) | 1 | alice
-- (0,2) | 2 | bob
```

注意 ctid **不是稳定的**。UPDATE 会生成新行，ctid 会变。VACUUM 整理后 ctid 也可能变。所以不能把 ctid 当作业务 ID 用。

### 7. PG 的 Tuple 结构了解吗？t_xmin 和 t_xmax 是什么？

PG 里每行数据叫 Tuple，由行头（HeapTupleHeader）+ 数据组成。行头里最重要的是这几个字段：

- **t_xmin**：插入这行的事务 ID，行可见的起点
- **t_xmax**：删除/更新这行的事务 ID，行不可见的起点（0 表示未被删除）
- **t_ctid**：当前行的 ctid，UPDATE 后指向新版本（形成版本链）

这三个字段是 MVCC 的核心：

```
INSERT：t_xmin = 当前事务ID, t_xmax = 0
DELETE：把 t_xmax 设为当前事务ID（行并不物理删除）
UPDATE：旧行 t_xmax = 当前事务ID，新行 t_xmin = 当前事务ID
         旧行 t_ctid 指向新行（形成版本链）
```

和 MySQL 的区别：MySQL 用 `trx_id` + `roll_pointer` 指向 undo log 中的旧版本。PG 的旧版本直接留在堆表里，通过 t_xmin/t_xmax 判断可见性。

### 8. TOAST 是什么？

TOAST（The Oversized-Attribute Storage Technique）是 PG 处理大字段的机制。当一行数据超过页大小的 1/4（约 2KB）时自动触发。

四种策略：
- **PLAIN**：不压缩不外存（定长类型如 int）
- **EXTENDED**：先压缩，还放不下就外存到单独的 TOAST 表（默认，text/jsonb 用这个）
- **EXTERNAL**：不压缩直接外存（需要频繁取子串时用）
- **MAIN**：先压缩，尽量不外存

简单理解：大字段自动拆到单独的 TOAST 表存储，主表只存一个指针。对应用层透明，不需要额外处理。

---

## 三、数据类型

### 9. PG 相比 MySQL 多了哪些特色数据类型？

这是 PG 的一大优势，很多 MySQL 需要应用层实现的功能，PG 在数据库层原生支持：

**JSONB**：二进制 JSON，支持 GIN 索引，查询效率高。MySQL 也有 JSON 类型但没有二进制格式，查询性能差很多。

**数组（Array）**：原生数组类型 `text[]`、`int[]`，支持 `ANY`、`@>`、`&&` 等操作符，配合 GIN 索引。MySQL 没有。

**范围类型（Range）**：`tstzrange`（时间范围）、`int4range`（整数范围）等。天然适合时间段、预订系统。配合排除约束（EXCLUDE）可以在数据库层直接防止重叠。MySQL 没有。

**tsvector/tsquery**：全文检索类型，不依赖 ES 就能做搜索。MySQL 也有 FULLTEXT 但功能弱很多。

**uuid**：原生 UUID 类型，`gen_random_uuid()` 直接生成。

**inet/cidr**：IP 地址和网段类型，做白名单、网络管理很方便。

**vector**（pgvector 扩展）：向量类型，AI embedding 检索，这是 PG 在 AI 时代的杀手级特性。

### 10. JSONB 和 JSON 有什么区别？JSONB 查询怎么加速？

**json** 是文本原样存储，每次查询都要重新解析，不支持索引。**jsonb** 是二进制格式预解析存储，查询快，支持 GIN 索引。**生产环境统一用 jsonb**。

JSONB 常用操作符：
- `->` 取 JSON 值（返回 jsonb，可链式取值）
- `->>` 取文本值（返回 text，最后一层用这个）
- `@>` 包含查询（走 GIN 索引，**最常用**）
- `?` key 存在判断
- `||` 合并、`-` 删除字段、`#-` 删除嵌套字段

加速方式：

```sql
-- 通用 GIN 索引（加速 @>、?、?| 等操作）
CREATE INDEX idx_data ON events USING GIN (data);
SELECT * FROM events WHERE data @> '{"type":"click"}';

-- 表达式索引（只索引特定字段，更小更快）
CREATE INDEX idx_type ON events ((data->>'type'));
SELECT * FROM events WHERE data->>'type' = 'click';
```

面试加分点：JSONB 的 `@>` 包含查询能走 GIN 索引，但 `->>` 取值后做 `=` 比较走不了 GIN，需要单独建表达式索引。

### 11. 范围类型有什么用？相比 MySQL 的两字段方案有什么优势？

范围类型把"一段区间"作为一等公民。典型场景是会议室预订：

```sql
CREATE TABLE bookings (
    id serial PRIMARY KEY,
    room text,
    during tstzrange,
    -- 排除约束：同房间时间不能重叠
    EXCLUDE USING GIST (room WITH =, during WITH &&)
);
```

插入重叠时间段数据库自动报错，不需要应用层加锁判断。

和 MySQL 用 `start_time / end_time` 两字段相比：

- **重叠检测**：PG 一个 `&&` 运算符搞定，MySQL 要写 `start < end2 AND end > start2`，容易写错
- **数据库层约束**：PG 用 EXCLUDE 约束自动拒绝冲突，MySQL 只能应用层加锁 → 查冲突 → 插入
- **开闭区间**：PG 原生支持 `[)` `(]` 语义精确，MySQL 靠 `<` / `<=` 人工控制
- **索引支持**：GiST 索引原生支持范围运算符，MySQL B+Tree 对两列时间范围查询效率低
- **交集/并集**：PG 原生 `*`（交集）、`+`（并集）运算符，MySQL 要应用层计算

### 12. PG 的数组类型和 unnest 是什么？

PG 原生支持数组类型，不需要拆关联表就能存多值：

```sql
CREATE TABLE articles (id serial, title text, tags text[]);
INSERT INTO articles VALUES (1, 'Go 入门', ARRAY['go','tutorial']);

-- 查询
SELECT * FROM articles WHERE 'go' = ANY(tags);         -- 包含某元素
SELECT * FROM articles WHERE tags @> ARRAY['go','redis']; -- 包含所有
```

配合 GIN 索引加速数组查询。

**unnest** 是数组最重要的配套函数，把数组展开成行：

```sql
SELECT unnest(ARRAY['go','redis','pgsql']);
-- go / redis / pgsql（三行）
```

**实战用法：批量写入**。传统批量 INSERT 要拼 VALUES，PG 用 unnest 更优雅：

```sql
INSERT INTO users (name, email)
SELECT * FROM unnest(ARRAY['alice','bob'], ARRAY['a@t.com','b@t.com']);
```

Go 里配合 pgx，直接传切片即可，pgx 自动映射为 PG 数组。比拼 SQL 更安全（防注入）、性能更好（一次网络往返）、无参数上限。

---

## 四、索引

### 13. PG 有哪些索引类型？分别适合什么场景？

PG 支持 **6 种**索引，远多于 MySQL（基本只有 B+Tree）：

| 索引 | 适合查询 | 典型场景 |
|------|---------|---------|
| **B-Tree** | =, <, >, BETWEEN, ORDER BY | 大多数场景（默认） |
| **Hash** | 纯 = | 超大表精确匹配（PG 10 之后安全可用） |
| **GIN** | @>, ?, @@, &&（多值包含） | JSONB、数组、全文检索 |
| **GiST** | &&, @>, <->（范围/空间） | 范围类型、PostGIS、排除约束 |
| **BRIN** | 范围查询（数据物理有序） | 时序数据、日志表 |
| **SP-GiST** | 空间分区 | 非均匀空间数据（较少用） |

面试重点讲 B-Tree、GIN、BRIN 三个就行，其他知道就好。

### 14. GIN 索引是什么？原理是什么？

GIN（Generalized Inverted Index）是**倒排索引**，和搜索引擎的原理一样。

原理：对每个值维护一个 posting list（包含该值的行 ID 列表）。

```
文档1: tags = ['redis','cache']
文档2: tags = ['redis','lock']
文档3: tags = ['mysql','index']

GIN 索引：
  'redis' → [文档1, 文档2]
  'cache' → [文档1]
  'lock'  → [文档2]
  'mysql' → [文档3]
```

适合**多值列**的查询：JSONB 的 `@>` 包含、数组的 `@>` `&&`、全文检索的 `@@`。

缺点是写入慢（每次写入要更新多个 posting list），所以 GIN 有个 fastupdate 机制：先把变更暂存到 pending list，积累一批后再合并到主索引，空间换时间。

### 15. BRIN 索引是什么？什么时候用？

BRIN（Block Range Index）是**块范围索引**。它不索引每一行，而是记录每个数据页范围（默认 128 页为一组）的**最小值和最大值**。

适合场景：**数据物理有序的大表**，最典型的就是按时间顺序插入的日志表、时序数据。

```sql
CREATE INDEX idx_logs_time ON logs USING BRIN (created_at);
```

BRIN 的体积极小。同样 1 亿行的时序表，B-Tree 索引约 2GB，BRIN 索引约 100KB，差 2 万倍。

但前提是数据物理有序。如果数据乱序插入，BRIN 记录的 min/max 范围会很大，扫描大量无效页，就退化了。

### 16. 部分索引、覆盖索引、表达式索引分别是什么？

**部分索引（Partial Index）**：只索引满足条件的行，体积更小：

```sql
-- 只索引未删除的用户（软删除场景超好用）
CREATE INDEX idx_active ON users (email) WHERE deleted_at IS NULL;
```

MySQL 不支持部分索引，这是 PG 的独有优势。

**覆盖索引（INCLUDE）**：索引额外携带列，避免回堆表：

```sql
CREATE INDEX idx_email ON users (email) INCLUDE (name);
-- SELECT name FROM users WHERE email = 'x' → Index Only Scan，不回表
```

MySQL 的覆盖索引是把查询列都放进联合索引，PG 用 INCLUDE 更灵活——INCLUDE 的列不参与索引排序，只是"搭便车"。

**表达式索引**：索引一个表达式的结果：

```sql
CREATE INDEX idx_lower ON users (lower(email));
-- WHERE lower(email) = 'alice@test.com' 走索引
```

### 17. PG 的索引和 MySQL 的索引有什么本质区别？

最核心的区别：**PG 索引叶子节点存 ctid（行指针），MySQL 聚簇索引叶子节点直接存数据**。

这意味着：

- PG **所有索引**（包括主键索引）查到的都是 ctid，都需要回堆表取数据。除非用覆盖索引（Index Only Scan）
- MySQL 主键查询直接在索引叶子节点取数据，不回表。二级索引取到主键后回表查主键索引

PG 主键查询比 MySQL 多一次 IO（多了回堆表），但 PG 的全表扫描效率更高（顺序扫描堆表是连续 IO）。

另外 PG 支持 6 种索引类型，MySQL 基本只有 B+Tree，这是 PG 在复杂查询场景下的重要优势。

---

## 五、MVCC

### 18. PG 的 MVCC 和 MySQL 的 MVCC 有什么区别？

这是面试高频题。两者都是 MVCC，但实现方式**完全不同**：

**MySQL 的 MVCC**：旧版本存在 undo log 里，堆表只有最新版本。读旧版本时通过 roll_pointer 沿着 undo log 版本链构造。

**PG 的 MVCC**：旧版本**直接留在堆表里**。每行通过 t_xmin/t_xmax 标记生命周期，读的时候根据事务快照判断哪个版本可见。

| 对比 | PostgreSQL | MySQL InnoDB |
|------|-----------|-------------|
| 旧版本存放 | 堆表内（原地） | undo log（单独区域） |
| UPDATE | 标记旧行删除 + 写入新行 | 原地更新 + 旧值写 undo log |
| 清理机制 | VACUUM | purge 线程自动 |
| 表膨胀 | 有 | 无 |
| 读旧版本 | 直接读堆表中的旧行（快） | 沿 undo log 链构造（慢） |

PG 的优势：读旧版本不需要构造，直接读堆表，快。
PG 的劣势：旧版本占堆表空间，表会膨胀，必须靠 VACUUM 清理。

### 19. PG 的可见性判断规则是什么？

一行数据对某个事务是否可见，由 t_xmin、t_xmax 和当前事务的**快照**决定。

简化规则：
1. t_xmin 的事务**已提交** 且在我的快照之前 → 行"出生了"，看得见
2. t_xmax = 0（未被删除） → 行还"活着"
3. t_xmax 的事务**已提交** 且在我的快照之前 → 行已"死亡"，看不见

综合：**可见 = t_xmin 已提交且在快照前 && (t_xmax=0 || t_xmax 未提交 || t_xmax 在快照后)**

和 MySQL 对比：MySQL 是通过 ReadView + undo log 版本链判断可见性，PG 是直接看行头的 t_xmin/t_xmax。原理相似但实现路径不同。

### 20. VACUUM 是什么？为什么 PG 必须有 VACUUM？

VACUUM 是 PG **独有的**维护操作，MySQL 没有。因为 PG 的 MVCC 把旧版本留在堆表里，DELETE/UPDATE 产生的**死元组（dead tuples）**不会自动回收。

不做 VACUUM 的后果：
- 表越来越大（膨胀），查询越来越慢
- 索引也跟着膨胀
- 最严重的：事务 ID 耗尽后数据库**强制关机**拒绝写入

三种 VACUUM：

| 类型 | 作用 | 是否锁表 |
|------|------|---------|
| VACUUM | 标记死元组空间可复用 | 不锁表 |
| VACUUM FULL | 重写整表，物理回收空间 | **锁表** |
| VACUUM FREEZE | 冻结旧事务 ID | 不锁表 |

生产环境靠 **autovacuum** 自动执行，默认开启。触发条件：

```
dead_tuples > 50 + 0.2 × 总行数（默认）
```

面试加分点：提到 autovacuum 必须保证正常运行。如果被关掉或卡住，事务 ID 回卷会导致数据库强制停机，这是 PG 运维中最危险的问题之一。

### 21. 什么是 HOT 更新？

HOT（Heap Only Tuple）是 PG 对 UPDATE 的优化。普通 UPDATE 要在堆表写新行 + 更新所有索引指向新行，开销大。

如果同时满足两个条件：
1. 更新的列**不在任何索引**中
2. 新行能放在**同一页**内

就可以做 HOT 更新：旧行 ctid 指向新行，索引不用动。索引仍指向旧行，通过 ctid 链跳转到新行。

表设计建议：高频更新的列尽量不要建索引，提高 HOT 命中率。

---

## 六、事务与锁

### 22. PG 的事务隔离级别和 MySQL 有什么不同？

PG 支持四个 SQL 标准隔离级别，但实现上只有三种（Read Uncommitted 行为等同于 Read Committed）：

| 隔离级别 | PG 行为 |
|---------|--------|
| Read Uncommitted | 等同于 Read Committed（PG 不允许脏读） |
| **Read Committed** | **默认**。每条 SQL 取最新快照 |
| Repeatable Read | 事务开始时取快照，整个事务复用 |
| Serializable | SSI（可序列化快照隔离） |

和 MySQL 的关键区别：

**默认级别不同**：PG 默认 Read Committed，MySQL 默认 Repeatable Read。

**幻读处理不同**：PG 的 Repeatable Read 是**真正的快照隔离**，直接没有幻读。MySQL 的 RR 通过 Gap Lock / Next-Key Lock 防幻读，但有边界情况（先快照读再当前读可能看到幻行）。

**Serializable 不同**：PG 用 SSI（Serializable Snapshot Isolation），基于快照检测冲突，不用加大量锁。MySQL 的 Serializable 是所有 SELECT 都加 S 锁，性能很差。

### 23. PG 支持 DDL 事务吗？

**支持**。这是 PG 相比 MySQL 的一个重要优势。

```sql
BEGIN;
CREATE TABLE test (id int);
INSERT INTO test VALUES (1);
ROLLBACK;  -- 整个事务回滚，test 表也不会创建
```

MySQL 的 DDL（CREATE TABLE、ALTER TABLE）会自动提交，无法回滚。这意味着 MySQL 做数据库迁移的时候，如果中间失败了，已经执行的 DDL 无法撤销。

PG 的 DDL 事务让数据库迁移更安全：一组 DDL 要么全成功要么全回滚。

### 24. PG 有哪些锁？和 MySQL 的锁有什么区别？

**表级锁**：PG 有 8 种表级锁模式，最常遇到的：
- ACCESS SHARE：普通 SELECT，只和 ACCESS EXCLUSIVE 冲突
- ROW EXCLUSIVE：INSERT/UPDATE/DELETE
- ACCESS EXCLUSIVE：ALTER TABLE、DROP、VACUUM FULL，和所有锁冲突

> 生产踩坑点：`ALTER TABLE` 需要 ACCESS EXCLUSIVE 锁，会锁全表。大表 DDL 变更必须小心。

**行级锁**：

```sql
SELECT * FROM users WHERE id = 1 FOR UPDATE;        -- 排他锁
SELECT * FROM users WHERE id = 1 FOR SHARE;          -- 共享锁
SELECT * FROM users WHERE id = 1 FOR UPDATE NOWAIT;   -- 加锁失败立即报错
SELECT * FROM tasks WHERE status = 'pending' 
  LIMIT 1 FOR UPDATE SKIP LOCKED;                    -- 跳过已锁行（队列场景）
```

`FOR UPDATE SKIP LOCKED` 是 PG 的特色，非常适合任务队列场景：多个 worker 竞争取任务，自动跳过被锁的行。

**和 MySQL 锁的区别**：
- MySQL 有 Gap Lock、Next-Key Lock 防幻读，PG 没有（PG 用快照隔离解决幻读）
- PG 有 Advisory Lock（咨询锁），MySQL 没有原生支持
- PG 行锁不依赖索引（MySQL 行锁必须走索引，否则升级为表锁）

### 25. Advisory Lock（咨询锁）是什么？

Advisory Lock 是 PG 提供的**应用层锁**。不锁表不锁行，锁的是一个自定义的 bigint key，业务自己决定 key 的含义。

三种类型：

```sql
-- 会话级（手动释放）
SELECT pg_advisory_lock(12345);
SELECT pg_advisory_unlock(12345);

-- 事务级（事务结束自动释放）
SELECT pg_advisory_xact_lock(12345);

-- 非阻塞（返回 bool）
SELECT pg_try_advisory_lock(12345);
```

典型场景：**分布式锁**。类似 Redis 的 SETNX，但不依赖额外组件。

和 Redis 分布式锁对比：
- 不需要额外部署 Redis
- 和 PG 事务原子（事务级锁事务结束自动释放）
- 没有 TTL 过期风险（Redis 锁如果 TTL 设短了，业务没执行完锁就过期了）
- 缺点是只能在同一个 PG 实例上竞争

### 26. PG 怎么处理死锁？

PG 有**死锁检测器**，默认每秒检测一次（`deadlock_timeout = 1s`）。发现死锁后自动回滚其中一个事务。

```sql
-- 事务A                          -- 事务B
UPDATE users SET age=1 WHERE id=1; UPDATE users SET age=2 WHERE id=2;
UPDATE users SET age=1 WHERE id=2; UPDATE users SET age=2 WHERE id=1;
-- 死锁！PG 回滚代价较小的事务
```

预防方式和 MySQL 一样：所有事务按相同顺序访问资源（比如按 id 升序加锁）。

---

## 七、查询优化

### 27. EXPLAIN 怎么看？和 MySQL 的 EXPLAIN 有什么区别？

```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT * FROM users WHERE age > 25;
```

输出关键字段：

| 字段 | 含义 |
|------|------|
| cost=0.00..1.05 | 启动代价..总代价 |
| rows=3 | 估算返回行数 |
| actual time | 实际执行时间（ANALYZE 才有） |
| Buffers: shared hit | 缓存命中页数 |
| Buffers: shared read | 磁盘读取页数 |

**常见扫描方式**（从好到差）：

| 扫描方式 | 含义 |
|---------|------|
| Index Only Scan | 覆盖索引，不回表（最快） |
| Index Scan | 索引扫描 → 回堆表 |
| Bitmap Index Scan + Bitmap Heap Scan | 索引构建位图 → 批量回表 |
| Seq Scan | 全表扫描 |

和 MySQL EXPLAIN 的区别：
- PG 的 EXPLAIN ANALYZE 会**真正执行**查询返回实际数据，MySQL 的 EXPLAIN 只是估算
- PG 输出是树形结构（嵌套的执行计划），MySQL 是表格
- PG 有 Buffers 信息可以看缓存命中率，MySQL 没有

### 28. CTE 递归查询怎么用？

CTE（Common Table Expression）递归查询是 PG 的强项，适合树形结构数据：

```sql
-- 查所有下级部门
WITH RECURSIVE dept_tree AS (
    SELECT id, name, parent_id, 0 AS depth
    FROM departments WHERE id = 1          -- 根节点
    UNION ALL
    SELECT d.id, d.name, d.parent_id, dt.depth + 1
    FROM departments d
    JOIN dept_tree dt ON d.parent_id = dt.id  -- 递归
)
SELECT * FROM dept_tree;
```

MySQL 8.0 也支持 CTE 递归，但 PG 支持更早、更成熟，优化器对 CTE 的处理也更好。

PG 12 之前 CTE 默认是**优化屏障**（不会被内联到主查询），12 之后默认会内联。如果需要强制物化可以加 `MATERIALIZED`。

### 29. 窗口函数了解吗？常用的有哪些？

窗口函数在**不减少行数**的前提下做聚合计算，比 GROUP BY 灵活得多：

```sql
-- 排名（并列有间隔）
RANK() OVER (ORDER BY salary DESC)

-- 行号（无间隔）
ROW_NUMBER() OVER (ORDER BY salary DESC)

-- 分组内排名
RANK() OVER (PARTITION BY dept ORDER BY salary DESC)

-- 累计求和
SUM(salary) OVER (ORDER BY id)

-- 前后行
LAG(salary, 1) OVER (ORDER BY id)   -- 上一行
LEAD(salary, 1) OVER (ORDER BY id)  -- 下一行

-- 分桶
NTILE(4) OVER (ORDER BY salary)     -- 分成4组
```

MySQL 8.0 也支持窗口函数，但 PG 支持更早（8.4 就有了）、功能更完整。

---

## 八、全文检索

### 30. PG 的全文检索是怎么实现的？能替代 Elasticsearch 吗？

PG 内置全文检索，核心是两个类型：

**tsvector**：文本的分词向量，存储分词结果和位置信息
**tsquery**：搜索查询表达式

```sql
-- 存储预计算的 tsvector（生成列）
tsv tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', title), 'A') ||
    setweight(to_tsvector('english', content), 'B')
) STORED

-- 查询
SELECT * FROM docs WHERE tsv @@ to_tsquery('redis & cache');
-- 排名
SELECT ts_rank(tsv, query) AS rank FROM docs, to_tsquery('redis') query;
-- 高亮
SELECT ts_headline('english', content, to_tsquery('redis'));
```

`setweight` 给不同字段设权重（A 最高），标题匹配排名比正文高。

**能不能替代 ES**：
- 几十万到百万级文档、搜索需求不复杂 → PG 全文检索够用，不需要额外部署 ES
- 需要中文分词 → 要安装 zhparser 或 pg_jieba 扩展
- 千万级以上、需要复杂聚合分析 → 还是上 ES

PG 的优势是**不需要额外组件**，数据和搜索在一个库里，事务一致性有保证。

---

## 九、AI 与向量检索

### 31. pgvector 是什么？怎么用于 AI 场景？

pgvector 是 PG 的向量相似度搜索扩展，是 PG 在 AI 时代的**杀手级特性**。

核心能力：存储高维向量，做最近邻搜索。

```sql
CREATE EXTENSION vector;

CREATE TABLE documents (
    id bigserial PRIMARY KEY,
    content text,
    embedding vector(1536)  -- 1536 维向量（OpenAI text-embedding-3-small 的输出）
);

-- 余弦相似度搜索 Top 5
SELECT content, 1 - (embedding <=> query_vec) AS similarity
FROM documents
ORDER BY embedding <=> query_vec
LIMIT 5;
```

三种距离函数：
- `<->` L2 距离（欧几里得）
- `<=>` 余弦距离（**推荐**，不受向量长度影响）
- `<#>` 内积距离

两种索引：
- **IVFFlat**：先聚类再搜索，适合中等数据量（10 万~百万）
- **HNSW**：分层图结构，召回率更高，适合大数据量，但内存占用更大

### 32. RAG 是什么？pgvector 在 RAG 中扮演什么角色？

RAG（Retrieval-Augmented Generation，检索增强生成）是当前 AI 应用的主流架构：

```
用户提问 → Embedding 模型 → 查询向量
    → pgvector 相似搜索 → Top K 文档
        → 文档 + 问题 → LLM → 生成回答
```

pgvector 在其中扮演**知识库**的角色：
1. 离线阶段：文档 → embedding 模型 → 向量 → 存入 pgvector
2. 在线阶段：用户问题 → 向量 → pgvector 搜索相似文档 → 喂给 LLM

为什么用 PG 而不是 Pinecone/Milvus 等专用向量数据库：
- **一站式**：业务数据和向量在同一个库，支持 SQL 联合查询（如按分类过滤 + 向量搜索）
- **事务保证**：文档和向量的更新在同一个事务中
- **运维简单**：不需要额外部署和维护向量数据库
- **成本低**：中小规模（百万级向量）PG 完全够用

百万级以下用 pgvector 就够了。千万级以上再考虑专用向量数据库。

---

## 十、高可用架构

### 33. PG 的主从复制原理是什么？和 MySQL 有什么区别？

PG 的**流复制（Streaming Replication）**基于 WAL 日志：

```
Primary (读写) ──── WAL 日志流 ────→ Standby (只读)
```

Primary 把 WAL 实时发给 Standby，Standby 回放 WAL 同步数据。

和 MySQL 主从复制的核心区别：

**复制内容不同**：PG 流复制是**物理复制**（WAL 记录数据页变更），MySQL 是**逻辑复制**（Binlog 记录 SQL 级别变更）。物理复制延迟更低、一致性更强。

**同步模式更灵活**：

| 模式 | 行为 |
|------|------|
| 异步（默认） | Primary 不等 Standby |
| remote_write | 等 Standby 接收到 WAL 到内存 |
| on | 等 Standby 写 WAL 到磁盘 |
| remote_apply | 等 Standby 应用 WAL（**最强一致**） |

MySQL 只有异步和半同步两种。

**Standby 不可写**：PG 的流复制 Standby 严格只读。MySQL 从库可以写（虽然不推荐）。PG 如果要选择性复制或双向写，需要用逻辑复制。

### 34. PG 的高可用方案怎么选？

PG 本身不自动切换，需要外部工具：

**Patroni + etcd（生产推荐）**：

```
            etcd 集群（元数据/选举）
           ↗       ↑       ↖
      Patroni   Patroni   Patroni
         ↓         ↓         ↓
    PG Primary  PG Standby  PG Standby
         ↑                      ↑
       写请求 ←── PgBouncer ──→ 读请求
```

- Patroni 负责监控、自动故障转移、管理复制
- etcd 存集群元数据和选主信息
- PgBouncer 做连接池 + 读写分离

和 MySQL 高可用对比：
- MySQL 有原生的 Group Replication / InnoDB Cluster，自带选主
- PG 必须依赖外部工具（Patroni），但 Patroni 方案非常成熟
- PG **必须**配连接池（PgBouncer），MySQL 连接池可选

### 35. LISTEN/NOTIFY 是什么？

PG 内置的轻量级 Pub/Sub 机制：

```sql
-- 订阅
LISTEN order_changes;

-- 发布（可以配合触发器自动发送）
NOTIFY order_changes, '{"id":1,"status":"paid"}';
```

特点：
- 不持久化（fire-and-forget），未监听就丢失
- 多个 listener 可订阅同一频道
- 走 WAL，性能有上限

适合缓存失效通知、配置变更广播这类轻量场景。需要持久化、可靠投递的消息还是用 Kafka/RabbitMQ。

---

## 十一、实战场景题

### 36. 转账怎么做？和 MySQL 有区别吗？

思路一样：事务 + 悲观锁。

```sql
BEGIN;
SELECT balance FROM accounts WHERE name = 'alice' FOR UPDATE;  -- 锁住转出方
-- 检查余额
SELECT balance FROM accounts WHERE name = 'bob' FOR UPDATE;    -- 锁住转入方
UPDATE accounts SET balance = balance - 100 WHERE name = 'alice';
UPDATE accounts SET balance = balance + 100 WHERE name = 'bob';
COMMIT;
```

和 MySQL 一样要注意**按固定顺序加锁**（如按 name 升序）防死锁。

PG 的区别在于：
- 默认隔离级别是 Read Committed（MySQL 是 Repeatable Read）
- PG 没有 Gap Lock，不会出现 MySQL 那种间隙锁导致的意外阻塞
- 用 `balance = balance - 100` 做原子更新，不要读出来减了再写回去

### 37. 用 PG 做任务队列怎么实现？

PG 有天然的任务队列能力，不需要额外的 MQ：

```sql
-- 取一个待处理任务（并加锁）
SELECT * FROM tasks
WHERE status = 'pending'
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED;  -- 关键：跳过已被其他 worker 锁住的行

-- 处理完更新状态
UPDATE tasks SET status = 'done' WHERE id = ?;
COMMIT;
```

`FOR UPDATE SKIP LOCKED` 是核心：多个 worker 并发取任务，自动跳过已锁行，不会阻塞。

也可以配合 Advisory Lock 实现更细粒度的控制：用 `pg_try_advisory_lock(task_id)` 非阻塞获取锁，获取失败说明别人在处理，跳过。

适合中等吞吐量的任务调度。高吞吐量还是用 Kafka/RabbitMQ。

### 38. PG 的表膨胀问题怎么处理？

表膨胀是 PG 运维中最常见的问题。原因是 UPDATE/DELETE 产生的死元组没有及时清理。

**排查**：

```sql
SELECT relname, n_dead_tup, n_live_tup,
       round(n_dead_tup::numeric / GREATEST(n_live_tup, 1) * 100, 2) AS dead_pct
FROM pg_stat_user_tables ORDER BY n_dead_tup DESC;
```

**处理**：
1. 轻微膨胀：确保 autovacuum 正常运行，调小触发阈值
   ```
   autovacuum_vacuum_scale_factor = 0.1  -- 10% 死元组触发（默认 0.2）
   ```
2. 严重膨胀：`VACUUM FULL` 重写整表回收空间，但**会锁表**。大表用 `pg_repack` 在线重建，不锁表
3. 预防：大批量 DELETE 后手动跑 VACUUM，避免长事务（长事务会阻止 VACUUM 清理）

### 39. PG 生产环境有哪些必须监控的指标？

| 监控项 | 告警阈值 |
|--------|---------|
| 连接数 | > max_connections × 80% |
| 长事务 | > 5 分钟 |
| 表膨胀（dead_tup 占比） | > 20% |
| 复制延迟 | > 1 秒 |
| 缓存命中率 | < 99% |
| 事务 ID 消耗（`age(datfrozenxid)`） | > 5 亿（必须告警） |

**事务 ID 消耗**是最危险的——PG 事务 ID 是 32 位，约 42 亿。如果 autovacuum 没跑好，事务 ID 耗尽会导致数据库**强制停机拒绝写入**。

---

## 十二、设计与选型

### 40. PG 和 MySQL 怎么选？

**选 PostgreSQL**：
- 需要复杂查询（CTE 递归、窗口函数、多表 JOIN）
- 需要 JSONB 半结构化存储 + 查询
- 需要地理信息（PostGIS）
- 需要向量检索（pgvector，AI 场景）
- 需要 DDL 事务（频繁变更表结构）
- 需要范围类型、数组类型等高级类型
- 数据完整性要求高（金融、政府）

**选 MySQL**：
- 简单 CRUD 为主的 Web 应用
- 团队对 MySQL 更熟悉，DBA 好招
- 需要更成熟的分库分表方案（ShardingSphere、Vitess）
- 需要原生多主方案（Group Replication）
- 已有 MySQL 生态

**简单总结**：高并发简单 OLTP → MySQL；复杂查询、高级特性、AI 场景 → PostgreSQL。

### 41. PG 的主键用什么类型？

**bigserial**（自增 BIGINT）：写入性能好（堆表追加无页分裂），简单。适合单库场景。

**UUID**：`gen_random_uuid()` 生成，全局唯一。PG 的堆表是无序存储，UUID 的随机性不会像 MySQL 那样导致页分裂（MySQL 聚簇索引要求主键有序）。所以 UUID 在 PG 里做主键的性能损失比 MySQL 小得多。

**雪花 ID**：全局唯一、趋势递增。适合分布式场景。

建议：单库用 bigserial，分布式用雪花 ID，对外暴露用 UUID 或雪花 ID（不暴露自增 ID）。

### 42. PG 的连接池怎么配？为什么必须用连接池？

PG 是多进程架构，每个连接 fork 一个进程，内存开销大。不用连接池的话几百个连接就能把服务器内存吃光。

**PgBouncer**（最常用）：

```
pool_mode = transaction    -- 事务级复用（推荐）
default_pool_size = 20     -- 每用户连接数
max_client_conn = 1000     -- 最大客户端连接
```

pool_mode 三种模式：
- **session**：连接绑定整个会话（最安全但复用率低）
- **transaction**：事务结束归还连接（推荐，复用率高）
- **statement**：每条 SQL 后归还（不支持事务，不推荐）

Go 应用里 pgxpool 本身就是连接池，但如果有多个服务连同一个 PG 实例，还是建议前面加 PgBouncer 做统一管理。

### 43. PG 有哪些常用扩展？

| 扩展 | 功能 | 场景 |
|------|------|------|
| **pgvector** | 向量相似搜索 | AI embedding、推荐系统 |
| **PostGIS** | 地理空间 | 地图、LBS、围栏 |
| **pg_trgm** | 三元组相似度 | `LIKE '%keyword%'` 模糊搜索加速 |
| **pg_stat_statements** | SQL 统计 | 慢查询分析、性能监控 |
| **btree_gist** | 让 GiST 支持普通类型 | 排除约束中用 `=` |
| **pgcrypto** | 加密函数 | 数据加密、密码哈希 |
| **postgres_fdw** | 外部数据源 | 跨库查询 |
| **timescaledb** | 时序数据库 | IoT、监控指标 |
| **citus** | 分布式分片 | 水平扩展 |

PG 的扩展生态是它最大的竞争优势之一。特别是 pgvector，让 PG 成为 AI 时代的热门数据库选择。
