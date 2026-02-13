package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // 隐式引入，会自动注册 /debug/pprof 路由
)

/*
1 运行程序后，访问 http://localhost:28080/leak 几次。
2 在浏览器打开 http://localhost:28080/debug/pprof/。

	你会看到 goroutine 那一栏后面的数字在不断增加。

3 看详细堆栈: 点击进入 goroutine 页面，或者直接访问 http://localhost:28080/debug/pprof/goroutine?debug=1。
4 参数不同: debug=1 分组信息  debug=2 协程完整堆栈
*/
func leakWorker() {
	ch := make(chan int)
	go func() {
		// 故意往没人读的 channel 写数据，会永远阻塞在这里
		ch <- 1
		fmt.Println("这行永远不会执行")
	}()
}

func handler(w http.ResponseWriter, r *http.Request) {
	leakWorker() // 每次访问这个 URL，就会泄露一个协程
	fmt.Fprint(w, "Done")
}

func main() {
	http.HandleFunc("/leak", handler)
	// 启动 HTTP 服务，监听 18080 端口
	http.ListenAndServe(":28080", nil)
}
