package topinterview150

import (
	"testing"
)

// 12. 整数转罗马数字 (Integer to Roman)
//
// 题目描述:
// 给你一个整数，将其转为罗马数字。
// 罗马数字包含以下七种字符： I， V， X， L，C，D 和 M。
// 字符          数值
// I             1
// V             5
// X             10
// L             50
// C             100
// D             500
// M             1000

// 🚀 性能细节 1：将映射表定义为全局定长数组 [...]string
// 而不是函数内部的切片 []string。这样可以保证这四张表只在程序启动时初始化一次，
// 且直接存储在数据段中，函数调用时 0 内存分配！
var (
	thousands = [...]string{"", "M", "MM", "MMM"}
	hundreds  = [...]string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
	tens      = [...]string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
	ones      = [...]string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
)

// intToRoman 整数转罗马数字
func intToRoman(num int) string {
	// 🚀 性能细节 2：利用 Go 编译器的字符串拼接优化
	// 在 Go 底层，当存在明确数量的字符串通过 '+' 相连时（如 s1 + s2 + s3 + s4），
	// 编译器会调用底层的 concatstrings 函数，预先计算好总长度，
	// 然后只进行【1次】内存分配，直接写入字节。比 strings.Builder 还要快！
	return thousands[num/1000] +
		hundreds[(num%1000)/100] +
		tens[(num%100)/10] +
		ones[num%10]
}

func TestIntToRoman(t *testing.T) {
	// 整数转罗马数字测试
}
