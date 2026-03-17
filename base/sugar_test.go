package base

import (
	"fmt"
	"testing"
)

type ServerConfig struct {
	Port int
	Host string
}

func (s *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// TestSugar 展示 Go 核心语法糖及其细节
func TestSugar(t *testing.T) {

	// 1. [...] 数组初始化语法糖
	// 编译器自动推导长度为 3，类型为 [3]string
	hosts := [...]string{"localhost", "127.0.0.1", "0.0.0.0"}

	// 2. 切片区间语法糖 [low:high]
	// 忽略 low 默认为 0
	subHosts := hosts[:2]

	// 3. 短变量声明 := (局部推导)
	// 4. range 迭代语法糖
	for i, host := range subHosts {

		// 5. 结构体部分初始化语法糖
		// 未定义的 Port 自动填充零值 0
		conf := ServerConfig{
			Host: host,
		}
		conf.Port = 8080 + i

		// 6. 方法接收器自动转换语法糖
		// conf 是值类型，但 GetAddr 需要指针类型 (*ServerConfig)
		// Go 自动将其处理为 (&conf).GetAddr()
		addr := conf.GetAddr()

		t.Logf("Index: %d, Generated Addr: %s", i, addr)
	}

	// 7. ... 打散切片语法糖 (Unpacking)
	extraHosts := []string{"192.168.1.1", "10.0.0.1"}

	// 模拟将两个切片合并
	allHosts := append(subHosts, extraHosts...)

	if len(allHosts) != 4 {
		t.Errorf("Expected 4 hosts, got %d", len(allHosts))
	}
}

// 性能进阶测试：展示 range 的优化细节
func BenchmarkRangePerformance(b *testing.B) {
	// 准备一个大数据量切片
	type BigStruct struct {
		Data [1024]int // 8KB 的超大结构体
	}
	items := make([]BigStruct, 1000)

	b.Run("Value_Copy_Sugar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// ❌ 这种语法糖会发生 8KB 的值拷贝，大数据下影响性能
			for _, item := range items {
				_ = item.Data[0]
			}
		}
	})

	b.Run("Index_Access_Optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// ✅ 这种语法糖只获取索引，性能极高
			for i := range items {
				_ = items[i].Data[0]
			}
		}
	})
}
