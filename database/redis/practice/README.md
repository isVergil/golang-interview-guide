# Redis Go 实战 Demo（go-redis）

基于 go-redis/v9 的 Redis 企业级实战项目，涵盖五种数据类型、Pipeline、Lua 脚本、分布式锁、缓存模式、Bitmap/HyperLogLog、Pub/Sub、延迟队列、限流器等核心特性。

代码按企业级规范编写，可直接复用到生产项目。

---

## 项目结构

```
practice/
├── cmd/server/main.go              # 启动入口：加载配置 → 连接 Redis → 验证连接
│
├── configs/config.yaml              # 配置文件：连接信息、连接池、超时配置
│
├── internal/                        # 业务代码（不对外暴露）
│   └── config/config.go            #   配置解析：yaml 读取 + 环境变量替换 + 字段校验
│
├── pkg/                             # 可复用公共包
│   └── redis/redis.go              #   go-redis 初始化、连接池配置、Ping 验证
│
├── examples/                        # 特性演示（每个特性一个独立可运行文件）
│   ├── setup/setup.go              #   公共初始化（加载配置 + 连接 Redis）
│   ├── 01_basic_types/main.go      #   五种基本数据类型（String/Hash/List/Set/ZSet）
│   ├── 02_expire_nx/main.go        #   过期策略 + 条件写入 + Key 管理
│   ├── 03_pipeline/main.go         #   Pipeline 批量操作 + 事务管道
│   ├── 04_lua/main.go              #   Lua 脚本原子操作（CAS 删除/限流/扣库存）
│   ├── 05_distributed_lock/main.go #   分布式锁（SETNX + Watchdog 续期）
│   ├── 06_cache_patterns/main.go   #   缓存模式（穿透/击穿/旁路缓存）
│   ├── 07_bitmap_hyperloglog/main.go #  Bitmap 签到 & HyperLogLog UV 统计
│   ├── 08_pubsub/main.go           #   发布/订阅 + 模式订阅
│   ├── 09_delayed_queue/main.go    #   ZSet 延迟队列
│   └── 10_rate_limiter/main.go     #   限流器（固定窗口/滑动窗口/令牌桶）
│
├── go.mod
└── go.sum
```

**分层架构**：

```
Example / cmd/server
        ↓
    examples/setup（公共初始化）
        ↓
    pkg/redis（连接管理、连接池）
        ↓
    internal/config（yaml 配置、环境变量）
        ↓
    configs/（yaml 配置文件）
```

---

## 环境准备

### 依赖

```
github.com/redis/go-redis/v9    # Redis 客户端
gopkg.in/yaml.v3                # YAML 配置解析
```

### Redis

项目使用 Docker 运行 Redis，连接信息：

| 配置项 | 值 |
|--------|-----|
| Host | 127.0.0.1 |
| Port | 16379 |
| Password | root（支持环境变量 `REDIS_PASSWORD` 覆盖） |
| Pool Size | 100 |

### 连接方式

**方式一：Docker 命令行（redis-cli）**

```bash
# 进入容器直接操作
docker exec -it guide-redis redis-cli -a root

# 常用命令
127.0.0.1:6379> PING              # 测试连通性，返回 PONG
127.0.0.1:6379> SET hello world   # 写入
127.0.0.1:6379> GET hello         # 读取
127.0.0.1:6379> DBSIZE            # 查看 key 数量
127.0.0.1:6379> INFO server       # 查看服务器信息
```

**方式二：RedisInsight 可视化客户端**

浏览器打开 http://localhost:18081，添加数据库时填写：

```
Connection URL: redis://default:root@redis:6379
```

> RedisInsight 和 Redis 都在 Docker 内部，容器间通信用服务名 `redis` + 内部端口 `6379`，不用 `127.0.0.1:16379`。

### 密码安全

`config.yaml` 中密码使用 `${REDIS_PASSWORD:root}` 语法，支持环境变量覆盖：

```bash
# 开发环境：不设置环境变量，使用默认值 root
go run cmd/server/main.go

# 生产环境：通过环境变量注入真实密码
export REDIS_PASSWORD=your_real_password
go run cmd/server/main.go
```

---

## 快速开始

```bash
cd practice

# 1. 安装依赖
go mod tidy

# 2. 验证 Redis 连接
go run cmd/server/main.go

# 3. 运行任意 Example
go run examples/01_basic_types/main.go
```

---

## 核心代码说明

### configs/config.yaml — 配置文件

- **连接池参数**：`pool_size`（最大连接数，默认 100）
- **超时配置**：`dial_timeout`（连接超时）、`read_timeout`（读超时）、`write_timeout`（写超时），单位秒
- **环境变量替换**：`${ENV_VAR:default}` 语法，敏感信息不硬编码

### internal/config/config.go — 配置解析

- `Config` / `RedisConfig` 结构体，字段与 yaml 一一对应
- `LoadConfig(path)` 读取文件 → 环境变量替换 → yaml 解析 → 字段校验
- `resolveEnvVars()` 解析 `${VAR:default}` 格式，逐个替换
- `Addr()` 方法拼接 `host:port` 连接地址
- `validate()` 校验 host / port 必填

### pkg/redis/redis.go — Redis 连接

- `NewRedis(cfg)` 根据配置初始化 go-redis 客户端，设置连接池、超时参数，Ping 验证连接
- `Close(rdb)` 关闭连接，用于优雅退出

---

## Examples 详细说明

所有 example 在 `practice/` 目录下运行，每个文件独立可执行。

### 01_basic_types — 五种基本数据类型

```bash
go run examples/01_basic_types/main.go
```

演示内容：
- **String**：SET/GET、INCR 原子计数、SETNX 不存在才写、MSET/MGET 批量操作
- **Hash**：HSET/HGET、HGETALL 获取全部字段、HINCRBY 字段级原子操作、HDEL/HEXISTS
- **List**：LPUSH/RPUSH 双端插入、LRANGE 范围查询、LPOP/RPOP、LLEN/LINDEX
- **Set**：SADD/SMEMBERS、SISMEMBER 判断存在、SINTER/SUNION/SDIFF 集合运算、SRANDMEMBER 随机抽取
- **ZSet**：ZADD/ZRANGE/ZREVRANGE、ZSCORE/ZRANK/ZREVRANK、ZINCRBY、ZRANGEBYSCORE 范围查询

关键知识点：
- 每种类型都展示了 `OBJECT ENCODING`，观察底层编码转换
- String 的 INCR 是原子的，天然线程安全
- Set 的集合运算适合共同好友、推荐去重等场景
- ZSet 的 Score 支持浮点数，适合排行榜、延迟队列

---

### 02_expire_nx — 过期策略与条件操作

```bash
go run examples/02_expire_nx/main.go
```

演示内容：
- **过期操作**：SET 带 EX、TTL/PTTL 查看剩余时间、EXPIRE 动态设置、PERSIST 取消过期
- **条件写入**：NX（不存在才写）、XX（存在才写）、SetArgs 高级参数
- **Key 管理**：EXISTS 判断存在、TYPE 查看类型、RENAME 重命名、DEL/UNLINK 删除
- **SCAN 遍历**：非阻塞遍历 key（替代 KEYS *）

关键知识点：
- 过期时间加随机抖动（jitter），防止缓存雪崩
- UNLINK 异步删除，不阻塞主线程，适合大 key
- 生产环境禁用 KEYS *，使用 SCAN 分批遍历

---

### 03_pipeline — Pipeline 批量操作

```bash
go run examples/03_pipeline/main.go
```

演示内容：
- **基础 Pipeline**：多条命令打包一次网络往返
- **事务 Pipeline（TxPipeline）**：MULTI/EXEC 包裹，原子执行
- **性能对比**：1000 次单独 SET vs Pipeline SET，展示 RTT 节省

关键知识点：
- Pipeline 减少网络往返，适合批量写入/读取
- TxPipeline = Pipeline + MULTI/EXEC，保证原子性
- Pipeline 不保证原子性，TxPipeline 才保证

---

### 04_lua — Lua 脚本原子操作

```bash
go run examples/04_lua/main.go
```

演示内容：
- **CAS 删除**：GET 并比较值后 DEL（分布式锁释放标准写法）
- **限流器**：INCR + EXPIRE 原子操作，一个窗口内限制请求次数
- **库存扣减**：DECRBY 前检查库存是否充足，不足返回 -1
- **EVALSHA 缓存**：go-redis 自动使用 SHA1 哈希缓存脚本

关键知识点：
- Lua 脚本在 Redis 内原子执行，不会被其他命令打断
- go-redis 的 `redis.NewScript()` 自动管理 EVALSHA → EVAL 回退
- 所有需要 "读-判断-写" 的操作都应该用 Lua 脚本

---

### 05_distributed_lock — 分布式锁

```bash
go run examples/05_distributed_lock/main.go
```

演示内容：
- **基础锁**：SETNX + TTL 加锁，UUID 持有者标识，Lua CAS 释放
- **Watchdog 续期**：后台 goroutine 每 TTL/3 续期一次，防止业务未完成锁已过期
- **并发竞争**：5 个 goroutine 竞争同一把锁，验证互斥性

关键知识点：
- 加锁必须用 SETNX + TTL 原子操作（`SET key value NX EX`）
- 释放锁必须用 Lua CAS（比较持有者 → 删除），防止误删他人的锁
- Watchdog 解决 "业务耗时超过 TTL" 的问题
- 生产环境推荐使用 Redisson（Java）或 redsync（Go）

---

### 06_cache_patterns — 缓存模式

```bash
go run examples/06_cache_patterns/main.go
```

演示内容：
- **Cache Aside（旁路缓存）**：读路径懒加载 + 写路径先更新 DB 再删缓存
- **缓存穿透**：缓存空值方案，短过期时间防污染
- **缓存击穿（互斥锁）**：SETNX 抢锁，只有一个线程查 DB 回填
- **缓存击穿（逻辑过期）**：数据永不过期，逻辑过期时间异步刷新

关键知识点：
- Cache Aside 写路径：先更新 DB → 再删除缓存（不是更新缓存）
- 缓存空值防穿透，但要设短过期防止 DB 有新数据后缓存仍为空
- 互斥锁方案一致性强但有等待开销
- 逻辑过期方案可用性高但有短暂数据不一致

---

### 07_bitmap_hyperloglog — Bitmap 签到 & HyperLogLog UV

```bash
go run examples/07_bitmap_hyperloglog/main.go
```

演示内容：
- **Bitmap 签到**：SETBIT 打卡、GETBIT 查询、BITCOUNT 统计月签到天数、BITPOS 首次签到日
- **连续签到统计**：遍历 Bitmap 位计算最长连续签到
- **BITOP 集合运算**：AND 计算两用户共同签到天数
- **HyperLogLog UV**：PFADD 添加访客、PFCOUNT 统计去重后的 UV
- **PFMERGE**：合并多页面 UV 统计全站 UV
- **内存对比**：100 万元素下 Bitmap vs HyperLogLog 的内存占用

关键知识点：
- Bitmap 每个用户每月仅占 4 字节（31 天），适合精确签到/在线状态
- HyperLogLog 始终约 12KB，适合大基数去重统计（允许 0.81% 误差）
- BITOP 可以做用户活跃度交叉分析

---

### 08_pubsub — 发布/订阅

```bash
go run examples/08_pubsub/main.go
```

演示内容：
- **基础 Pub/Sub**：SUBSCRIBE 订阅频道 → PUBLISH 发布消息 → 消费者接收
- **模式订阅（PSubscribe）**：`order:*` 匹配所有订单状态变更频道
- **多订阅者广播**：一条消息所有订阅者都收到

关键知识点：
- Pub/Sub 是发后即忘模型，消息不持久化，订阅者离线会丢失消息
- 需要可靠消息队列用 Redis Stream（5.0+）或 Kafka
- PSubscribe 适合事件驱动架构（订单状态变更、日志收集）

---

### 09_delayed_queue — ZSet 延迟队列

```bash
go run examples/09_delayed_queue/main.go
```

演示内容：
- **基础延迟队列**：ZADD Score=执行时间戳、ZRANGEBYSCORE 取到期任务、ZREM 原子消费
- **Lua 原子消费**：脚本内 ZRANGEBYSCORE + ZREM，防止多消费者重复消费
- **生产者-消费者模式**：后台轮询 + Lua 取任务，展示完整流程

关键知识点：
- ZSet 的 Score 存执行时间（毫秒时间戳），ZRANGEBYSCORE 取已到期任务
- ZREM 返回值判断是否已被其他消费者消费（乐观锁思想）
- 典型场景：订单超时取消、延迟通知、定时任务调度
- 生产环境可用 Redisson DelayedQueue 或 RocketMQ 延迟消息替代

---

### 10_rate_limiter — 限流器

```bash
go run examples/10_rate_limiter/main.go
```

演示内容：
- **固定窗口**：INCR + EXPIRE，窗口内计数，超限拒绝
- **滑动窗口**：ZSet 存每次请求时间戳，ZREMRANGEBYSCORE 清理过期、ZCARD 判断限额
- **令牌桶**：Hash 存上次时间和令牌数，按时间补充令牌，请求消耗令牌
- **并发压测**：50 个 goroutine 并发请求，验证限流器正确性

关键知识点：
- 固定窗口实现简单，但有边界突发问题（两个窗口交界处可通过 2 倍流量）
- 滑动窗口精确但内存占用随请求量线性增长
- 令牌桶支持突发流量（桶满时可瞬间消费），适合 API 网关
- 所有限流逻辑用 Lua 脚本保证原子性

---

## 运行命令汇总

```bash
cd practice

# 初始化（首次运行）
go mod tidy
go run cmd/server/main.go

# 特性演示
go run examples/01_basic_types/main.go       # 五种数据类型
go run examples/02_expire_nx/main.go         # 过期策略 + 条件操作
go run examples/03_pipeline/main.go          # Pipeline 批量操作
go run examples/04_lua/main.go               # Lua 脚本原子操作
go run examples/05_distributed_lock/main.go  # 分布式锁
go run examples/06_cache_patterns/main.go    # 缓存模式
go run examples/07_bitmap_hyperloglog/main.go # Bitmap + HyperLogLog
go run examples/08_pubsub/main.go            # 发布/订阅
go run examples/09_delayed_queue/main.go     # 延迟队列
go run examples/10_rate_limiter/main.go      # 限流器

# 验证数据
redis-cli -h 127.0.0.1 -p 16379 -a root
```

---

## 特性覆盖清单

| 特性 | 所在文件 | 生产场景 |
|------|----------|----------|
| String CRUD | 01_basic_types | 缓存、计数器、分布式锁 |
| Hash CRUD | 01_basic_types | 对象缓存、用户信息 |
| List 操作 | 01_basic_types | 消息队列、最新列表 |
| Set 集合运算 | 01_basic_types | 共同好友、标签去重 |
| ZSet 排序 | 01_basic_types | 排行榜、延迟队列 |
| 过期策略 | 02_expire_nx | 缓存自动失效 |
| 条件写入 NX/XX | 02_expire_nx | 分布式锁、防重复写 |
| SCAN 遍历 | 02_expire_nx | 批量 key 处理 |
| Pipeline 批量 | 03_pipeline | 批量写入、减少 RTT |
| TxPipeline 事务 | 03_pipeline | 原子批量操作 |
| Lua CAS 删除 | 04_lua | 分布式锁安全释放 |
| Lua 限流 | 04_lua | API 限流 |
| Lua 库存扣减 | 04_lua | 秒杀库存 |
| EVALSHA 缓存 | 04_lua | 脚本性能优化 |
| SETNX 分布式锁 | 05_distributed_lock | 分布式互斥 |
| Watchdog 续期 | 05_distributed_lock | 长任务锁保活 |
| Cache Aside | 06_cache_patterns | 标准缓存读写模式 |
| 缓存穿透防护 | 06_cache_patterns | 恶意请求防护 |
| 缓存击穿防护 | 06_cache_patterns | 热点 key 保护 |
| Bitmap 签到 | 07_bitmap_hyperloglog | 用户签到、在线状态 |
| BITOP 运算 | 07_bitmap_hyperloglog | 活跃度交叉分析 |
| HyperLogLog UV | 07_bitmap_hyperloglog | 大基数去重统计 |
| Pub/Sub | 08_pubsub | 事件通知、广播 |
| 模式订阅 | 08_pubsub | 事件驱动架构 |
| ZSet 延迟队列 | 09_delayed_queue | 订单超时、定时任务 |
| Lua 原子消费 | 09_delayed_queue | 多消费者防重复 |
| 固定窗口限流 | 10_rate_limiter | 简单 QPS 限制 |
| 滑动窗口限流 | 10_rate_limiter | 精确流量控制 |
| 令牌桶限流 | 10_rate_limiter | API 网关限流 |
| 连接池配置 | pkg/redis | 生产环境调优 |
| 环境变量配置 | internal/config | 敏感信息不硬编码 |
