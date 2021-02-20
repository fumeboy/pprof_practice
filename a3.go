package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

// 运行一段时间：fatal error: runtime: out of memory
func main() {
	// 开启pprof
	go func() {
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()

	tick := time.Tick(time.Second / 100)
	var buf []byte
	for range tick {
		buf = append(buf, make([]byte, 1024*1024)...)
	}
}

/*
	怎么用heap发现内存问题
		使用pprof的heap能够获取程序运行时的内存信息
		在程序平稳运行的情况下，每个一段时间使用heap获取内存的profile，然后使用base能够对比两个profile文件的差别
		就像diff命令一样显示出增加和减少的变化

	将上面代码运行起来，执行以下命令获取profile文件，Ctrl-D退出，过一段时间后再获取一次。
	go tool pprof http://localhost:6060/debug/pprof/heap

	可以获取到两个profile文件（pprof 打印出的红色字体标注了文件地址）
	使用base把001文件作为基准，然后用002和001对比，先执行top看top的对比，然后执行list main列出main函数的内存对比，结果如下：
	```
	go tool pprof -base ./file1 ./file2

	(pprof) top
	Showing nodes accounting for 970.34MB, 32.30% of 3003.99MB total
		  flat  flat%   sum%        cum   cum%
	  970.34MB 32.30% 32.30%   970.34MB 32.30%  main.main   // 这里的 970.34MB 是对比前一个文件的数据而言的
			 0     0% 32.30%   970.34MB 32.30%  runtime.main
	```

	heap能显示内存的分配情况，以及哪行代码占用了多少内存，我们能轻易的找到占用内存最多的地方
	如果这个地方的数值还在不断怎大，基本可以认定这里就是内存泄露的位置
*/
