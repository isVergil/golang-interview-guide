# PostgreSQL 核心理论

---

## 1. PostgreSQL 是什么

PostgreSQL 是开源的**关系型数据库**，以标准合规性、可扩展性和数据完整性著称。常被称为"最先进的开源关系型数据库"。

核心特点：
- **SQL 标准合规**：对 SQL 标准支持最完整的数据库之一（支持 CTE、窗口函数、LATERAL JOIN 等）
- **MVCC 并发控制**：基于多版本实现高并发读写互不阻塞
- **丰富的数据类型**：原生支持 JSON/JSONB、数组、范围类型、几何类型、全文检索等
- **可扩展性**：支持自定义类型、自定义函数、自定义索引方法、扩展插件（如 PostGIS、pg_vector）
- **事务与数据完整性**：严格的 ACID 保证，支持 DDL 事务（建表/加列也可以回滚）

### 1.1 PostgreSQL vs MySQL 的定位

```
                 MySQL                          PostgreSQL
定位           Web 后端首选，简单高效            企业级/复杂查询首选
设计哲学       简单够用就好                     标准合规、功能完备
MVCC 实现      undo log（回滚段）               行的多版本直接存堆表
索引           B+Tree 为主                     B-Tree、GIN、GiST、BRIN、Hash、SP-GiST
JSON           JSON + 虚拟列索引               JSONB 二进制存储 + GIN 索引（更强）
全文检索       内置简单全文索引                  内置 tsvector + GIN（不依赖 ES 就能用）
DDL 事务       不支持（ALTER TABLE 自动提交）    支持（CREATE TABLE 也能回滚）
复制           binlog 逻辑复制                  WAL 物理复制 + 逻辑复制
```

### 1.2 核心架构

```
客户端连接
    ↓
Postmaster（主进程）
    ↓ fork 子进程
Backend Process（每连接一个进程）
    ↓
┌─────────────────────────────────────┐
│           共享内存                     │
│  ┌──────────┐  ┌──────────────────┐  │
│  │ Shared    │  │ WAL Buffer       │  │
│  │ Buffer    │  │（WAL 日志缓冲区）  │  │
│  │ Pool      │  │                  │  │
│  │（数据页    │  └──────────────────┘  │
│  │  缓冲区）  │  ┌──────────────────┐  │
│  │           │  │ CLOG             │  │
│  └──────────┘  │（事务提交状态）     │  │
│                └──────────────────┘  │
└─────────────────────────────────────┘
    ↓                    ↓
  数据文件             WAL 日志文件
 （堆表/索引）         （预写日志）
```

**进程模型 vs 线程模型**：

PostgreSQL 采用**多进程**架构（每个连接 fork 一个进程），MySQL 采用**多线程**架构。

| 对比 | PG（多进程） | MySQL（多线程） |
|------|-------------|----------------|
| 隔离性 | 一个进程崩溃不影响其他连接 | 一个线程崩溃可能影响整个服务 |
| 内存开销 | 每进程独占一份内存，开销大 | 线程共享内存，开销小 |
| 连接数 | 默认 100，需要连接池（pgbouncer） | 默认 151，可以开更多 |
| 上下文切换 | 进程切换开销大 | 线程切换开销小 |

> 生产环境 PostgreSQL **必须** 配合连接池（PgBouncer / Pgpool-II）使用，否则几百个连接就会把服务器内存吃光。

### 1.3 核心组件

**Shared Buffer Pool**：

数据页的缓存，类似 MySQL 的 InnoDB Buffer Pool。读数据时先查 buffer pool，命中就不用读磁盘。

```
读数据流程：
  SELECT → 查 Shared Buffer → 命中？直接返回
                              ↓ 未命中
                          从磁盘读取数据页 → 放入 Buffer → 返回
```

**WAL（Write-Ahead Logging）**：

预写日志，类似 MySQL 的 redo log。先写日志再写数据，保证崩溃后能恢复。

```
写数据流程：
  UPDATE → 修改 Buffer 中的数据页（脏页）
         → 写 WAL 日志到 WAL Buffer
         → WAL Buffer fsync 到磁盘（事务提交时）
         → 脏页由 Background Writer / Checkpointer 异步刷盘
```

**CLOG（Commit Log）**：

记录每个事务的提交状态（IN_PROGRESS / COMMITTED / ABORTED），MVCC 判断可见性时需要查这个。

---

## 2. 存储结构

PostgreSQL 的存储和 MySQL InnoDB 有本质区别：PG 是**堆表**（Heap Table），MySQL 是**聚簇索引表**（索引即数据）。

### 2.1 堆表 vs 聚簇索引表

```
PostgreSQL（堆表）：
  数据文件（Heap）：行按插入顺序存放，无序
  索引（B-Tree）：叶子节点存 ctid（行指针），指向堆表中的行
  
  索引 → ctid(0,5) → 堆表第0页第5行 → 取数据

MySQL InnoDB（聚簇索引）：
  主键索引（B+Tree）：叶子节点直接存整行数据
  二级索引：叶子节点存主键值，需要回表
  
  二级索引 → 主键值 → 主键索引 → 取数据
```

| 对比 | PG 堆表 | MySQL 聚簇索引 |
|------|---------|---------------|
| 主键查询 | 走索引 → ctid → 回堆表（多一次 IO） | 直接在主键索引叶子节点取数据 |
| 写入性能 | 堆表无序追加，写入快 | 主键有序插入才快，乱序写入导致页分裂 |
| UPDATE | 写新版本到堆表任意位置（可能同页 HOT） | 原地更新（undo log 存旧版本） |
| 表膨胀 | 有（旧版本残留在堆表） | 无（undo log 单独管理） |
| 全表扫描 | 顺序扫描堆表，效率高（连续 IO） | 扫主键索引叶子节点 |

### 2.2 页结构（Page）

PostgreSQL 数据文件按 **8KB** 页（block）为单位组织（MySQL InnoDB 是 16KB）。

```
8KB 页结构：

+--------------------+
| PageHeader (24B)   |  页头：LSN、校验和、空闲空间指针
+--------------------+
| ItemId 数组         |  行指针数组（每个 4 字节），从前往后增长
| (item1, item2, ...)| 
+--------------------+
|     空闲空间         |  ItemId 从前往后长，Tuple 从后往前长
|     Free Space      |  两者相遇 = 页满
+--------------------+
| Tuple 数据          |  实际行数据，从页尾往前存
| (tuple2, tuple1)   |
+--------------------+
| Special Space      |  索引页使用（数据页通常为空）
+--------------------+
```

**ctid（行指针）**：

每行数据的物理地址，格式 `(page, offset)`，比如 `(0, 5)` 表示第 0 页第 5 个 item。索引叶子节点存的就是 ctid。

```sql
-- 查看行的物理位置
SELECT ctid, * FROM users LIMIT 5;
-- (0,1)  | 1 | alice | ...
-- (0,2)  | 2 | bob   | ...
```

### 2.3 行结构（Tuple）

每一行数据在 PG 里叫 Tuple，包含**行头**和**数据**两部分：

```
Tuple 结构：

+-------------------------------------------+
| HeapTupleHeader (23 bytes)                |
|   t_xmin    - 插入该行的事务 ID              |
|   t_xmax    - 删除/更新该行的事务 ID          |
|   t_cid     - 命令序号（同一事务内的顺序）     |
|   t_ctid    - 当前行的 ctid（更新后指向新版本）|
|   t_infomask - 标志位（事务状态、NULL 信息等） |
+-------------------------------------------+
| NULL bitmap (可选)                         |
+-------------------------------------------+
| 用户数据 (col1, col2, col3, ...)           |
+-------------------------------------------+
```

**t_xmin 和 t_xmax 是 MVCC 的核心**：

- `t_xmin`：创建这行的事务 ID，行可见的起点
- `t_xmax`：删除这行的事务 ID，行不可见的起点（0 表示未被删除）

```
INSERT：t_xmin = 当前事务ID, t_xmax = 0
DELETE：把 t_xmax 设为当前事务ID（行并不物理删除）
UPDATE：旧行 t_xmax = 当前事务ID + 插入新行 t_xmin = 当前事务ID
         旧行 t_ctid 指向新行（形成版本链）
```

### 2.4 TOAST（大字段存储）

当一行数据超过页大小的 1/4（约 2KB）时，PG 会自动用 TOAST（The Oversized-Attribute Storage Technique）处理：

| 策略 | 说明 |
|------|------|
| PLAIN | 不压缩不外存（定长类型如 int） |
| EXTENDED | 先压缩，还放不下就外存到 TOAST 表（默认，text/jsonb） |
| EXTERNAL | 不压缩，直接外存（需要快速取子串时用） |
| MAIN | 先压缩，尽量不外存 |

```
普通 Tuple：  [header | col1 | col2 | col3(小)]
TOAST Tuple： [header | col1 | col2 | col3 → 指向 TOAST 表的指针]
                                              ↓
                                         TOAST 表（单独的存储）
                                         [chunk1][chunk2][chunk3]
```

---

## 3. 数据类型

PostgreSQL 的类型系统远比 MySQL 丰富，很多需要应用层实现的功能，PG 在数据库层就原生支持了。

### 3.1 基础类型

**数值类型**：

| 类型 | 大小 | 范围 | 说明 |
|------|------|------|------|
| smallint (int2) | 2 字节 | -32768 ~ 32767 | 小整数 |
| integer (int4) | 4 字节 | -21 亿 ~ 21 亿 | 常用整数 |
| bigint (int8) | 8 字节 | -922 京 ~ 922 京 | 大整数、雪花 ID |
| numeric(p,s) | 可变 | 任意精度 | 金额（精确计算） |
| real | 4 字节 | 6 位精度 | 浮点（不精确） |
| double precision | 8 字节 | 15 位精度 | 浮点（不精确） |
| serial / bigserial | 4/8 字节 | 自增 | 自增主键（本质是 sequence） |

> 金额用 `numeric`（不是 float），或者和 MySQL 一样用 `bigint` 存分。

**字符类型**：

| 类型 | 说明 |
|------|------|
| varchar(n) | 变长，最大 n 个字符（n 可选，不填则不限长度） |
| char(n) | 定长，不足补空格（几乎不用） |
| text | 变长，无长度限制（和不带 n 的 varchar 等价） |

> PG 中 `varchar` 和 `text` 性能完全相同，底层实现一样。不需要像 MySQL 那样纠结 varchar(255)。

**时间类型**：

| 类型 | 说明 |
|------|------|
| timestamp | 不带时区（存什么就是什么） |
| timestamptz | 带时区（存入时转 UTC，取出时按会话时区转换）**推荐** |
| date | 日期 |
| interval | 时间间隔（如 '3 days 2 hours'） |

> 生产环境统一用 `timestamptz`，避免时区问题。

### 3.2 JSON / JSONB

PG 原生支持两种 JSON 类型：

| 类型 | 存储 | 查询速度 | 索引 | 使用场景 |
|------|------|---------|------|---------|
| json | 文本原样保存 | 慢（每次解析） | 不支持 | 只存不查 |
| jsonb | 二进制格式 | 快（预解析） | GIN 索引 | **推荐，查询+索引** |

```sql
-- 创建 JSONB 列
CREATE TABLE events (
    id bigserial PRIMARY KEY,
    data jsonb NOT NULL
);

-- 写入
INSERT INTO events (data) VALUES ('{"type":"click","page":"/home","user_id":1001}');

-- 查询 JSONB 字段
SELECT data->>'type' AS event_type FROM events;           -- 取文本值
SELECT data->'user_id' FROM events;                       -- 取 JSON 值
SELECT * FROM events WHERE data @> '{"type":"click"}';    -- 包含查询
SELECT * FROM events WHERE data->>'page' = '/home';       -- 字段值查询

-- GIN 索引加速 JSONB 查询
CREATE INDEX idx_events_data ON events USING GIN (data);
-- 加索引后 @>（包含）、?（key 存在）、?|（任一 key 存在）等操作走索引
```

**JSONB 常用运算符**：

| 运算符 | 含义 | 示例 |
|--------|------|------|
| `->` | 取 JSON 对象字段（返回 JSON） | `data->'name'` |
| `->>` | 取 JSON 对象字段（返回文本） | `data->>'name'` |
| `@>` | 左边包含右边 | `data @> '{"type":"click"}'` |
| `<@` | 左边被右边包含 | `data <@ '{"type":"click","page":"/"}'` |
| `?` | 是否存在某个 key | `data ? 'type'` |
| `\|\|` | 合并两个 JSONB | `data \|\| '{"extra":1}'` |
| `- key` | 删除某个 key | `data - 'page'` |

### 3.3 数组类型

PG 原生支持数组，不需要拆表：

```sql
-- 数组列
CREATE TABLE articles (
    id serial PRIMARY KEY,
    title text,
    tags text[]    -- 文本数组
);

INSERT INTO articles (title, tags) VALUES ('Go 入门', ARRAY['go', 'tutorial', 'backend']);

-- 查询
SELECT * FROM articles WHERE 'go' = ANY(tags);           -- 包含某个元素
SELECT * FROM articles WHERE tags @> ARRAY['go','backend']; -- 包含所有
SELECT * FROM articles WHERE tags && ARRAY['go','python'];  -- 有交集

-- GIN 索引加速数组查询
CREATE INDEX idx_articles_tags ON articles USING GIN (tags);
```

**unnest —— 数组展开为行**：

`unnest()` 把数组的每个元素展开成独立的一行，是数组类型最重要的配套函数：

```sql
-- 基本用法：一维数组展开
SELECT unnest(ARRAY['go','redis','pgsql']);
--  go
--  redis
--  pgsql

-- 配合表数据：把每篇文章的标签展开成行
SELECT id, title, unnest(tags) AS tag FROM articles;
-- 1 | Go 入门 | go
-- 1 | Go 入门 | tutorial
-- 1 | Go 入门 | backend

-- 统计所有标签出现次数（先展开再聚合）
SELECT tag, count(*) AS cnt
FROM (SELECT unnest(tags) AS tag FROM articles) t
GROUP BY tag ORDER BY cnt DESC;

-- 反向操作：array_agg 把行聚合回数组
SELECT array_agg(DISTINCT tag ORDER BY tag)
FROM (SELECT unnest(tags) AS tag FROM articles) t;
```

**unnest 用于批量写入（pgx 推荐写法）**：

传统批量 INSERT 需要拼 `VALUES ($1,$2),($3,$4),...`，PG 的 `unnest` 可以把多个数组"zip"展开成行，一条 SQL 完成批量插入：

```sql
-- 三个数组对齐展开，等价于多行 VALUES
INSERT INTO users (name, email, age)
SELECT * FROM unnest(
    ARRAY['alice','bob','carol'],          -- name[]
    ARRAY['a@test.com','b@test.com','c@test.com'],  -- email[]
    ARRAY[25, 30, 28]                      -- age[]
);

-- 等价于：
-- INSERT INTO users (name, email, age) VALUES
--   ('alice', 'a@test.com', 25),
--   ('bob',   'b@test.com', 30),
--   ('carol', 'c@test.com', 28);
```

> Go 中配合 pgx 使用时，直接传 Go 切片即可：`pool.Exec(ctx, sql, names, emails, ages)`，pgx 自动把 `[]string` / `[]int` 映射为 PG 数组。这种方式比拼 SQL 更安全（防注入）、性能更好（一次网络往返）、无参数数量上限。

### 3.4 范围类型

表示一个值的范围，天然适合时间段、价格区间等场景：

```sql
-- 会议室预订（时间范围不允许重叠）
CREATE TABLE bookings (
    id serial PRIMARY KEY,
    room text,
    during tstzrange,  -- 时间戳范围
    EXCLUDE USING GIST (room WITH =, during WITH &&)  -- 排除约束：同房间时间不重叠
);

INSERT INTO bookings (room, during) VALUES ('A101', '[2026-03-24 09:00, 2026-03-24 11:00)');
-- 再插入重叠时间段会自动报错
INSERT INTO bookings (room, during) VALUES ('A101', '[2026-03-24 10:00, 2026-03-24 12:00)');
-- ERROR: conflicting key value violates exclusion constraint

-- 范围运算符
SELECT * FROM bookings WHERE during @> '2026-03-24 10:00'::timestamptz;  -- 包含某时间点
SELECT * FROM bookings WHERE during && '[2026-03-24 08:00, 2026-03-24 10:00)'; -- 有重叠
```

| 范围类型 | 说明 |
|---------|------|
| int4range | 整数范围 |
| int8range | 大整数范围 |
| numrange | 数值范围 |
| tsrange | 时间戳范围（无时区） |
| tstzrange | 时间戳范围（带时区） |
| daterange | 日期范围 |

**范围运算符一览**：

| 运算符 | 含义 | 示例 |
|--------|------|------|
| `@>` 值 | 范围包含某个点 | `during @> '2026-03-24 10:00'::timestamptz` |
| `@>` 范围 | 范围完全包含另一个范围 | `during @> '[09:00, 10:00)'` |
| `&&` | 两个范围有重叠 | `during && '[08:00, 10:00)'` |
| `<<` | 严格在左边（无重叠） | `during << '[12:00, 14:00)'` |
| `*` | 交集 | `'[9:00,12:00)' * '[10:00,14:00)'` → `[10:00,12:00)` |
| `+` | 并集 | `'[9:00,11:00)' + '[10:00,12:00)'` → `[9:00,12:00)` |

**范围类型 vs MySQL 两字段方案（start_time / end_time）**：

MySQL 没有范围类型，只能用两个字段模拟时间段。PG 的范围类型在以下方面有明显优势：

1）**重叠检测**：
```sql
-- MySQL：多个 AND 条件，容易写错边界
WHERE room = 'A101' AND start_time < '12:00' AND end_time > '10:00'

-- PG：一个运算符
WHERE room = 'A101' AND during && '[10:00, 12:00)'
```

2）**数据库层约束**：
```sql
-- MySQL：无法在数据库层防止时间重叠，必须应用层加锁 → 查冲突 → 插入
-- PG：排除约束自动拒绝冲突数据，并发安全，无需应用层加锁
EXCLUDE USING GIST (room WITH =, during WITH &&)
```

3）**开闭区间原生支持**：
```sql
-- PG 原生区分开闭区间，语义精确
'[9:00, 11:00)'   -- 左闭右开：包含 9:00，不包含 11:00
'(9:00, 11:00]'   -- 左开右闭
'[9:00, 11:00]'   -- 全闭

-- MySQL 没有范围类型，开闭全靠 < / <= 人工控制，不同人写法不一致容易出 bug
```

4）**索引支持**：PG 的 GiST 索引原生支持范围运算符（`&&`、`@>`），MySQL 的 B+Tree 只能索引单列，两列时间范围查询无法高效走索引。

5）**交集/并集/差集运算**：PG 原生 `*`、`+` 运算符直接计算，MySQL 需要 `GREATEST/LEAST` 在应用层拼接。

> 总结：PG 把"范围"作为一等公民类型，查询语义清晰、约束由数据库保证、索引原生支持；MySQL 只能用两个字段模拟，逻辑散落在应用层。

### 3.5 其他特色类型

| 类型 | 说明 | 场景 |
|------|------|------|
| uuid | UUID v4 | 分布式主键（用 `gen_random_uuid()` 生成） |
| inet / cidr | IP 地址 / 网段 | 网络配置、IP 白名单 |
| hstore | 简单键值对 | 轻量级 key-value（JSONB 功能更全） |
| tsvector / tsquery | 全文检索 | 搜索引擎（不依赖 ES 就能用） |
| point / polygon / geometry | 几何类型 | GIS（配合 PostGIS 扩展） |
| enum | 枚举 | 状态字段（不推荐，改起来麻烦） |

---

## 4. 索引

PostgreSQL 支持**6 种**索引类型，远多于 MySQL（基本只有 B+Tree）。

### 4.1 B-Tree 索引（默认）

> 本质和 MySQL 的 B+Tree 类似，适合等值查询和范围查询。PG 的 B-Tree 叶子节点存 ctid（行指针），不像 MySQL 聚簇索引叶子节点直接存数据。

```sql
-- 默认就是 B-Tree
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_age ON users (age);

-- 适合的查询
SELECT * FROM users WHERE email = 'alice@test.com';  -- 等值
SELECT * FROM users WHERE age BETWEEN 20 AND 30;     -- 范围
SELECT * FROM users ORDER BY age;                     -- 排序
```

**复合索引与最左前缀**：

```sql
CREATE INDEX idx_age_balance ON users (age, balance);

-- 走索引
WHERE age = 25                        -- 最左列
WHERE age = 25 AND balance > 1000     -- 最左 + 后续列
WHERE age > 20 ORDER BY age           -- 最左列范围

-- 不走索引
WHERE balance > 1000                  -- 跳过了最左列
```

### 4.2 Hash 索引

> 本质：哈希表，只支持等值查询（=），不支持范围、排序。PG 10 之后 Hash 索引支持 WAL，可以安全使用了。

```sql
CREATE INDEX idx_users_email_hash ON users USING HASH (email);

-- 走索引
WHERE email = 'alice@test.com'

-- 不走索引
WHERE email LIKE 'alice%'    -- 范围/模糊
ORDER BY email               -- 排序
```

适用场景：超大表的等值查询，比 B-Tree 少一点空间和稍快一点。但大多数场景 B-Tree 就够了。

### 4.3 GIN 索引（通用倒排索引）

> 本质：倒排索引（Inverted Index），和搜索引擎的原理一样。每个值维护一个 posting list（包含该值的行列表）。

适合：**多值列**的查询，如 JSONB、数组、全文检索。

```sql
-- JSONB 索引
CREATE INDEX idx_events_data ON events USING GIN (data);
SELECT * FROM events WHERE data @> '{"type":"click"}';   -- 包含查询走 GIN

-- 数组索引
CREATE INDEX idx_tags ON articles USING GIN (tags);
SELECT * FROM articles WHERE tags @> ARRAY['go'];         -- 数组包含走 GIN

-- 全文检索索引
CREATE INDEX idx_content_fts ON articles USING GIN (to_tsvector('english', content));
SELECT * FROM articles WHERE to_tsvector('english', content) @@ to_tsquery('redis & cache');
```

**GIN 的工作原理**：

```
文档1: "Redis 缓存设计"  →  tsvector: 'redis' 'cache' 'design'
文档2: "Redis 分布式锁"  →  tsvector: 'redis' 'distributed' 'lock'
文档3: "MySQL 索引优化"  →  tsvector: 'mysql' 'index' 'optimize'

GIN 倒排索引：
  'redis'       → [文档1, 文档2]
  'cache'       → [文档1]
  'distributed' → [文档2]
  'mysql'       → [文档3]
  'index'       → [文档3]
```

### 4.4 GiST 索引（通用搜索树）

> 本质：平衡搜索树，但不是 B-Tree。支持重叠、包含、距离等二维/多维查询。

适合：范围类型、几何类型、全文检索、PostGIS 空间查询。

```sql
-- 范围类型排除约束
EXCLUDE USING GIST (room WITH =, during WITH &&);

-- PostGIS 空间查询
CREATE INDEX idx_geom ON places USING GIST (location);
SELECT * FROM places WHERE ST_DWithin(location, ST_MakePoint(116.4, 39.9), 1000);

-- 全文检索（GiST 也支持，比 GIN 小但查询慢）
CREATE INDEX idx_content_gist ON articles USING GIST (to_tsvector('english', content));
```

### 4.5 BRIN 索引（块范围索引）

> 本质：不索引每一行，而是记录每个**数据页范围**（如每 128 页）的最小值和最大值。体积极小。

适合：数据物理有序的大表（如时序数据，按时间插入的日志表）。

```sql
-- 时序表，数据按时间顺序插入
CREATE TABLE logs (
    id bigserial,
    created_at timestamptz DEFAULT now(),
    message text
);

-- BRIN 索引（比 B-Tree 小几百倍）
CREATE INDEX idx_logs_time ON logs USING BRIN (created_at);

-- 查最近一小时的日志
SELECT * FROM logs WHERE created_at > now() - interval '1 hour';
```

**BRIN vs B-Tree 大小对比（1 亿行时序表）**：

```
B-Tree 索引：~2 GB
BRIN 索引：  ~100 KB（小 2 万倍）
```

前提是数据物理有序。如果数据乱序插入，BRIN 会退化，扫描大量无效页。

### 4.6 SP-GiST 索引

> 本质：空间分区搜索树（四叉树、kd-tree 等）。适合非平衡数据分布的空间查询。

使用较少，了解即可。

### 4.7 六种索引对比

| 索引类型 | 适合查询 | 典型场景 | 体积 |
|---------|---------|---------|------|
| B-Tree | =, <, >, BETWEEN, ORDER BY | 大多数场景（默认） | 中等 |
| Hash | = | 超大表精确匹配 | 小 |
| GIN | @>, ?, @@, &&（多值包含） | JSONB、数组、全文检索 | 大 |
| GiST | &&, @>, <->（范围/空间） | 范围类型、PostGIS | 中等 |
| BRIN | 范围查询（数据物理有序） | 时序数据、日志表 | 极小 |
| SP-GiST | 空间分区 | 非均匀空间数据 | 中等 |

### 4.8 索引进阶用法

**部分索引（Partial Index）**：

只索引满足条件的行，减小索引体积：

```sql
-- 只索引未删除的行
CREATE INDEX idx_active_users ON users (email) WHERE deleted_at IS NULL;

-- 只索引 VIP 用户
CREATE INDEX idx_vip ON users (name) WHERE level = 'vip';
```

**表达式索引**：

索引一个表达式的结果：

```sql
-- 索引小写 email（解决大小写不敏感查询）
CREATE INDEX idx_lower_email ON users (lower(email));
SELECT * FROM users WHERE lower(email) = 'alice@test.com';  -- 走索引

-- 索引 JSONB 某个字段
CREATE INDEX idx_data_type ON events ((data->>'type'));
SELECT * FROM events WHERE data->>'type' = 'click';  -- 走索引
```

**覆盖索引（INCLUDE）**：

B-Tree 索引可以额外携带列，避免回堆表：

```sql
-- 索引 email，额外携带 name（不用回表）
CREATE INDEX idx_email_include ON users (email) INCLUDE (name);
SELECT name FROM users WHERE email = 'alice@test.com';  -- Index Only Scan
```

---

## 5. MVCC（多版本并发控制）

MVCC 是 PostgreSQL 实现高并发的核心机制：**读不阻塞写，写不阻塞读**。PG 的 MVCC 实现和 MySQL 有本质区别。

### 5.1 PG 的 MVCC vs MySQL 的 MVCC

```
PostgreSQL MVCC：
  旧版本直接留在堆表中（通过 t_xmin/t_xmax 标记可见性）
  需要 VACUUM 清理旧版本（死元组）

MySQL InnoDB MVCC：
  旧版本存在 undo log 中，堆表只有最新版本
  undo log 由 purge 线程自动清理
```

| 对比 | PostgreSQL | MySQL InnoDB |
|------|-----------|-------------|
| 旧版本存放 | 堆表内（原地） | undo log（单独区域） |
| UPDATE 实现 | 标记旧行删除 + 写入新行 | 原地更新 + 旧值写 undo log |
| 清理机制 | VACUUM（手动/autovacuum） | purge 线程自动清理 |
| 表膨胀 | 有（旧版本占空间） | 无 |
| 读旧版本 | 直接读堆表中的旧行 | 通过 undo log 链构造旧版本 |
| 优点 | 旧版本在堆表，读取快 | 不膨胀，undo 自动清理 |
| 缺点 | 需要 VACUUM，表会膨胀 | 长事务 undo log 积压 |

### 5.2 可见性判断

一行数据对某个事务是否可见，由 t_xmin、t_xmax 和当前事务的**快照**决定：

```
规则（简化版）：
  1. t_xmin 的事务已提交 且 在我的快照之前 → 行"出生"了，看得见
  2. t_xmax = 0（未被删除） → 行还"活着"
  3. t_xmax 的事务已提交 且 在我的快照之前 → 行已"死亡"，看不见

综合：一行可见 = t_xmin 已提交且在快照前 && (t_xmax=0 || t_xmax 未提交 || t_xmax 在快照后)
```

**具体例子**：

```
时间线：
  事务 100: INSERT INTO users VALUES (1, 'alice');   -- t_xmin=100, t_xmax=0
  事务 100: COMMIT;

  事务 200: UPDATE users SET name='bob' WHERE id=1;  
            -- 旧行: t_xmin=100, t_xmax=200
            -- 新行: t_xmin=200, t_xmax=0
  事务 200: COMMIT;

  事务 300: SELECT * FROM users WHERE id=1;
            -- 快照包含已提交的 100 和 200
            -- 旧行: t_xmin=100 ✓ t_xmax=200 已提交 → 不可见
            -- 新行: t_xmin=200 ✓ t_xmax=0       → 可见 → 返回 'bob'
```

### 5.3 事务 ID 回卷（Transaction ID Wraparound）

PG 的事务 ID 是 32 位无符号整数，最大约 42 亿。用完后会回卷（wraparound），导致旧数据"消失"（新事务 ID 比旧事务 ID 小，可见性判断出错）。

**FREEZE 机制**：

VACUUM FREEZE 会把非常老的行的 t_xmin 标记为 `FrozenTransactionId`（特殊值 2），表示"对所有事务都可见"，这样就不受事务 ID 大小比较的影响了。

```
autovacuum_freeze_max_age = 200000000（默认 2 亿）
→ 当表中最老的事务 ID 和当前事务 ID 差距超过 2 亿时
→ 自动触发 aggressive VACUUM FREEZE
```

> 生产环境必须保证 autovacuum 正常运行。如果 autovacuum 被关闭或卡住，事务 ID 耗尽会导致数据库**强制关机**（拒绝一切写入）。

### 5.4 VACUUM

VACUUM 是 PostgreSQL **独有的**维护操作，MySQL 没有。它负责清理堆表中的**死元组**（dead tuples，被 DELETE/UPDATE 产生的旧版本）。

**为什么需要 VACUUM**：

```
初始状态：users 表 1000 行，占 10 个页

执行 DELETE FROM users WHERE age < 18;  -- 删了 200 行

表状态：1000 行（其中 200 行是死元组），仍占 10 个页
         ↑ 死元组不释放空间，新 INSERT 也不能复用这些位置

VACUUM 后：死元组标记为可复用，新 INSERT 可以使用这些空间
VACUUM FULL：重写整张表，物理上回收空间（但会锁表）
```

| VACUUM 类型 | 作用 | 是否锁表 | 场景 |
|------------|------|---------|------|
| VACUUM | 标记死元组空间可复用 | 不锁表 | autovacuum 自动执行 |
| VACUUM FULL | 重写整表，物理回收空间 | **锁表** | 表膨胀严重时手动执行 |
| VACUUM FREEZE | 冻结旧事务 ID | 不锁表 | 防止事务 ID 回卷 |

**autovacuum（自动清理）**：

PG 内置 autovacuum worker，默认开启，根据死元组比例自动触发：

```
触发条件：dead_tuples > autovacuum_vacuum_threshold + autovacuum_vacuum_scale_factor × 总行数
默认值：  dead_tuples > 50 + 0.2 × 总行数

例：1 万行的表，死元组超过 50 + 2000 = 2050 行时触发
```

### 5.5 HOT 更新（Heap Only Tuple）

HOT 是 PG 对 UPDATE 的优化：如果更新的列**没有索引**，且新行能放在**同一页**内，就不用更新索引。

```
普通 UPDATE：
  旧行标记删除 → 新行写入（可能不同页）→ 所有索引更新指向新行

HOT UPDATE：
  旧行标记删除 → 新行写入同一页 → 旧行 ctid 指向新行 → 索引不用动
  索引仍指向旧行，通过 ctid 链跳转到新行
```

HOT 的条件：
1. 更新的列不在任何索引中
2. 新行能放进同一页（页有足够空间）

> 表设计时，高频更新的列尽量不要建索引，以提高 HOT 命中率。

---

## 6. 事务与锁

### 6.1 事务隔离级别

PG 支持四个 SQL 标准隔离级别，但实现上只有三种（Read Uncommitted 行为等同于 Read Committed）：

| 隔离级别 | 脏读 | 不可重复读 | 幻读 | PG 实现 |
|---------|------|----------|------|---------|
| Read Uncommitted | ✗ | ✓ | ✓ | 等同于 Read Committed |
| Read Committed | ✗ | ✓ | ✓ | **默认**，每条 SQL 取最新快照 |
| Repeatable Read | ✗ | ✗ | ✗ | 事务开始时取快照，整个事务用同一个 |
| Serializable | ✗ | ✗ | ✗ | SSI（可序列化快照隔离） |

> PG 的 Repeatable Read 比 MySQL 更严格——PG 没有幻读（真正的快照隔离），MySQL 的 RR 通过 gap lock 防幻读但仍有边界情况。

**Read Committed vs Repeatable Read**：

```sql
-- Read Committed（默认）：每条 SQL 看到最新已提交数据
BEGIN;
SELECT count(*) FROM users;  -- 快照1: 100
-- 另一个事务 INSERT 一行并 COMMIT
SELECT count(*) FROM users;  -- 快照2: 101（看到了新提交的数据）
COMMIT;

-- Repeatable Read：整个事务用同一个快照
BEGIN ISOLATION LEVEL REPEATABLE READ;
SELECT count(*) FROM users;  -- 快照: 100
-- 另一个事务 INSERT 一行并 COMMIT
SELECT count(*) FROM users;  -- 还是 100（快照不变）
COMMIT;
```

### 6.2 锁机制

**表级锁**：

PG 有 8 种表级锁，常用的：

| 锁模式 | 触发语句 | 冲突 |
|--------|---------|------|
| ACCESS SHARE | SELECT | 只和 ACCESS EXCLUSIVE 冲突 |
| ROW SHARE | SELECT FOR UPDATE/SHARE | |
| ROW EXCLUSIVE | INSERT/UPDATE/DELETE | |
| ACCESS EXCLUSIVE | ALTER TABLE, DROP, VACUUM FULL | 和所有锁冲突 |

> `ALTER TABLE` 需要 ACCESS EXCLUSIVE 锁——生产环境加列/改列会**锁全表**。大表 DDL 变更需要特别注意。

**行级锁**：

```sql
-- SELECT FOR UPDATE：排他锁（写锁），其他事务读可以但改会阻塞
SELECT * FROM users WHERE id = 1 FOR UPDATE;

-- SELECT FOR SHARE：共享锁（读锁），其他事务读可以但改会阻塞
SELECT * FROM users WHERE id = 1 FOR SHARE;

-- FOR UPDATE NOWAIT：加锁失败立即报错（不等待）
SELECT * FROM users WHERE id = 1 FOR UPDATE NOWAIT;

-- FOR UPDATE SKIP LOCKED：跳过已锁定的行（队列消费场景）
SELECT * FROM tasks WHERE status = 'pending' LIMIT 1 FOR UPDATE SKIP LOCKED;
```

**Advisory Lock（咨询锁）**：

应用层自定义锁，PG 只提供锁基础设施，业务自己决定锁的含义：

```sql
-- 获取锁（key 是一个 bigint）
SELECT pg_advisory_lock(12345);

-- 尝试获取（非阻塞）
SELECT pg_try_advisory_lock(12345);

-- 释放
SELECT pg_advisory_unlock(12345);
```

典型场景：分布式锁（类似 Redis SETNX，但不依赖额外组件）。

### 6.3 死锁

PG 有**死锁检测器**，默认每 1 秒检测一次（`deadlock_timeout = 1s`）。发现死锁后自动回滚其中一个事务。

```sql
-- 事务A                              -- 事务B
BEGIN;                                BEGIN;
UPDATE users SET age=1 WHERE id=1;    UPDATE users SET age=2 WHERE id=2;
UPDATE users SET age=1 WHERE id=2;    UPDATE users SET age=2 WHERE id=1;
-- 等待 B 释放 id=2                    -- 等待 A 释放 id=1
-- 死锁！PG 检测到后回滚其中一个事务
```

预防：所有事务按**相同顺序**访问资源（如按 id 升序加锁）。

---

## 7. 查询优化

### 7.1 EXPLAIN 分析

```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT * FROM users WHERE age > 25;
```

关键字段：

```
Seq Scan on users  (cost=0.00..1.05 rows=3 width=72) (actual time=0.01..0.02 rows=3 loops=1)
  Filter: (age > 25)
  Rows Removed by Filter: 2
  Buffers: shared hit=1
Planning Time: 0.05 ms
Execution Time: 0.03 ms
```

| 字段 | 含义 |
|------|------|
| cost | 启动代价..总代价（单位：page fetch 的开销） |
| rows | 估算返回行数 |
| actual time | 实际执行时间（毫秒） |
| Buffers: shared hit | 命中缓存的页数 |
| Buffers: shared read | 从磁盘读取的页数 |

**常见扫描方式**：

| 扫描方式 | 含义 | 好坏 |
|---------|------|------|
| Seq Scan | 全表扫描 | 小表正常，大表要优化 |
| Index Scan | B-Tree 索引扫描 → 回堆表取数据 | 正常 |
| Index Only Scan | 覆盖索引，不回表 | 最快 |
| Bitmap Index Scan | 索引扫描构建位图 → 批量回表 | 多条件 OR / 返回行多时 |
| Bitmap Heap Scan | 按位图批量取堆表数据 | 配合 Bitmap Index Scan |

### 7.2 常见优化手段

**统计信息**：

PG 的查询计划依赖**统计信息**，如果统计信息过时，优化器会选错计划。

```sql
-- 手动更新统计信息
ANALYZE users;

-- 查看统计信息
SELECT * FROM pg_stats WHERE tablename = 'users';
```

**常见优化**：

| 问题 | 解决方案 |
|------|---------|
| 全表扫描 | 加索引，或检查索引是否失效（类型不匹配、函数运算） |
| Index Scan 但仍然慢 | 检查是否返回太多行（选择性低），考虑覆盖索引 |
| Nested Loop 太慢 | 大表 JOIN 考虑 Hash Join（确保 work_mem 足够） |
| 排序耗内存 | 增加 work_mem，或建带排序的索引 |
| 统计信息过时 | 执行 ANALYZE，或调整 autovacuum 频率 |
| LIKE '%keyword%' | 使用 pg_trgm 扩展 + GIN 索引 |

### 7.3 CTE（公用表表达式）

```sql
-- 递归查询（树形结构：查所有下级部门）
WITH RECURSIVE dept_tree AS (
    -- 根节点
    SELECT id, name, parent_id, 0 AS depth
    FROM departments WHERE id = 1
    
    UNION ALL
    
    -- 递归
    SELECT d.id, d.name, d.parent_id, dt.depth + 1
    FROM departments d
    JOIN dept_tree dt ON d.parent_id = dt.id
)
SELECT * FROM dept_tree;
```

### 7.4 窗口函数

```sql
-- 排名（并列排名有间隔）
SELECT name, balance, RANK() OVER (ORDER BY balance DESC) AS rk FROM users;

-- 密集排名（并列排名无间隔）
SELECT name, balance, DENSE_RANK() OVER (ORDER BY balance DESC) AS rk FROM users;

-- 行号
SELECT name, ROW_NUMBER() OVER (ORDER BY balance DESC) AS rn FROM users;

-- 分组内排名（每个部门内的工资排名）
SELECT name, dept, salary,
       RANK() OVER (PARTITION BY dept ORDER BY salary DESC) AS dept_rank
FROM employees;

-- 累计求和
SELECT name, balance,
       SUM(balance) OVER (ORDER BY id) AS running_total
FROM users;

-- 前后行
SELECT name, balance,
       LAG(balance, 1) OVER (ORDER BY id) AS prev_balance,
       LEAD(balance, 1) OVER (ORDER BY id) AS next_balance
FROM users;
```

---

## 8. 高可用架构

### 8.1 流复制（Streaming Replication）

> PG 原生的高可用方案，基于 WAL 日志流复制，类似 MySQL 的 binlog 复制。

```
                    WAL 日志流
Primary (读写) ──────────────→ Standby (只读)
     ↑                              ↑
   写请求                         读请求
```

**同步模式**：

| 模式 | 特点 |
|------|------|
| 异步复制（默认） | Primary 不等 Standby 确认，性能好但可能丢数据 |
| 同步复制 | Primary 等 Standby 写入 WAL 后才返回，零丢失但延迟高 |
| 准同步复制 | Primary 等 Standby 接收 WAL（不用写盘），折中方案 |

```
-- postgresql.conf (Primary)
synchronous_commit = on           -- 同步级别
synchronous_standby_names = 'sb1' -- 同步 Standby 名称

-- 同步级别选项：
-- off          : 异步（最快，可能丢数据）
-- on           : 等待 Standby 写 WAL 到磁盘
-- remote_write : 等待 Standby 接收 WAL 到内存
-- remote_apply : 等待 Standby 应用 WAL（最强一致）
```

**故障切换**：

PG 本身不自动切换，需要配合工具：

```
Patroni（推荐）：
  Patroni + etcd/ZooKeeper/Consul
  自动故障检测、自动切换、自动注册新 Standby

pg_auto_failover：
  PG 官方的 HA 扩展，配置简单

手动切换：
  pg_ctl promote -D /data/standby   -- 手动提升 Standby 为 Primary
```

### 8.2 逻辑复制

> 基于行级别的变更复制（类似 MySQL 的 row-based binlog），可以选择性复制部分表。

```sql
-- 发布端（Primary）
CREATE PUBLICATION my_pub FOR TABLE users, orders;

-- 订阅端（Standby 或另一个集群）
CREATE SUBSCRIPTION my_sub
    CONNECTION 'host=primary port=5432 dbname=mydb'
    PUBLICATION my_pub;
```

| 对比 | 流复制 | 逻辑复制 |
|------|--------|---------|
| 复制粒度 | 整个实例 | 可选表 |
| 版本要求 | 主从版本一致 | 可跨版本 |
| Standby 可写 | 不可以 | 可以（独立实例） |
| 用途 | HA 高可用 | 数据同步、跨版本升级、多活 |

### 8.3 高可用架构方案

**Patroni + etcd（生产推荐）**：

```
                   etcd 集群
                   (元数据/选举)
                  ↗     ↑     ↖
            Patroni   Patroni   Patroni
               ↓         ↓         ↓
          PG Primary  PG Standby  PG Standby
               ↑                      ↑
              写请求    ←── PgBouncer ──→  读请求
                          (连接池+路由)
```

- Patroni 负责监控 PG 状态、自动故障转移、管理复制
- etcd 存储集群元数据和选主信息
- PgBouncer 做连接池 + 读写分离路由

### 8.4 PG vs MySQL 高可用对比

| 对比 | PostgreSQL | MySQL |
|------|-----------|-------|
| 原生复制 | WAL 流复制 | binlog 复制 |
| 自动故障切换 | 需要 Patroni 等外部工具 | MySQL Group Replication / InnoDB Cluster |
| 连接池 | 必须（PgBouncer） | 可选（ProxySQL） |
| 多主写入 | 不原生支持（BDR 扩展） | Group Replication 支持 |
| 复制延迟 | WAL 物理复制延迟低 | binlog 逻辑复制延迟略高 |

---

## 9. PostgreSQL vs MySQL 深度对比

### 9.1 功能对比

| 特性 | PostgreSQL | MySQL |
|------|-----------|-------|
| SQL 标准 | 支持最完整 | 支持常用子集 |
| DDL 事务 | 支持（CREATE/ALTER 可回滚） | 不支持（DDL 自动提交） |
| JSONB | 原生支持 + GIN 索引 | JSON + 虚拟列索引 |
| 数组 | 原生支持 | 不支持 |
| 全文检索 | 内置 tsvector + GIN | 内置简单 FULLTEXT |
| CTE 递归 | 完整支持 | 8.0+ 支持 |
| 窗口函数 | 完整支持 | 8.0+ 支持 |
| CHECK 约束 | 完整支持 | 8.0.16+ 才真正支持 |
| 物化视图 | 支持 MATERIALIZED VIEW | 不支持 |
| 表继承 | 支持 | 不支持 |
| 自定义类型 | 支持 | 不支持 |
| 扩展 | 插件生态丰富（PostGIS、pg_vector） | 插件机制有限 |

### 9.2 性能对比

| 场景 | PostgreSQL | MySQL |
|------|-----------|-------|
| 简单 CRUD | 略慢（多进程开销） | 快（多线程轻量） |
| 复杂查询 | 优化器更智能，复杂 JOIN 更快 | 优化器简单，复杂查询需手动优化 |
| 写入性能 | 堆表追加快，但 VACUUM 有开销 | 聚簇索引顺序写快，乱序写有页分裂 |
| 并发连接 | 进程模型，需连接池 | 线程模型，支持更多连接 |
| JSON 查询 | JSONB + GIN，查询快 | JSON 无二进制格式，查询慢 |

### 9.3 选型建议

**选 PostgreSQL 的场景**：
- 需要复杂查询（CTE 递归、窗口函数、多表 JOIN）
- 需要 JSONB 半结构化数据存储和查询
- 需要地理信息（PostGIS）
- 需要向量检索（pg_vector，AI 场景）
- 需要 DDL 事务（频繁变更表结构）
- 数据完整性要求高（金融、政府）

**选 MySQL 的场景**：
- 简单 CRUD 为主的 Web 应用
- 团队对 MySQL 更熟悉
- 需要更多的连接数（没有连接池的限制）
- 需要成熟的多主方案（Group Replication）
- 已有 MySQL 生态（中间件、监控）

---

## 10. 常用扩展

| 扩展 | 功能 | 场景 |
|------|------|------|
| PostGIS | 地理空间数据 | 地图、LBS、地理围栏 |
| pg_vector | 向量相似度搜索 | AI 嵌入向量检索、推荐系统 |
| pg_trgm | 三元组相似度 | 模糊搜索（LIKE '%keyword%' 加速） |
| pg_stat_statements | SQL 统计 | 慢查询分析、性能监控 |
| pgcrypto | 加密函数 | 数据加密、密码哈希 |
| postgres_fdw | 外部数据源 | 跨库查询 |
| timescaledb | 时序数据库 | IoT、监控指标 |
| citus | 分布式（分片） | 水平扩展、多租户 |

```sql
-- 启用扩展
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- pg_trgm 模糊搜索
CREATE INDEX idx_name_trgm ON users USING GIN (name gin_trgm_ops);
SELECT * FROM users WHERE name LIKE '%ali%';  -- 走 GIN 索引

-- uuid 生成
SELECT gen_random_uuid();  -- PG 13+ 内置，不需要扩展
```

---

## 11. 最佳实践

### 11.1 连接管理

```
生产环境必须使用连接池：

应用 → PgBouncer（连接池）→ PostgreSQL

PgBouncer 配置建议：
  pool_mode = transaction    -- 事务级别复用（推荐）
  default_pool_size = 20     -- 每个用户的连接数
  max_client_conn = 1000     -- 最大客户端连接数
```

### 11.2 表设计规范

| 规范 | 说明 |
|------|------|
| 主键 | bigserial 或 UUID（`gen_random_uuid()`） |
| 时间字段 | 用 `timestamptz`，不用 `timestamp` |
| 金额 | 用 `numeric` 或 `bigint` 存分 |
| 状态字段 | 用 `smallint` + CHECK 约束，不用 enum |
| 软删除 | `deleted_at timestamptz` + 部分索引 |
| JSON 字段 | 用 `jsonb`，不用 `json` |
| 数组字段 | 适合标签等简单场景，复杂关系还是拆表 |
| 默认值 | 字段尽量 NOT NULL + DEFAULT |

### 11.3 VACUUM 调优

```
-- 关键参数
autovacuum = on                        -- 必须开启
autovacuum_vacuum_scale_factor = 0.1   -- 10% 死元组触发（默认 0.2）
autovacuum_analyze_scale_factor = 0.05 -- 5% 变更触发 ANALYZE
autovacuum_max_workers = 3             -- 并行 worker 数

-- 大表单独调优
ALTER TABLE huge_table SET (autovacuum_vacuum_scale_factor = 0.01);

-- 监控表膨胀
SELECT schemaname, relname, n_dead_tup, n_live_tup,
       round(n_dead_tup::numeric / GREATEST(n_live_tup, 1) * 100, 2) AS dead_pct
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;
```

### 11.4 监控要点

| 监控项 | SQL | 告警阈值 |
|--------|-----|---------|
| 连接数 | `SELECT count(*) FROM pg_stat_activity` | > max_connections × 80% |
| 长事务 | `SELECT * FROM pg_stat_activity WHERE state='active' AND age(now(), xact_start) > '5 min'` | > 5 分钟 |
| 表膨胀 | `SELECT n_dead_tup FROM pg_stat_user_tables` | dead_tup > live_tup × 20% |
| 复制延迟 | `SELECT replay_lag FROM pg_stat_replication` | > 1 秒 |
| 缓存命中率 | `SELECT sum(heap_blks_hit) / sum(heap_blks_hit + heap_blks_read) FROM pg_statio_user_tables` | < 99% |
| 事务 ID 消耗 | `SELECT age(datfrozenxid) FROM pg_database` | > 5 亿 |
