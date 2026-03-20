# MySQL 面试高频题

> 以下回答按照面试中口语化表达整理，重点在于思路清晰、逻辑连贯，面试官听得舒服。

---

## 一、架构与原理

### 1. 一条 SELECT 语句在 MySQL 里是怎么执行的？

大致分五步走：

**连接器 → 解析器 → 优化器 → 执行器 → 存储引擎**

先说连接器，客户端连进来的时候，连接器负责做身份验证、权限校验这些事。连上以后你的这条 SQL 就交给解析器了。

解析器干的事就是词法分析和语法分析，简单说就是看你这条 SQL 写得对不对、是 SELECT 还是 UPDATE、查的哪张表哪些字段。如果语法有问题，这一步就直接报错了。

然后到优化器，这一步很关键。同样一条 SQL 可能有很多种执行方案，比如走哪个索引、先查哪张表，优化器会选一个它认为成本最低的方案出来。

最后执行器拿着优化器给的执行计划，去调存储引擎的接口，一行一行地读数据、做过滤，把结果返回给客户端。

补充一点，MySQL 8.0 之前还有个查询缓存，但是命中率太低了，8.0 直接给删掉了。

### 2. 那 UPDATE 语句呢？跟 SELECT 有什么不一样？

UPDATE 前面那几步是一样的，但到存储引擎层的时候多了很多事情。我按顺序说一下：

1. 先把要改的那一行从磁盘读到 Buffer Pool 里（如果已经在内存里就不用读了）
2. 写 Undo Log，记录旧值，方便回滚和 MVCC 用
3. 在 Buffer Pool 里把数据改了，这时候这个页就变成脏页了
4. 写 Redo Log，先标记为 prepare 状态
5. 写 Binlog
6. 最后把 Redo Log 标记为 commit

为什么要搞两阶段提交呢？就是为了保证 Redo Log 和 Binlog 的一致性。你想啊，如果 Redo 写了 Binlog 没写，主库崩了恢复以后有这条数据，但从库没收到 Binlog 就少了这条数据，主从就不一致了。

### 3. InnoDB 和 MyISAM 有什么区别？

面试最常问的几个点：

**事务**：InnoDB 支持事务，MyISAM 不支持。这是最核心的区别。

**锁粒度**：InnoDB 支持行锁，MyISAM 只有表锁。并发一上来 MyISAM 就扛不住了。

**外键**：InnoDB 支持，MyISAM 不支持。

**崩溃恢复**：InnoDB 有 Redo Log，崩了能恢复到一致状态。MyISAM 没有，崩了可能数据就坏了。

**索引区别**：InnoDB 用聚簇索引，数据和主键索引存在一起。MyISAM 用非聚簇索引，数据和索引是分开的文件。

所以现在基本上都用 InnoDB，MyISAM 几乎不用了。MySQL 5.5 之后默认引擎就是 InnoDB。

### 4. Buffer Pool 了解吗？LRU 是怎么改进的？

Buffer Pool 是 InnoDB 的内存缓存区，用来缓存数据页和索引页，减少磁盘 IO。

普通 LRU 有个问题：如果你做了一次全表扫描，大量冷数据会涌进来把热数据挤出去，这叫**缓存污染**。

InnoDB 的改进 LRU 把链表分成了两段：**young 区（5/8）**和 **old 区（3/8）**。

工作原理是这样的：
- 新读进来的页先放到 old 区的头部，不会直接进 young 区
- 只有在 old 区待了超过 1 秒，再次被访问时，才会被提升到 young 区
- 全表扫描的数据虽然会进 old 区，但因为扫完就不会再访问了，待不到 1 秒就被淘汰了

这样热数据就安全了，不会被一次全表扫描冲掉。

---

## 二、索引

### 5. 为什么 MySQL 用 B+Tree 做索引，不用 B-Tree 或者 Hash？

先说 B+Tree 相比 B-Tree 的优势：

B-Tree 所有节点都存数据，B+Tree 只有叶子节点存数据。这就意味着 B+Tree 的非叶子节点能塞更多的 key 进去，同样的数据量树会更矮。树更矮意味着查一次数据的磁盘 IO 次数更少。你想一下，三层的 B+Tree 大概能存两千多万条数据，绝大多数查询 3 次 IO 就够了。

第二个优势是范围查询。B+Tree 的叶子节点之间用双向链表串起来了，做范围查询比如 `WHERE age BETWEEN 20 AND 30` 的时候，找到起点直接沿着链表遍历就行，非常快。B-Tree 就得做中序遍历，效率差很多。

再说 Hash，Hash 索引等值查询确实 O(1) 很快，但它不支持范围查询，不支持排序，也不支持最左前缀匹配。实际业务里范围查询太常见了，所以不适合做通用索引。

### 6. 聚簇索引和二级索引有什么区别？什么是回表？什么是覆盖索引？

**聚簇索引**就是主键索引，它的叶子节点存的是完整的行数据。每张表只能有一个聚簇索引，因为数据只能按一种方式物理排序。

**二级索引**（也叫辅助索引、非聚簇索引）叶子节点存的不是完整行，而是主键值。可以建很多个。

**回表**就是说，你用二级索引查到了主键值，但你 SELECT 的列不全在索引里，就得拿着主键值再去聚簇索引查一次完整的行数据。这个过程就叫回表。

**覆盖索引**就是反过来，你查的列全都在索引里了，不需要回表。EXPLAIN 的时候 Extra 列会显示 `Using index`。

举个例子，有个联合索引 `(name, age)`：
- `SELECT name, age FROM users WHERE name = 'alice'` 这个就是覆盖索引，name 和 age 都在索引里
- `SELECT * FROM users WHERE name = 'alice'` 这个就得回表，因为 `SELECT *` 肯定有索引里没有的列

生产环境里尽量利用覆盖索引，能少一次回表查询，性能差距很大。

### 7. 联合索引的最左前缀原则是什么？

假设有个联合索引 `(a, b, c)`，它在 B+Tree 里是按 a 排序，a 相同按 b 排序，b 相同按 c 排序。

所以查询的时候：
- `WHERE a = 1`：能用索引
- `WHERE a = 1 AND b = 2`：能用索引
- `WHERE a = 1 AND b = 2 AND c = 3`：完美用上全部
- `WHERE b = 2`：用不了，因为 b 是在 a 确定之后才有序的，单独查 b 没有排序可言
- `WHERE a = 1 AND c = 3`：只能用到 a，c 用不上，因为中间跳过了 b

还有一个点，`WHERE a = 1 AND b > 2 AND c = 3`，a 和 b 能用上索引，但 c 用不上。因为 b 是范围查询，b 确定范围之后 c 就不是有序的了。

### 8. 索引失效的场景有哪些？

这个面试很爱问，我列几个常见的：

**在索引列上做函数运算**：`WHERE YEAR(create_time) = 2024`，优化器没法用索引了，得全表扫描。正确写法是 `WHERE create_time >= '2024-01-01' AND create_time < '2025-01-01'`。

**在索引列上做计算**：`WHERE age + 1 = 20`，应该写成 `WHERE age = 19`。

**隐式类型转换**：字段是 varchar 你传了个数字 `WHERE phone = 13800138000`，MySQL 会把 phone 列转成数字来比较，相当于对索引列做了函数运算，索引就失效了。

**左模糊查询**：`WHERE name LIKE '%abc'`，B+Tree 是按前缀排序的，你前面都是通配符没法走索引。`LIKE 'abc%'` 可以。

**OR 条件**：`WHERE a = 1 OR b = 2`，如果 b 没有索引，整个查询就得全表扫描。

**不满足最左前缀**：联合索引 `(a, b, c)` 你直接查 `WHERE b = 1`，用不上。

实际遇到问题就用 EXPLAIN 看一下，关注 type 列是不是 ALL（全表扫描），key 列有没有用上索引。

### 9. EXPLAIN 怎么看？重点关注哪些字段？

最核心的几个字段：

**type**：查询类型，性能从好到差大概是 `system > const > eq_ref > ref > range > index > ALL`。看到 ALL 就是全表扫描，得优化了。

**key**：实际用了哪个索引。NULL 就是没用索引。

**rows**：预估要扫描多少行。这个值越小越好。

**Extra**：
- `Using index` → 覆盖索引，不用回表，好事
- `Using where` → 在存储引擎取数据之后还要在 Server 层做过滤
- `Using filesort` → 用了文件排序，没用上索引的排序，性能差
- `Using temporary` → 用了临时表，通常出现在 GROUP BY 没有合适索引的时候

**key_len**：判断联合索引用了几个列。计算方式是 int 4 字节、bigint 8 字节、varchar(n) 是 n × 字符集字节数 + 2（长度标识），允许 NULL 的再 +1。

---

## 三、事务

### 10. 事务的 ACID 怎么理解？InnoDB 怎么实现的？

**A（原子性）**：事务要么全做要么全不做。靠 Undo Log 实现的，出错了就用 Undo Log 回滚。

**C（一致性）**：这个是最终目标，AID 三个都是为了保证 C。比如转账，不管成功还是失败，两个人的钱加起来应该不变。

**I（隔离性）**：多个事务并发的时候互相不影响。靠 MVCC + 锁 实现的。

**D（持久性）**：事务提交了数据就不会丢。靠 Redo Log 实现的。就算数据库崩了，重启后用 Redo Log 恢复。

### 11. 四个隔离级别分别解决什么问题？

从低到高：

**READ UNCOMMITTED（读未提交）**：啥都不管，能读到别人还没提交的数据（脏读）。基本没人用。

**READ COMMITTED（读已提交）**：只能读到别人已提交的数据，解决了脏读。但同一个事务里两次读同一行结果可能不一样（不可重复读），因为中间别人改了并提交了。Oracle 默认就是这个级别。

**REPEATABLE READ（可重复读）**：InnoDB 默认级别。同一个事务里多次读同一行结果一致，解决了不可重复读。而且 InnoDB 在这个级别下还通过 MVCC + Next-Key Lock 解决了幻读。

**SERIALIZABLE（串行化）**：完全串行执行，所有问题都解决了但性能最差，基本不用。

### 12. MVCC 是什么？怎么实现的？

MVCC 就是多版本并发控制。核心思想是读不加锁，写不阻塞读，通过版本链来实现快照读。

实现靠三个东西：

**第一个是隐藏列**。每一行数据都有两个隐藏字段：`trx_id`（最后修改这行的事务 ID）和 `roll_pointer`（指向 Undo Log 中这行的上一个版本）。

**第二个是 Undo Log 版本链**。每次修改一行数据，旧版本就通过 roll_pointer 串成一个链表。比如一行数据被改了三次，就有三个版本串在一起。

**第三个是 ReadView**。事务做快照读的时候会创建一个 ReadView，里面记录了当前活跃事务的 ID 列表。然后沿着版本链去找，找到第一个对当前事务可见的版本。

判断可见性的规则：
- 如果这个版本的 trx_id 就是我自己，那肯定可见
- 如果 trx_id 比 ReadView 里最小的活跃事务 ID 还小，说明改这行的事务早就提交了，可见
- 如果 trx_id 大于等于 ReadView 创建时的下一个事务 ID，说明是之后才开始的事务，不可见
- 如果 trx_id 在活跃事务列表里，说明那个事务还没提交，不可见

**RC 和 RR 的区别**就是 ReadView 的创建时机不同。RC 每次 SELECT 都创建新的 ReadView，所以能看到别人新提交的数据。RR 只在第一次 SELECT 时创建，后面复用，所以整个事务看到的都是同一个快照。

### 13. InnoDB 是怎么解决幻读的？

先说什么是幻读：同一个事务里，两次相同的范围查询，第二次多出了新行（别的事务插入的）。

InnoDB 在 RR 级别下分两种情况解决：

**快照读（普通 SELECT）**：靠 MVCC 解决。因为 RR 下 ReadView 是第一次 SELECT 创建后复用的，后面别人插入的新行 trx_id 肯定大于 ReadView 的 max_trx_id，按可见性规则判定为不可见，所以读不到新行。

**当前读（SELECT FOR UPDATE / UPDATE / DELETE）**：靠 Next-Key Lock 解决。当前读会加锁，不是只锁已有的行，而是用 Next-Key Lock 锁住一个范围（左开右闭区间）。别的事务想在这个范围内插入新行，会被 Gap Lock 阻塞住，插不进来，所以也不会出现幻行。

但说实话，InnoDB 在 RR 下解决幻读不是百分百的。有一种边界情况：事务 A 先快照读，然后事务 B 插入新行并提交，事务 A 再做当前读（比如 UPDATE），这时候就能看到 B 插入的行。因为快照读不加锁，不会阻止 B 插入。

---

## 四、锁

### 14. MySQL 有哪些锁？

先按**粒度**分：

- **表锁**：锁整张表，开销小但并发低。MyISAM 就是表锁。
- **行锁**：锁具体的行，开销大但并发高。InnoDB 支持行锁。

再按**模式**分：

- **共享锁（S锁）**：读锁，多个事务可以同时加 S 锁。`SELECT ... LOCK IN SHARE MODE`
- **排他锁（X锁）**：写锁，只有一个事务能加。`SELECT ... FOR UPDATE` 或者 UPDATE/DELETE 自动加。

InnoDB 还有三种**行级锁的实现**：

- **Record Lock**：锁住一条记录
- **Gap Lock**：锁住两条记录之间的间隙，防止插入。这是为了解决幻读
- **Next-Key Lock**：Record Lock + Gap Lock 的组合，锁住记录本身加上它前面的间隙，左开右闭区间 `(a, b]`

还有**意向锁（IS/IX）**：这是表级别的标记锁。事务要加行级 S 锁之前先加表级 IS 锁，加行级 X 锁之前先加表级 IX 锁。作用是让加表锁的时候不用逐行检查有没有行锁，直接看意向锁就知道了。

### 15. Next-Key Lock 的加锁规则能详细说说吗？

加锁的基本单位是 Next-Key Lock，左开右闭区间。但在特定条件下会退化：

**等值查询唯一索引**：找到记录就退化为 Record Lock，只锁那一行。因为唯一索引保证不会有重复值，不需要锁间隙。找不到记录就退化为 Gap Lock。

**等值查询非唯一索引**：命中后向右遍历到第一个不满足条件的值，这个 Next-Key Lock 退化为 Gap Lock（不锁那个不满足条件的记录本身）。

**范围查询**：一般不退化，就是 Next-Key Lock。但 MySQL 8.0.18 之后，范围查询唯一索引在边界值存在的情况下也有些优化。

举个例子，表里 id 是主键，有 5、10、15、20、25 这几条记录：

`SELECT * FROM t WHERE id = 12 FOR UPDATE`：12 不存在，落在 (10, 15) 区间，加 Gap Lock (10, 15)，防止别人在这个区间插入。

`SELECT * FROM t WHERE id = 15 FOR UPDATE`：15 存在且是唯一索引，退化为 Record Lock，只锁 id=15 这一行。

`SELECT * FROM t WHERE id >= 10 AND id < 12 FOR UPDATE`：锁 [10, 15) 这个范围的 Next-Key Lock。

### 16. 什么是死锁？怎么处理？

死锁就是两个事务互相等待对方释放锁，谁也进行不下去。

经典场景：事务 A 锁了 id=1 要去锁 id=2，事务 B 锁了 id=2 要去锁 id=1。互相等着，卡住了。

InnoDB 有两种处理方式：

**超时机制**：等锁等超过 `innodb_lock_wait_timeout`（默认 50 秒）就放弃。但 50 秒太长了，线上一般会调到 5~10 秒。

**死锁检测**：默认开启 `innodb_deadlock_detect = ON`。InnoDB 会主动检测锁等待的依赖关系，发现环了就选一个代价最小的事务回滚掉。这个更及时。

**怎么预防**：
- 不同的事务访问表和行的顺序保持一致，比如都先锁 id 小的再锁 id 大的
- 缩短事务，锁持有时间越短死锁概率越低
- 合理建索引，没有索引的话行锁会升级为表锁，更容易死锁

### 17. 乐观锁和悲观锁有什么区别？分别什么场景用？

这两个不是 MySQL 原生的锁类型，是并发控制的两种策略。

**悲观锁**：假设一定会有冲突，操作前先加锁。MySQL 里就是 `SELECT ... FOR UPDATE`，锁住这行，别人改不了。

适合冲突率高的场景，比如转账扣款——余额这种东西并发改的概率很高，不锁住可能超扣。

**乐观锁**：假设一般不会冲突，不加锁。给表加个 version 字段，更新的时候带上 `WHERE version = 旧版本`，如果 RowsAffected 是 0 说明别人先改了，冲突了就重试。

适合冲突率低的场景，比如更新商品信息、修改用户资料。大部分时候不会冲突，偶尔冲突重试一下就行。

**关键区别**：悲观锁真的会阻塞别人等你释放锁，有性能开销。乐观锁不阻塞任何人，冲突的时候才有额外开销（重试）。

---

## 五、日志

### 18. Redo Log、Undo Log、Binlog 分别是什么？

**Redo Log**：InnoDB 引擎层面的，物理日志，记录的是"某个数据页的某个偏移量写了什么值"。作用是**崩溃恢复**，保证持久性。事务提交时数据可能还在 Buffer Pool 没刷盘，但 Redo Log 先写了，崩了用 Redo Log 恢复就行。这就是 WAL（Write-Ahead Logging）。

**Undo Log**：也是 InnoDB 层面的，逻辑日志，记录的是反向操作，比如 INSERT 就记 DELETE，UPDATE 就记旧值。两个作用：**回滚**（事务失败了用 Undo Log 撤销）和 **MVCC**（版本链就是 Undo Log 串起来的）。

**Binlog**：Server 层面的，逻辑日志，记录的是 SQL 语句级别的变更。两个作用：**主从复制**（从库拉 Binlog 重放）和**数据恢复**（基于全量备份 + Binlog 恢复到任意时间点）。

面试官经常会问 Redo Log 和 Binlog 的区别：
- 层级不同：Redo 是引擎层，Binlog 是 Server 层
- 类型不同：Redo 是物理日志，Binlog 是逻辑日志
- 写入方式不同：Redo 是循环写固定大小，写满了就覆盖。Binlog 是追加写，写满一个文件开新的
- 用途不同：Redo 负责崩溃恢复，Binlog 负责复制和归档

### 19. 两阶段提交是什么？为什么需要？

两阶段提交是 Redo Log 和 Binlog 之间的协调机制：

1. 写 Redo Log，标记为 **prepare** 状态
2. 写 Binlog
3. 把 Redo Log 标记为 **commit** 状态

为什么需要？假设不做两阶段提交：

**先写 Redo 再写 Binlog**：Redo 写完 Binlog 没写就崩了。主库用 Redo 恢复，数据有了。但从库没收到 Binlog，数据没了。主从不一致。

**先写 Binlog 再写 Redo**：Binlog 写完 Redo 没写就崩了。从库收到 Binlog 有了这条数据。主库没 Redo 恢复不了，数据没了。还是不一致。

两阶段提交的恢复规则：
- Redo 是 prepare 但 Binlog 完整 → 提交（两边都有，安全）
- Redo 是 prepare 但 Binlog 不完整 → 回滚（Binlog 没写成功，丢弃）

这样就保证了主从数据一致。

---

## 六、高可用与复制

### 20. MySQL 主从复制的原理是什么？

三个线程配合：

1. **主库 Binlog Dump Thread**：主库有数据变更就把 Binlog 发给从库
2. **从库 I/O Thread**：接收主库发来的 Binlog，写到本地的 Relay Log（中继日志）
3. **从库 SQL Thread**：读 Relay Log 里的事件，回放执行，数据就同步过来了

默认是异步复制，主库写完 Binlog 就返回客户端成功了，不管从库有没有收到。好处是性能好，坏处是主库挂了可能有数据丢失。

半同步复制是改进方案，主库要等至少一个从库确认收到 Binlog 才返回成功，安全性更好但会增加一点延迟。

### 21. 主从延迟怎么解决？

主从延迟说白了就是从库的 SQL Thread 回放速度跟不上主库的写入速度。

常见原因：
- 从库单线程回放（老版本），主库是并发写的，从库串行回放当然慢
- 大事务，比如一个删除几百万行的 DELETE，从库也得执行那么久
- 从库机器配置差
- 网络延迟

解决方案：

**并行复制**：MySQL 5.7 引入了基于 LOGICAL_CLOCK 的并行复制，同一组提交的事务可以在从库并行回放。8.0 引入了 writeset 并行复制，效果更好。

**业务层面**：写后立即读的场景直接走主库查。不紧急的读走从库。

**GTID 方案**：`WAIT_FOR_EXECUTED_GTID_SET`，写操作返回 GTID，读操作带着 GTID 去从库查，从库如果还没回放到这个 GTID 就等一下，确保读到最新数据。

---

## 七、性能优化

### 22. 慢查询怎么排查和优化？

排查流程：

**第一步**：开慢查询日志。`slow_query_log = ON`，设一个阈值 `long_query_time = 1`（超过 1 秒就记录）。

**第二步**：用 `mysqldumpslow` 或者 `pt-query-digest` 分析慢查询日志，找出 Top N 的慢 SQL。

**第三步**：EXPLAIN 看执行计划。重点看 type 有没有全表扫描、key 有没有用上索引、rows 扫了多少行、Extra 有没有 filesort 或 temporary。

**第四步**：针对性优化。

常见优化手段：
- 加索引或调整索引（联合索引顺序、覆盖索引）
- 改 SQL 写法（避免索引失效的那些情况）
- 减少返回列（`SELECT *` 改成只查需要的列）
- 分页优化（深分页用游标或延迟关联）
- 大表拆分查询（一次删 100 万行改成分批每次 1000 行）

### 23. 深分页怎么优化？

`SELECT * FROM orders ORDER BY id LIMIT 1000000, 10` 这种查询，MySQL 要先扫描 100 万零 10 行，丢掉前 100 万行，只返回最后 10 行。越往后翻页越慢。

**方案一：游标分页**

```sql
SELECT * FROM orders WHERE id > #{上一页最后一条的id} ORDER BY id LIMIT 10
```

用上一页最后一条记录的 id 作为起点，直接定位到那个位置开始取 10 条。走主键索引，不管翻到第几页都是只扫 10 行。

缺点是只能上一页下一页，不能跳页。

**方案二：延迟关联**

```sql
SELECT * FROM orders
INNER JOIN (SELECT id FROM orders ORDER BY id LIMIT 1000000, 10) t
ON orders.id = t.id
```

子查询只查 id，走的是覆盖索引，不用回表，快很多。拿到 10 个 id 之后再用 JOIN 去取完整数据。

### 24. 连接池参数怎么配？

生产环境几个关键参数：

**max_open_conns（最大连接数）**：根据业务 QPS 和查询耗时估算。比如 QPS 1000，平均查询 10ms，理论上 10 个连接就够。一般设 50~200，太大了连接切换开销反而高。

**max_idle_conns（最大空闲连接数）**：保持一些空闲连接避免频繁创建销毁。一般是 max_open 的 10%~25%。

**conn_max_lifetime（连接最大存活时间）**：建议小于 MySQL 的 `wait_timeout`（默认 8 小时），设 1 小时比较常见。防止用到被 MySQL 已经断开的僵尸连接。

**conn_max_idle_time（空闲连接超时）**：空闲多久就关掉。设 10 分钟左右，回收长期闲置的连接。

---

## 八、分库分表

### 25. 什么时候需要分库分表？

**分库**：单库连接数到瓶颈了。MySQL 单实例默认最大连接数 151，调大了也有上限。多个服务都连同一个库，连接数不够用了就得分库。

**分表**：单表数据量太大了。一般来说 B+Tree 三层结构能装大概两千万行。超过之后索引树会变高，查询性能开始下降。阿里巴巴开发规范建议单表不超过 500 万行，这个比较保守但也是个参考值。

能不分就不分。先考虑加索引、优化 SQL、读写分离、缓存这些手段，都搞不定了再考虑分库分表。

### 26. 分库分表有哪些策略？

**垂直拆分**：
- 垂直分库：按业务拆，用户库、订单库、商品库各自独立
- 垂直分表：把一张宽表的列拆开，常用字段一张表，大字段（TEXT/BLOB）一张表

**水平拆分**：
- 水平分库：同一张表的数据按规则分散到多个库
- 水平分表：同一个库里把表拆成多个（user_0, user_1, user_2...）

分片路由策略：
- **Hash 取模**：`user_id % 4` 分到 4 张表。数据均匀但扩容麻烦（要迁数据）
- **Range 范围**：`id 1~1000万 → 表1，1000万~2000万 → 表2`。方便扩容但可能热点集中
- **一致性哈希**：比 Hash 取模扩容友好，只需迁移少量数据

### 27. 分库分表之后有哪些问题？

**分布式 ID**：自增 ID 在多库下会重复。用雪花算法、美团 Leaf、或者号段模式生成全局唯一 ID。

**跨库 JOIN**：分到不同库的表没法 JOIN 了。解决方案是在应用层做聚合，或者做数据冗余（把需要 JOIN 的字段冗余存一份）。

**分布式事务**：跨库操作没法用本地事务了。用 XA（两阶段提交）、TCC（Try-Confirm-Cancel）或 Saga 模式。但都比较复杂，能避免尽量避免。

**跨分片分页**：`ORDER BY + LIMIT` 在分库场景下，得从每个分片都查出来然后在应用层归并排序。深分页性能很差。

**全局唯一约束**：唯一索引只能保证单库唯一，跨库的唯一性需要额外处理。

中间件选择：ShardingSphere（Java 生态主流）、Vitess（Go 生态、YouTube 在用）、TiDB（直接用分布式数据库，不用手动分库分表）。

---

## 九、实战场景题

### 28. 转账怎么保证数据一致性？

用事务 + 悲观锁。代码思路是这样的：

```
开启事务
  SELECT * FROM users WHERE id = fromID FOR UPDATE   -- 锁住转出方
  检查余额够不够
  SELECT * FROM users WHERE id = toID FOR UPDATE     -- 锁住转入方
  UPDATE users SET balance = balance - amount WHERE id = fromID  -- 原子扣款
  UPDATE users SET balance = balance + amount WHERE id = toID    -- 原子加款
提交事务（任何一步失败自动回滚）
```

几个关键点：
- `FOR UPDATE` 悲观锁锁住行，防止并发读到旧余额
- `balance = balance - amount` 用 SQL 表达式做原子更新，不要在应用层读出来减了再写回去（读-改-写有竞态）
- 事务内的所有操作走同一个数据库连接
- 余额用整数存分，不要用浮点数

### 29. 库存扣减怎么防止超卖？

两种方案：

**悲观锁方案**：`SELECT stock FROM products WHERE id = ? FOR UPDATE`，锁住这行，读出库存判断够不够，够就扣减。简单直接但并发性能差，高并发场景锁竞争严重。

**乐观锁方案**：给商品加个 version 字段。

```sql
UPDATE products SET stock = stock - 1, version = version + 1
WHERE id = ? AND version = ? AND stock >= 1
```

不加锁，用版本号检测冲突。RowsAffected 为 0 就说明被别人抢先了，重新读最新 version 再试。配合重试机制，一般重试 3~5 次。

冲突率低用乐观锁，冲突率高用悲观锁。秒杀这种场景其实更适合用 Redis 原子操作做预扣减，MySQL 兜底。

### 30. 幂等写入怎么做？

用 MySQL 的 `INSERT ... ON DUPLICATE KEY UPDATE` 语法。

```sql
INSERT INTO products (sku, name, price, stock) VALUES ('BOOK-001', 'Go圣经', 6800, 100)
ON DUPLICATE KEY UPDATE name = VALUES(name), price = VALUES(price), stock = VALUES(stock)
```

按唯一键 sku 判断：不存在就插入，存在就更新指定列。一条 SQL 搞定，不需要先查再判断。

典型场景：外部数据同步（定时拉取第三方数据）、消息消费去重（MQ 消费者重试）、批量数据导入。

GORM 里对应的写法是 `Clauses(clause.OnConflict{...})`。

### 31. 分布式 ID 方案怎么选？

常见方案对比：

**数据库自增**：最简单，但单点瓶颈，分库后 ID 重复。

**UUID**：不依赖任何服务，但 36 个字符太长、完全无序导致 B+Tree 频繁页分裂、索引性能差。

**雪花算法（Snowflake）**：64 位整数，趋势递增、全局唯一、高性能（单节点 400 万/秒）。缺点是依赖时钟（时钟回拨会出问题）、需要分配机器 ID。

**美团 Leaf**：支持号段模式和雪花模式，高可用。但需要额外部署服务。

**Redis INCR**：利用 Redis 的原子自增。简单有序，但依赖 Redis，有持久化风险。

我们项目里用的是**双 ID 模式**：数据库主键用自增（保证 B+Tree 写入性能），业务号用雪花 ID（全局唯一，对外暴露）。前端收到的是雪花 ID 字符串（因为 JavaScript 的 Number 只有 53 位精度，超过会丢精度）。

### 32. 大事务有什么危害？怎么避免？

危害很多：

- **锁持有时间长**：事务不提交锁就不释放，阻塞其他事务
- **Undo Log 膨胀**：长事务导致版本链很长，MVCC 查询时要沿着链表往回找，性能下降
- **主从延迟**：大事务的 Binlog 很大，从库回放也很慢
- **回滚代价大**：事务失败了要回滚，改了 100 万行就要撤销 100 万行
- **占用连接**：一个长事务占着一个连接不放，连接池资源浪费

怎么避免：
- 拆！大批量操作拆成小批次，每次处理 1000 条提交一次
- 事务里不要做 RPC 调用、不要做文件 IO，这些耗时操作放到事务外面
- 设置 `innodb_lock_wait_timeout` 超时时间
- 监控长事务：`SELECT * FROM information_schema.INNODB_TRX WHERE TIME_TO_SEC(TIMEDIFF(NOW(), trx_started)) > 60`

---

## 十、GORM 实战

### 33. GORM 的零值更新问题怎么解决？

这个是 GORM 的经典坑。你用 struct 更新的时候，Go 的零值（int 的 0、string 的 ""、bool 的 false）会被 GORM 忽略，不会生成到 SQL 里。

比如你想把 age 更新为 0：

```go
db.Model(&user).Updates(User{Age: 0})
// 实际 SQL: UPDATE users SET updated_at = '...'  ← age 没了！
```

解决方案就是**用 map 更新**：

```go
db.Model(&user).Updates(map[string]interface{}{"age": 0})
// 实际 SQL: UPDATE users SET age = 0, updated_at = '...'  ✓
```

map 里的值不存在"零值忽略"的问题，传什么就更新什么。

### 34. GORM 的软删除是怎么回事？

如果模型里有 `gorm.DeletedAt` 字段（`gorm.Model` 自带），调用 `db.Delete()` 的时候不会真的执行 `DELETE`，而是：

```sql
UPDATE users SET deleted_at = '2024-01-01 12:00:00' WHERE id = 1
```

后续所有查询会自动加上 `WHERE deleted_at IS NULL`，所以"看起来"数据被删了，实际还在表里。

想真的删除用 `db.Unscoped().Delete()`。想查出已删除的数据用 `db.Unscoped().Find()`。

企业里软删除很常见，好处是数据可追溯可恢复，坏处是表会越来越大，要定期清理或归档。

### 35. GORM 里事务怎么用？有什么要注意的？

推荐用闭包模式：

```go
db.Transaction(func(tx *gorm.DB) error {
    // 事务内的所有操作都用 tx，不要用 db
    if err := tx.Create(&order).Error; err != nil {
        return err  // 返回 error → 自动 Rollback
    }
    if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
        return err
    }
    return nil  // 返回 nil → 自动 Commit
})
```

注意事项：
- 事务内所有操作用 `tx`，不能用原始的 `db`，否则操作不在同一个事务里
- 如果 DAO 层要在事务中使用，需要 `WithTx(tx)` 方法创建绑定事务的新 DAO
- 不要在事务内做耗时的非 DB 操作（调接口、发消息），缩短事务持有时间

### 36. Preload 和 Joins 有什么区别？怎么选？

**Preload** 实际上是执行两条 SQL：

```sql
SELECT * FROM orders WHERE user_id = 1;
SELECT * FROM users WHERE id IN (1);
```

先查主表，收集关联 ID，再用 IN 查关联表。

**Joins** 是执行一条 JOIN SQL：

```sql
SELECT orders.*, users.* FROM orders
LEFT JOIN users ON orders.user_id = users.id
WHERE users.name = 'bob';
```

怎么选：
- 一对多关联、不需要用关联字段做 WHERE → **Preload**
- 需要用关联字段做 WHERE/ORDER BY → **Joins**
- 关联数据量很大时 → **Preload**（IN 比 JOIN 效率高）
- 一对一关联且要减少 SQL 往返 → **Joins**

### 37. GORM Hooks 有什么用？什么时候会不触发？

Hooks 是 GORM 在 Create/Update/Delete/Query 前后自动调用的方法。定义在 Model 上。

典型用法：
- `BeforeCreate`：数据清洗（email 转小写）、自动生成 UUID
- `AfterCreate`：发消息通知、写审计日志
- `BeforeUpdate`：参数标准化

**不触发的情况**（这个面试会追问）：
- `db.Updates(map[string]interface{}{...})` → 不触发
- `db.Exec("UPDATE ...")` → 不触发
- `db.UpdateColumn()` / `db.UpdateColumns()` → 不触发
- 批量操作 `db.Where(...).Updates(...)` → 不触发

只有 `db.Create(&model)` 和 `db.Save(&model)` 这种传 struct 指针的方式才触发。

---

## 十一、设计与选型

### 38. 为什么金额要用整数存分，不用 DECIMAL？

首先浮点数（FLOAT/DOUBLE）肯定不能用，`0.1 + 0.2 != 0.3` 精度问题，涉及钱的场景不可接受。

DECIMAL 可以精确表示小数，是个选择。但用整数存分更简单：
- 整数运算比 DECIMAL 快
- 不用担心小数位数对齐问题
- 代码里全程整数操作，展示的时候除以 100 就行
- BIGINT 8 字节，DECIMAL(10,2) 也差不多

大部分互联网公司用的都是 BIGINT 存分。

### 39. 主键用自增还是 UUID 还是雪花？

**自增 BIGINT**：写入性能最好（B+Tree 追加到末尾，零页分裂），简单。缺点是分库后重复、可被猜测数据量。

**UUID**：全局唯一、不依赖任何服务。缺点是 36 字符太长、完全无序导致大量页分裂、索引性能差。**不推荐做主键**。

**雪花 ID**：趋势递增（偶尔不严格连续），全局唯一，性能好。缺点是依赖时钟。

**最佳实践**：双 ID 模式。主键用自增 BIGINT（内部使用，保证 B+Tree 性能），业务号用雪花 ID（对外暴露，全局唯一）。

### 40. MySQL 和 PostgreSQL 怎么选？

**MySQL**：
- 互联网行业 OLTP 首选，生态成熟
- 主从复制简单、读写分离方案成熟
- 分库分表中间件丰富（ShardingSphere、Vitess）
- 上手简单，DBA 好招

**PostgreSQL**：
- 功能更丰富：JSONB（可索引）、数组类型、全文搜索、CTE、窗口函数更强
- MVCC 不同：用堆表多版本，不依赖 Undo Log，但需要 VACUUM 清理
- 扩展性极强，支持自定义类型、函数、索引
- 适合复杂业务和分析场景

简单说：高并发 OLTP、互联网项目 → MySQL；复杂查询、数据分析、需要高级特性 → PostgreSQL。
