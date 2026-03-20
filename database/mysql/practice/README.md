# MySQL Go 实战 Demo（GORM）

基于 GORM 的 MySQL 企业级实战项目，涵盖生产环境常用的 CRUD、事务、关联查询、分布式 ID、乐观锁、幂等写入等核心特性。

代码按企业级规范编写，可直接复用到生产项目。

---

## 项目结构

```
practice/
├── cmd/server/main.go              # 启动入口：加载配置 → 连接数据库 → 自动建表
│
├── configs/config.yaml              # 配置文件：连接信息、连接池、日志级别
│
├── internal/                        # 业务代码（不对外暴露）
│   ├── config/config.go            #   配置解析：yaml 读取 + 环境变量替换 + 字段校验
│   ├── model/user.go               #   数据模型：User / Order / Product + Hooks
│   ├── dao/                        #   数据访问层（Data Access Object）
│   │   ├── user_dao.go             #     User CRUD + Scopes 可复用查询
│   │   ├── order_dao.go            #     Order CRUD + Preload / Joins 关联查询
│   │   └── product_dao.go          #     Product Upsert 幂等写入 + 乐观锁扣库存
│   └── service/user_service.go     #   业务逻辑层：参数校验、事务转账
│
├── pkg/                             # 可复用公共包
│   ├── database/mysql.go           #   GORM 初始化、连接池配置、日志级别
│   └── snowflake/snowflake.go      #   Twitter 雪花算法：分布式 ID 生成器
│
├── examples/                        # 特性演示（每个特性一个独立可运行文件）
│   ├── setup/setup.go              #   公共初始化（加载配置 + 连接数据库 + 建表）
│   ├── 01_migrate/main.go          #   建库建表最佳实践
│   ├── 02_hooks/main.go            #   GORM Hooks 生命周期
│   ├── 03_scopes/main.go           #   Scopes 可复用查询条件
│   ├── 04_association/main.go      #   Preload vs Joins 关联查询
│   ├── 05_upsert/main.go           #   INSERT ON DUPLICATE KEY UPDATE 幂等写入
│   ├── 06_optimistic_lock/main.go  #   乐观锁并发控制
│   ├── 07_aggregation/main.go      #   聚合统计 + Raw SQL
│   ├── 08_snowflake/main.go        #   雪花算法分布式 ID
│   └── 09_transaction/main.go      #   事务转账 + 悲观锁
│
├── go.mod
└── go.sum
```

**分层架构**：

```
Example / cmd/server
        ↓
    Service（业务逻辑、事务编排）
        ↓
    DAO（数据访问、SQL 操作、Scopes）
        ↓
    Model（数据结构、表映射、Hooks）
        ↓
    pkg/database（连接管理、连接池）
        ↓
    configs/（yaml 配置、环境变量）
```

---

## 环境准备

### 依赖

```
gorm.io/gorm           # ORM 框架
gorm.io/driver/mysql    # MySQL 驱动
gopkg.in/yaml.v3        # YAML 配置解析
```

### 数据库

项目使用 Docker 运行 MySQL 8.0，连接信息：

| 配置项 | 值 |
|--------|-----|
| Host | 127.0.0.1 |
| Port | 13306 |
| User | root |
| Password | root（支持环境变量 `MYSQL_PASSWORD` 覆盖） |
| Database | guide_db |
| phpMyAdmin | http://localhost:18080 |

### 密码安全

`config.yaml` 中密码使用 `${MYSQL_PASSWORD:root}` 语法，支持环境变量覆盖：

```bash
# 开发环境：不设置环境变量，使用默认值 root
go run cmd/server/main.go

# 生产环境：通过环境变量注入真实密码
export MYSQL_PASSWORD=your_real_password
go run cmd/server/main.go
```

---

## 快速开始

```bash
cd practice

# 1. 安装依赖
go mod tidy

# 2. 初始化数据库（建表）
go run cmd/server/main.go

# 3. 运行任意 Example
go run examples/01_migrate/main.go
```

---

## 核心代码说明

### configs/config.yaml — 配置文件

企业级配置，每个参数都有明确的取值依据：

- **连接池参数**：`max_open_conns`（最大连接数）、`max_idle_conns`（空闲连接数）、`conn_max_lifetime`（连接最大存活时间）、`conn_max_idle_time`（空闲连接超时）
- **环境变量替换**：`${ENV_VAR:default}` 语法，敏感信息不硬编码
- **慢查询阈值**：`slow_threshold: 200`（毫秒），超过自动打印日志

### internal/config/config.go — 配置解析

- `Config` / `MySQLConfig` 结构体，字段与 yaml 一一对应
- `LoadConfig(path)` 读取文件 → 环境变量替换 → yaml 解析 → 字段校验
- `resolveEnvVars()` 解析 `${VAR:default}` 格式，逐个替换
- `DSN()` 方法拼接 GORM 连接字符串
- `validate()` 校验 host / port / user / dbname 必填

### internal/model/user.go — 数据模型

定义了三个模型：

**User**（用户）
- 内嵌 `gorm.Model`（自增 ID + CreatedAt + UpdatedAt + DeletedAt 软删除）
- 字段：Name、Age、Email（唯一索引）、Balance（余额，单位分）
- Hooks：`BeforeCreate` / `BeforeUpdate` 自动清洗数据（Email 转小写去空格，Name 去空格）

**Order**（订单）— 双 ID 模式
- `ID uint`：自增主键，保证 B+Tree 顺序写入，内部使用
- `OrderNo int64`：雪花 ID 业务号，全局唯一，对外暴露
- `User` 关联：`belongs to`，通过 `UserID` 外键关联

**Product**（商品）
- `SKU string`：唯一索引，用于 Upsert 冲突判断
- `Version int`：乐观锁版本号，并发更新时检测冲突

### internal/dao/ — 数据访问层

**user_dao.go**

| 方法 | 功能 | 要点 |
|------|------|------|
| `Create` | 插入单条 | 触发 BeforeCreate Hook |
| `BatchCreate` | 批量插入 | `CreateInBatches` 分批写入 |
| `GetByID` | 主键查询 | `First(&user, id)` |
| `GetByEmail` | 唯一索引查询 | `Where("email = ?", email)` |
| `ListByAge` | 分页查询 | `Count` + `Offset` + `Limit`，返回列表+总数 |
| `Update` | 更新指定字段 | 用 `map` 解决 GORM 零值不更新问题 |
| `Delete` | 软删除 | `gorm.Model` 自带 `DeletedAt`，实际执行 UPDATE |
| `WithTx` | 事务支持 | 返回绑定 tx 的新 DAO 实例 |
| `ListWithScopes` | Scopes 组合查询 | 接收可变参数 Scope 函数 |

Scopes（可复用查询条件）：

| Scope | 功能 |
|-------|------|
| `AgeRange(min, max)` | 年龄区间过滤 |
| `BalanceGte(amount)` | 余额大于等于 |
| `Paginate(page, size)` | 通用分页，自带参数保护（page<=0 修正为 1，size 上限 100） |
| `OrderByCreated(desc)` | 按创建时间排序 |

**order_dao.go**

| 方法 | 功能 | 要点 |
|------|------|------|
| `GetByOrderNo` | 按雪花业务号查询 | 对外接口使用，Preload 加载关联用户 |
| `GetByID` | 按自增主键查询 | 内部使用 |
| `ListByUserID` | 查用户的所有订单 | Preload 方式，两次 SQL |
| `ListByUserName` | 按用户名查订单 | Joins 方式，单次 SQL，关联字段做 WHERE |
| `ListPaidWithUser` | 分页查已支付订单 | Scopes + Preload 组合 |

**product_dao.go**

| 方法 | 功能 | 要点 |
|------|------|------|
| `Upsert` | 单条幂等写入 | `ON DUPLICATE KEY UPDATE`，按 SKU 冲突更新 |
| `BatchUpsert` | 批量幂等写入 | `CreateInBatches` + `OnConflict` |
| `DeductStockOptimistic` | 乐观锁扣库存 | `WHERE version = ?` 检测冲突，`RowsAffected == 0` 表示失败 |

### internal/service/user_service.go — 业务逻辑层

| 方法 | 功能 | 要点 |
|------|------|------|
| `CreateUser` | 创建用户 | 参数校验 + 调用 DAO |
| `Transfer` | 转账 | `db.Transaction()` 闭包，FOR UPDATE 悲观锁，`gorm.Expr` 原子更新 |

Transfer 事务流程：
1. `tx.Set("gorm:query_option", "FOR UPDATE").First(&from, fromID)` — 锁住转出方
2. 检查余额是否充足
3. `tx.First(&to, toID)` — 锁住转入方
4. `gorm.Expr("balance - ?", amount)` — 原子扣款
5. `gorm.Expr("balance + ?", amount)` — 原子加款
6. 返回 nil 自动 Commit，返回 error 自动 Rollback

### pkg/database/mysql.go — 数据库连接

- `NewMySQL(cfg)` 根据配置初始化 GORM，设置日志级别、连接池参数，Ping 验证连接
- `Close(db)` 获取底层 `*sql.DB` 并关闭，用于优雅退出
- `parseLogLevel()` 将 yaml 字符串转为 GORM 日志枚举

### pkg/snowflake/snowflake.go — 雪花算法

Twitter Snowflake 算法 Go 实现，生成 64 位全局唯一 ID：

```
64 位结构：1 bit 符号 | 41 bits 时间戳 | 10 bits 机器ID | 12 bits 序列号
```

- **单节点 QPS**：每毫秒 4096 个 ID，约 409.6 万/秒
- **时钟回拨检测**：发现时间倒退则阻塞等待恢复
- **线程安全**：`sync.Mutex` 保护
- `NewNode(nodeID)` 创建节点，nodeID 范围 0~1023
- `Generate()` 生成 ID
- `ParseID(id)` 反解出时间、节点、序列号（调试用）

---

## Examples 详细说明

所有 example 在 `practice/` 目录下运行，每个文件独立可执行。

### 01_migrate — 建库建表最佳实践

```bash
go run examples/01_migrate/main.go
```

演示内容：
- **查看 DDL**：`SHOW CREATE TABLE` 查看 AutoMigrate 生成的建表语句
- **企业建表规范**：字符集 utf8mb4、InnoDB 引擎、字段 NOT NULL + DEFAULT、每列必须有 COMMENT
- **手动创建复合索引**：先查 `information_schema` 判断是否存在，不存在才创建
- **EXPLAIN 分析**：对比主键查询、唯一索引查询、复合索引查询、索引失效（函数运算）的执行计划
- **迁移工具推荐**：golang-migrate、goose，以及迁移文件命名规范

关键知识点：
- 生产环境不用 AutoMigrate，用版本化 SQL 迁移文件
- 复合索引遵循最左前缀原则
- 索引列上做函数运算会导致索引失效

---

### 02_hooks — GORM Hooks 生命周期

```bash
go run examples/02_hooks/main.go
```

演示内容：
- **BeforeCreate**：创建用户时故意传入带空格和大写的 Name/Email，验证 Hook 自动清洗
- **BeforeUpdate**：通过 `db.Save()` 更新时验证 Hook 同样生效
- **不触发的情况**：`db.Updates(map)`、`db.Exec()`、`db.UpdateColumn()` 均不触发 Hook

关键知识点：
- Hook 适合做数据清洗、自动填充，不适合放业务逻辑
- 只有 `db.Create(&model)` 和 `db.Save(&model)` 会触发 Hook
- Map 更新和 Raw SQL 跳过 Hook

---

### 03_scopes — Scopes 可复用查询条件

```bash
go run examples/03_scopes/main.go
```

演示内容：
- **AgeRange**：年龄区间过滤
- **BalanceGte**：余额下限过滤
- **Paginate**：通用分页，内置参数保护（page<=0 修正为 1，size 上限 100）
- **组合查询**：多个 Scope 链式组合，`db.Scopes(scope1, scope2, scope3).Find(&users)`

关键知识点：
- Scope 是 `func(db *gorm.DB) *gorm.DB` 签名的函数，可任意组合
- 企业常见 Scope：分页、排序、状态过滤、时间范围、权限过滤
- Paginate 加 size 上限保护，防止前端传超大值打爆内存

---

### 04_association — Preload vs Joins 关联查询

```bash
go run examples/04_association/main.go
```

演示内容：
- **Preload**：查单个订单自动加载用户信息（两次 SQL）
- **Preload**：查用户的所有订单列表
- **Joins**：按用户名查订单（单次 JOIN SQL，关联字段做 WHERE）
- **组合**：Scopes 分页 + Preload 加载关联

关键知识点：
- Preload 适合一对多、不需要关联字段做 WHERE 的场景，两次 SQL
- Joins 适合一对一、需要关联字段做 WHERE/ORDER BY 的场景，单次 SQL
- 关联数据量大时 Preload（IN 查询）比 Joins 效率更高

---

### 05_upsert — 幂等写入

```bash
go run examples/05_upsert/main.go
```

演示内容：
- **单条 Upsert**：第一次插入，第二次相同 SKU 自动更新价格和库存
- **批量 Upsert**：`CreateInBatches` + `OnConflict` 批量幂等写入

SQL 原理：
```sql
INSERT INTO products (sku, name, price, stock) VALUES (...)
ON DUPLICATE KEY UPDATE name=VALUES(name), price=VALUES(price), stock=VALUES(stock)
```

关键知识点：
- 按唯一键冲突时更新，不重复插入，一条 SQL 搞定
- 典型场景：外部数据同步、消息消费去重、配置批量导入、数据修复脚本
- 不怕重复执行，天然幂等

---

### 06_optimistic_lock — 乐观锁并发控制

```bash
go run examples/06_optimistic_lock/main.go
```

演示内容：
- **正常扣减**：读取 version → 带 version 条件更新 → version+1
- **版本冲突**：模拟两个用户用旧 version 去更新，后者失败
- **带重试的乐观锁**：生产标准写法，失败后重新读取最新 version 重试
- **20 并发压测**：20 个 goroutine 抢 10 个库存，验证不超卖

SQL 原理：
```sql
UPDATE products SET stock = stock - 1, version = version + 1
WHERE id = ? AND version = ? AND stock >= ?
```

关键知识点：
- 乐观锁不加锁，用 version 检测冲突，`RowsAffected == 0` 表示被别人抢先改了
- 适合冲突率低的场景，比悲观锁性能好
- 生产环境必须配合重试机制使用

---

### 07_aggregation — 聚合统计 + Raw SQL

```bash
go run examples/07_aggregation/main.go
```

演示内容：
- **Group By**：按年龄分组统计人数、总余额、平均余额
- **Having**：过滤分组结果（人数 >= 2 的年龄段）
- **Raw SQL**：窗口函数 `RANK() OVER` 做余额排名
- **Raw Exec**：批量更新（年龄 >= 35 的用户余额翻倍）
- **Select 指定列**：只查 id 和 name，减少 IO

关键知识点：
- GORM 的 `Select` + `Group` + `Having` + `Scan` 组合做聚合统计
- 复杂报表、窗口函数等用 `db.Raw()` 直接写 SQL
- 写操作用 `db.Exec()`，返回 `RowsAffected`
- 大表查询只 Select 需要的列，减少网络传输和内存

---

### 08_snowflake — 雪花算法分布式 ID

```bash
go run examples/08_snowflake/main.go
```

演示内容：
- **基本使用**：生成 5 个雪花 ID，反解出时间、节点、序列号
- **双 ID 模式**：订单使用自增主键（内部）+ 雪花业务号（对外）
- **并发安全测试**：10000 个 goroutine 并发生成，验证零重复
- **多节点**：3 个节点各自生成 ID，验证全局不冲突
- **方案对比**：DB 自增 / UUID / 雪花 / Leaf / Redis INCR 的优缺点

关键知识点：
- 雪花 ID 是 64 位整数（BIGINT），趋势递增、全局唯一
- 双 ID 模式：自增主键保证 B+Tree 写入性能，雪花 ID 做业务号对外暴露
- 前端 JavaScript 只有 53 位精度，传给前端需要转成字符串

---

### 09_transaction — 事务转账 + 悲观锁

```bash
go run examples/09_transaction/main.go
```

演示内容：
- **正常转账**：alice → bob 30 元，验证双方余额变化
- **余额不足**：转账失败，验证自动回滚，双方余额不变
- **10 并发转账**：验证 FOR UPDATE 悲观锁保证资金守恒（总和始终不变）
- **要点总结**：Transaction 闭包、FOR UPDATE、gorm.Expr、WithTx 四个核心

关键知识点：
- `db.Transaction(func(tx *gorm.DB) error { ... })` 闭包，nil 提交，error 回滚
- `FOR UPDATE` 悲观锁：锁住行，防止并发读到旧余额
- `gorm.Expr("balance - ?", amount)`：SQL 层面原子更新，避免读-改-写竞态
- `DAO.WithTx(tx)`：事务内创建临时 DAO，保证操作走同一个连接

---

## 运行命令汇总

```bash
cd practice

# 初始化（首次运行）
go mod tidy
go run cmd/server/main.go

# 特性演示
go run examples/01_migrate/main.go          # 建表规范 + EXPLAIN
go run examples/02_hooks/main.go            # Hooks 生命周期
go run examples/03_scopes/main.go           # Scopes 可复用查询
go run examples/04_association/main.go      # Preload vs Joins
go run examples/05_upsert/main.go           # 幂等写入
go run examples/06_optimistic_lock/main.go  # 乐观锁
go run examples/07_aggregation/main.go      # 聚合 + Raw SQL
go run examples/08_snowflake/main.go        # 雪花算法
go run examples/09_transaction/main.go      # 事务转账

# 验证数据
# phpMyAdmin: http://localhost:18080
# 命令行:
mysql -h 127.0.0.1 -P 13306 -u root -proot guide_db -e "SELECT * FROM users"
mysql -h 127.0.0.1 -P 13306 -u root -proot guide_db -e "SELECT * FROM orders"
mysql -h 127.0.0.1 -P 13306 -u root -proot guide_db -e "SELECT * FROM products"
```

---

## 特性覆盖清单

| 特性 | 所在文件 | 生产场景 |
|------|----------|----------|
| CRUD | dao/user_dao.go | 所有业务 |
| 软删除 | model (gorm.Model) | 数据可追溯、防误删 |
| 批量插入 | dao/user_dao.go | 数据导入、批量处理 |
| 分页查询 | dao/user_dao.go | 列表接口 |
| Map 更新 | dao/user_dao.go | 解决零值不更新 |
| Hooks | model/user.go | 数据清洗、自动填充 |
| Scopes | dao/user_dao.go | 可复用查询条件 |
| Preload | dao/order_dao.go | 一对多关联加载 |
| Joins | dao/order_dao.go | 关联字段做 WHERE |
| Upsert | dao/product_dao.go | 幂等写入、数据同步 |
| 乐观锁 | dao/product_dao.go | 低冲突并发控制 |
| 悲观锁 | service/user_service.go | 高冲突并发控制（转账） |
| 事务 | service/user_service.go | 转账、下单等原子操作 |
| gorm.Expr | service/user_service.go | 原子更新避免竞态 |
| Raw SQL | examples/07 | 复杂报表、窗口函数 |
| Group/Having | examples/07 | 聚合统计 |
| Select 指定列 | examples/07 | 大表优化 |
| 雪花算法 | pkg/snowflake | 分布式全局唯一 ID |
| 双 ID 模式 | model/Order | 自增主键+雪花业务号 |
| EXPLAIN | examples/01 | 索引验证 |
| 复合索引 | examples/01 | 查询优化 |
| 连接池配置 | pkg/database | 生产环境调优 |
| 环境变量配置 | internal/config | 敏感信息不硬编码 |
