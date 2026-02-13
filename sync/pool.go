package main

import (
	"fmt"
	"sync"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			return &Student{
				Name: "zhangsan",
				Age:  18,
			}
		},
	}

	st := pool.Get().(*Student)
	println(st.Name, st.Age)
	fmt.Printf("addr is %p\n", st)

	// 修改
	st.Name = "lisi"
	st.Age = 20

	// 回收
	pool.Put(st)

	st1 := pool.Get().(*Student)
	println(st1.Name, st1.Age)
	fmt.Printf("addr1 is %p\n", st1)
}
