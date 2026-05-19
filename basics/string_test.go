package basics

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

/*
Q1: Go 字符串的底层结构是什么？
Q2: 字符串遍历有什么要注意的？
Q3: 如何修改字符串？
Q4: 字符串拼接有哪些方式？性能区别？
Q5: string 和 []byte 互转有拷贝吗？怎么零拷贝？
Q6: rune 是什么？和 byte 有什么区别？
Q7: 字符串比较的底层机制？

---
Q1: Go 字符串的底层结构是什么？
【理解】
string 底层是一个只读的字节数组引用：
  type stringHeader struct {
      Data unsafe.Pointer  // 指向底层字节数组
      Len  int             // 字节长度
  }

关键特性：
  1. 不可变（immutable）：创建后不能修改任何字节，修改需要转 []byte 或 []rune
  2. 值传递时只拷贝 header（16 字节），不拷贝底层数据
  3. 可以包含任意字节（包括 \0），不要求是合法 UTF-8
  4. len(s) 返回的是字节数，不是字符数

为什么设计成不可变？
  - 并发安全：多个 goroutine 共享字符串不需要加锁
  - 可以安全地做子串切片（共享底层数组）
  - 可以作为 map 的 key（可比较、可哈希）
【回答】
string 底层是 {指针, 长度} 的结构体，指向一个只读字节数组。
不可变设计的好处：并发安全无需加锁、子串切片可以共享底层数组、可以作为 map key。
值传递只拷贝 16 字节的 header，不拷贝底层数据。len(s) 返回字节数不是字符数。

---
Q2: 字符串遍历有什么要注意的？
【理解】
两种遍历方式，语义完全不同：

for i 遍历（按字节）：
  for i := 0; i < len(s); i++ { s[i] }  // 类型是 byte(uint8)
  中文会被截断成多个字节，看到乱码

for range 遍历（按字符）：
  for i, r := range s { }  // r 的类型是 rune(int32)
  自动做 UTF-8 解码，正确处理多字节字符
  i 是该字符的起始字节索引（不是字符索引！）

统计字符数：
  len("Go语言") = 8（字节：2+3+3）
  utf8.RuneCountInString("Go语言") = 4（字符数）
  []rune("Go语言") 的长度 = 4
【回答】
for i 遍历的是字节（byte），中文会被截断成多个字节。
for range 遍历的是字符（rune），自动 UTF-8 解码，正确处理多字节字符。
len(s) 返回字节数，要统计字符数用 utf8.RuneCountInString 或 len([]rune(s))。
注意 range 的索引 i 是字节位置不是字符位置。

---
Q3: 如何修改字符串？
【理解】
字符串不可变，修改必须转成可变类型再转回来：

方案 A：转 []byte（适用于纯 ASCII）
  s := "hello"
  b := []byte(s)   // 拷贝一份
  b[0] = 'H'
  s = string(b)    // 再拷贝回来

方案 B：转 []rune（适用于含多字节字符）
  s := "中文"
  r := []rune(s)   // 按字符拆开
  r[0] = '日'
  s = string(r)

注意：两次转换都有内存拷贝（string -> []byte/[]rune -> string）。
如果只是替换子串，用 strings.Replace 或 strings.Builder 更高效。
【回答】
字符串不可变，修改要转 []byte（纯 ASCII）或 []rune（含中文）再转回来。
两次转换都有内存拷贝。如果只是替换子串，用 strings.Replace 更高效。

---
Q4: 字符串拼接有哪些方式？性能区别？
【理解】
五种拼接方式，性能差异巨大：

1. + 拼接：每次创建新字符串，拷贝旧内容，O(n²)
   s += "a"  // 循环中极慢

2. fmt.Sprintf：反射 + 格式解析，最慢，适合格式化不适合纯拼接

3. strings.Builder：内部 []byte，append 追加，2 倍扩容，O(n)
   最后 String() 零拷贝转换

4. bytes.Buffer：和 Builder 类似，但 String() 有一次拷贝
   Builder 是 Go 1.10 引入的 Buffer 的替代品

5. strings.Join：预计算总长度，一次分配，最优（已知所有片段时）

性能排序（循环拼接场景）：
  strings.Join ≈ Builder（预分配）> Builder > bytes.Buffer >> + >> fmt.Sprintf

选择原则：
  已知所有片段 -> strings.Join
  循环拼接 -> strings.Builder（可 Grow 预分配）
  2~3 个短串 -> + 就行，可读性优先
  需要格式化 -> fmt.Sprintf
【回答】
+ 拼接每次创建新字符串 O(n²)，循环中极慢。
strings.Builder 内部用 []byte append，2 倍扩容 O(n)，String() 零拷贝。
strings.Join 预计算总长度一次分配，已知所有片段时最优。
选择：循环拼接用 Builder，已知片段用 Join，2~3 个短串直接 + 就行。

---
Q5: string 和 []byte 互转有拷贝吗？怎么零拷贝？
【理解】
标准转换有拷贝：
  s := "hello"
  b := []byte(s)   // 分配新内存，拷贝 5 字节
  s2 := string(b)  // 分配新内存，拷贝 5 字节

为什么要拷贝？
  string 是不可变的，如果 []byte 和 string 共享底层数组，
  修改 []byte 就会破坏 string 的不可变性。

零拷贝（unsafe，慎用）：
  // string -> []byte（危险：不能修改返回的 []byte！）
  func stringToBytes(s string) []byte {
      return unsafe.Slice(unsafe.StringData(s), len(s))
  }
  // []byte -> string（危险：之后不能修改原 []byte！）
  func bytesToString(b []byte) string {
      return unsafe.String(&b[0], len(b))
  }

编译器优化（自动零拷贝的场景）：
  - map 查找：m[string(b)] 不会真的分配新 string
  - 字符串比较：string(b) == "hello" 不会分配
  - 字符串拼接：string(b) + string(b2) 中间不分配
【回答】
标准转换有拷贝（保证 string 不可变性）。零拷贝可以用 unsafe 实现但很危险——不能修改共享的底层数组。
编译器会在 map 查找、字符串比较等场景自动优化为零拷贝，不需要手动 unsafe。
大多数场景标准转换就够了，只有性能极度敏感的热路径才考虑 unsafe。

---
Q6: rune 是什么？和 byte 有什么区别？
【理解】
byte = uint8，表示一个字节（0~255）
rune = int32，表示一个 Unicode 码点（可以表示任何字符）

UTF-8 编码规则：
  ASCII（0~127）：1 个字节
  中文/日文等：3 个字节
  emoji 等：4 个字节

  "A" -> 1 byte, 1 rune
  "中" -> 3 bytes, 1 rune
  "😀" -> 4 bytes, 1 rune

什么时候用 rune？
  - 需要按字符处理（统计字符数、截取前 N 个字符、反转字符串）
  - 处理非 ASCII 文本

什么时候用 byte？
  - 处理纯 ASCII 文本
  - 网络/文件 IO（传输的是字节流）
  - 性能敏感场景（rune 转换有开销）
【回答】
byte 是 uint8 表示一个字节，rune 是 int32 表示一个 Unicode 字符。
一个中文字符 = 3 个 byte = 1 个 rune。
按字符处理（截取、反转、统计）用 rune；按字节处理（IO、网络）用 byte。
for range string 自动按 rune 遍历，for i 按 byte 遍历。

---
Q7: 字符串比较的底层机制？
【理解】
Go 字符串比较是逐字节比较（memcmp），不是按字符比较。

== 比较：
  先比长度，长度不同直接 false（O(1)）
  长度相同再逐字节比较（O(n)）
  如果底层指针相同（同一个字符串），直接 true（O(1)）

< > 比较（字典序）：
  逐字节比较，第一个不同的字节决定大小
  "abc" < "abd"（c < d）
  "ab" < "abc"（前缀相同，短的小）

性能注意：
  长字符串频繁比较开销大（O(n)）
  如果只需要判断相等，先比长度可以快速排除大部分情况
  需要大小写不敏感比较用 strings.EqualFold
【回答】
字符串比较是逐字节比较（memcmp）。== 先比长度再比内容，底层指针相同直接返回 true。
< > 是字典序逐字节比较。长字符串频繁比较开销是 O(n)。
大小写不敏感比较用 strings.EqualFold，不要 ToLower 后再比（多一次分配）。

*/

// TestStringIter 字符串遍历
func TestStringIter(t *testing.T) {
	s := "Go语言"

	// for i 遍历 (byte)
	fmt.Print("byte 遍历: ")
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i])
	}
	fmt.Println()

	// for range 遍历 (rune)
	fmt.Print("rune 遍历: ")
	for _, r := range s {
		fmt.Printf("%c ", r)
	}
	fmt.Println()

	fmt.Printf("字节长度: %d, 字符数: %d\n", len(s), utf8.RuneCountInString(s))
}

// TestStringImmutable 修改字符串
func TestStringImmutable(t *testing.T) {
	// []byte 修改 ASCII
	b := []byte("hello")
	b[0] = 'H'
	fmt.Println(string(b))

	// []rune 修改中文
	r := []rune("中文")
	r[0] = '日'
	fmt.Println(string(r))
}

// TestStringConcat 拼接方式对比
func TestStringConcat(t *testing.T) {
	// strings.Builder（推荐）
	var b strings.Builder
	for i := 0; i < 100; i++ {
		b.WriteString("a")
	}
	fmt.Println("Builder len:", b.Len())

	// strings.Join（已知所有片段时最优）
	parts := []string{"Go", "is", "fast"}
	fmt.Println("Join:", strings.Join(parts, " "))
}

// BenchmarkConcat 拼接性能对比
func BenchmarkConcat(b *testing.B) {
	b.Run("Plus", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := ""
			for j := 0; j < 100; j++ {
				s += "a"
			}
			_ = s
		}
	})

	b.Run("Builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			for j := 0; j < 100; j++ {
				sb.WriteString("a")
			}
			_ = sb.String()
		}
	})
}
