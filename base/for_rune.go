package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "abc汉字"
	for i := 0; i < len(s); i++ { // byte
		fmt.Printf("%c,", s[i])
	}
	// a,b,c,æ,±,,å,,,

	fmt.Println()
	for _, r := range s { // rune
		fmt.Printf("%c,", r)
	}
	//a,b,c,汉,字,

	//计算含汉字的字符长度
	fmt.Println(utf8.RuneCountInString(s))
}
