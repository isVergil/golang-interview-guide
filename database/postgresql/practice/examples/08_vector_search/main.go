package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 08_vector_search：pgvector 向量检索（AI 应用核心）
//
// pgvector 是 PostgreSQL 的向量相似度搜索扩展
// 核心场景：RAG（检索增强生成）、语义搜索、推荐系统、图片搜索
//
// 工作流程（RAG）：
//   1. 文档 → embedding 模型 → 向量 → 存入 pgvector
//   2. 用户提问 → embedding 模型 → 查询向量
//   3. pgvector 向量相似度搜索 → 返回 Top K 文档
//   4. 文档 + 问题 → LLM → 生成回答
//
// 距离函数：
//   <->  L2 距离（欧几里得）    值越小越相似
//   <=>  余弦距离               值越小越相似（推荐，不受向量长度影响）
//   <#>  内积距离（负内积）      值越小越相似
//
// 注意：需要先安装 pgvector 扩展
//   Docker: docker exec -it pg bash -c "apt update && apt install -y postgresql-17-pgvector"
//   然后重启容器
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	if !initExtension(ctx, pool) {
		return
	}
	initData(ctx, pool)
	distanceSearch(ctx, pool)
	indexDemo(ctx, pool)
	ragDemo(ctx, pool)

	log.Println("[08_vector_search] 演示完成")
}

func initExtension(ctx context.Context, pool *pgxpool.Pool) bool {
	_, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector`)
	if err != nil {
		log.Printf("pgvector 扩展未安装: %v", err)
		fmt.Println("请先安装 pgvector 扩展:")
		fmt.Println("  docker exec -it pg bash -c \"apt update && apt install -y postgresql-17-pgvector\"")
		fmt.Println("  docker restart pg")
		return false
	}
	fmt.Println("pgvector 扩展已启用")
	return true
}

func initData(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 初始化文档数据 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS documents`)
	pool.Exec(ctx, `
		CREATE TABLE documents (
			id        bigserial PRIMARY KEY,
			title     text NOT NULL,
			content   text NOT NULL,
			category  text NOT NULL,
			-- vector(384) 表示 384 维向量（常见 embedding 模型输出维度）
			-- 实际使用 text-embedding-3-small 是 1536 维，这里用 384 做演示
			embedding vector(384)
		)
	`)

	// 模拟文档 + embedding（实际场景由 OpenAI/本地模型生成）
	docs := []struct {
		title, content, category string
	}{
		{"Redis 缓存设计", "Redis 常用作缓存层，支持多种数据结构，通过 TTL 控制过期", "database"},
		{"Redis 分布式锁", "SETNX + 过期时间实现分布式锁，Redisson 提供看门狗机制", "database"},
		{"PostgreSQL JSONB", "PostgreSQL 原生支持 JSONB 二进制格式，配合 GIN 索引查询", "database"},
		{"Go 并发模型", "Go 使用 goroutine 和 channel 实现 CSP 并发模型", "golang"},
		{"Go 内存管理", "Go 使用三色标记清除 GC，配合写屏障实现并发垃圾回收", "golang"},
		{"Docker 容器原理", "Docker 基于 Linux namespace 和 cgroup 实现进程隔离", "devops"},
		{"Kubernetes 调度", "K8s scheduler 通过 predicate 和 priority 两阶段选择节点", "devops"},
		{"向量数据库原理", "向量数据库通过 ANN 算法实现高维向量的近似最近邻搜索", "ai"},
		{"RAG 检索增强生成", "RAG 结合向量检索和 LLM，减少幻觉提高回答准确性", "ai"},
		{"Embedding 模型", "Embedding 模型将文本映射到高维向量空间，语义相近的文本向量距离近", "ai"},
	}

	for _, d := range docs {
		// 生成模拟 embedding（实际使用 OpenAI API 或本地模型）
		vec := mockEmbedding(d.title + " " + d.content)
		pool.Exec(ctx,
			`INSERT INTO documents (title, content, category, embedding) VALUES ($1, $2, $3, $4)`,
			d.title, d.content, d.category, vec,
		)
	}
	fmt.Printf("插入 %d 篇文档（含模拟 embedding）\n", len(docs))
}

// ============================================================
// 1. 向量距离搜索
// ============================================================
func distanceSearch(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 向量相似度搜索 ==========")

	// 模拟用户查询的 embedding
	queryVec := mockEmbedding("Redis 缓存和分布式锁怎么用")

	// 余弦相似度搜索（推荐）
	fmt.Println("--- 余弦相似度 Top 5 ---")
	rows, _ := pool.Query(ctx, `
		SELECT id, title, category,
		       1 - (embedding <=> $1::vector) AS similarity
		FROM documents
		ORDER BY embedding <=> $1::vector
		LIMIT 5`, queryVec)
	defer rows.Close()

	for rows.Next() {
		var id int64
		var title, category string
		var similarity float64
		rows.Scan(&id, &title, &category, &similarity)
		fmt.Printf("  id=%d [%s] %s (相似度: %.4f)\n", id, category, title, similarity)
	}

	// L2 距离搜索
	fmt.Println("\n--- L2 距离 Top 3 ---")
	rows2, _ := pool.Query(ctx, `
		SELECT title, embedding <-> $1::vector AS l2_distance
		FROM documents
		ORDER BY embedding <-> $1::vector
		LIMIT 3`, queryVec)
	defer rows2.Close()

	for rows2.Next() {
		var title string
		var dist float64
		rows2.Scan(&title, &dist)
		fmt.Printf("  %s (L2距离: %.4f)\n", title, dist)
	}

	// 按分类过滤 + 向量搜索
	fmt.Println("\n--- 分类过滤 + 向量搜索 ---")
	rows3, _ := pool.Query(ctx, `
		SELECT title, 1 - (embedding <=> $1::vector) AS similarity
		FROM documents
		WHERE category = 'database'
		ORDER BY embedding <=> $1::vector
		LIMIT 3`, queryVec)
	defer rows3.Close()

	for rows3.Next() {
		var title string
		var sim float64
		rows3.Scan(&title, &sim)
		fmt.Printf("  [database] %s (相似度: %.4f)\n", title, sim)
	}
}

// ============================================================
// 2. 向量索引
// ============================================================
func indexDemo(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 向量索引 ==========")
	fmt.Println(`
两种索引类型：
  IVFFlat：倒排文件索引
    - 原理：先聚类，搜索时只查最近的几个聚类中心
    - 适合：中等数据量（10 万~百万级）
    - 参数：lists（聚类数），probes（搜索时查几个聚类）

  HNSW：分层可导航小世界图
    - 原理：构建多层跳表式图结构
    - 适合：大数据量，召回率更高
    - 参数：m（每层连接数），ef_construction（构建精度）
    - 缺点：占内存更多，构建更慢`)

	// IVFFlat 索引（需要先有数据才能创建）
	pool.Exec(ctx, `
		CREATE INDEX idx_docs_ivfflat ON documents 
		USING ivfflat (embedding vector_cosine_ops) 
		WITH (lists = 4)`)
	fmt.Println("\n创建 IVFFlat 索引 (lists=4)")

	// 调整搜索精度
	pool.Exec(ctx, `SET ivfflat.probes = 2`)
	fmt.Println("设置 ivfflat.probes = 2（搜索时查 2 个聚类）")

	// 删除 IVFFlat，创建 HNSW
	pool.Exec(ctx, `DROP INDEX idx_docs_ivfflat`)
	pool.Exec(ctx, `
		CREATE INDEX idx_docs_hnsw ON documents 
		USING hnsw (embedding vector_cosine_ops) 
		WITH (m = 16, ef_construction = 64)`)
	fmt.Println("创建 HNSW 索引 (m=16, ef_construction=64)")

	// 调整 HNSW 搜索精度
	pool.Exec(ctx, `SET hnsw.ef_search = 40`)
	fmt.Println("设置 hnsw.ef_search = 40（搜索时探索 40 个候选）")
}

// ============================================================
// 3. RAG 完整流程演示
// ============================================================
func ragDemo(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== RAG 检索增强生成 ==========")
	fmt.Println(`
RAG 工作流：
  ┌──────────┐     ┌────────────┐     ┌────────────┐
  │ 用户提问  │ ──→ │ Embedding  │ ──→ │ pgvector   │
  │          │     │ 模型       │     │ 相似搜索    │
  └──────────┘     └────────────┘     └─────┬──────┘
                                            │ Top K 文档
                                            ↓
  ┌──────────┐     ┌────────────┐     ┌────────────┐
  │ 最终回答  │ ←── │ LLM 生成   │ ←── │ 构造 Prompt│
  │          │     │ (GPT/...)  │     │ 文档+问题   │
  └──────────┘     └────────────┘     └────────────┘`)

	question := "PostgreSQL 怎么存储 JSON 数据？"
	fmt.Printf("\n用户问题: %s\n", question)

	// Step 1: 问题 → embedding
	queryVec := mockEmbedding(question)
	fmt.Println("Step 1: 问题转 embedding (模拟)")

	// Step 2: 向量检索 Top 3
	fmt.Println("Step 2: pgvector 检索相关文档")
	rows, _ := pool.Query(ctx, `
		SELECT title, content, 1 - (embedding <=> $1::vector) AS similarity
		FROM documents
		ORDER BY embedding <=> $1::vector
		LIMIT 3`, queryVec)
	defer rows.Close()

	var contexts []string
	for rows.Next() {
		var title, content string
		var sim float64
		rows.Scan(&title, &content, &sim)
		fmt.Printf("  检索到: %s (相似度: %.4f)\n", title, sim)
		contexts = append(contexts, fmt.Sprintf("标题: %s\n内容: %s", title, content))
	}

	// Step 3: 构造 Prompt
	prompt := fmt.Sprintf(`基于以下参考资料回答用户问题。如果资料中没有相关信息，请说明。

参考资料:
%s

用户问题: %s

请回答:`, strings.Join(contexts, "\n---\n"), question)

	fmt.Println("\nStep 3: 构造 Prompt")
	fmt.Printf("  Prompt 长度: %d 字符\n", len(prompt))
	fmt.Println("  (实际应用中将此 Prompt 发送给 LLM 获取回答)")

	fmt.Println(`
RAG 的优势:
  1. 减少幻觉：LLM 基于检索到的真实文档回答，不是凭空编造
  2. 实时数据：新文档入库后立即可被检索，无需重新训练模型
  3. 成本低：不用微调大模型，只需维护向量数据库
  4. 可溯源：回答对应的参考文档可追踪
  5. 用 PG 就够了：不需要额外的 Pinecone/Milvus，PG + pgvector 一站式搞定`)
}

// ============================================================
// 工具函数
// ============================================================

// mockEmbedding 生成模拟 embedding 向量
// 实际使用时替换为 OpenAI text-embedding-3-small 等模型调用
func mockEmbedding(text string) string {
	dim := 384
	vec := make([]string, dim)

	// 用文本内容做种子，保证相同文本生成相同向量
	seed := int64(0)
	for _, c := range text {
		seed += int64(c)
	}
	r := rand.New(rand.NewSource(seed))

	var norm float64
	vals := make([]float64, dim)
	for i := 0; i < dim; i++ {
		vals[i] = r.Float64()*2 - 1 // [-1, 1]
		norm += vals[i] * vals[i]
	}
	// L2 归一化（余弦相似度需要）
	norm = math.Sqrt(norm)
	for i := 0; i < dim; i++ {
		vec[i] = fmt.Sprintf("%.6f", vals[i]/norm)
	}

	return "[" + strings.Join(vec, ",") + "]"
}
