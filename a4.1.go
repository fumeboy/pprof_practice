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
		buf = append(buf, make([]byte, 1024*4)...)
	}
}

/*
$ go tool pprof http://0.0.0.0:6060/debug/pprof/goroutine
Fetching profile over HTTP from http://localhost:6061/debug/pprof/goroutine
Saved profile in /home/ubuntu/pprof/pprof.leak_demo.goroutine.001.pb.gz  // profile文件保存位置
File: leak_demo
Type: goroutine
Time: May 16, 2019 at 2:44pm (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)

同样是 top list traces

top 方法同上，直接看 traces

	```
	(pprof) traces
	File: leak_demo
	Type: goroutine
	Time: May 16, 2019 at 2:44pm (CST)
	-----------+-------------------------------------------------------
		 20312   runtime.gopark
				 runtime.goparkunlock
				 runtime.chansend
				 runtime.chansend1 // channel发送
				 main.alloc2.func1 // alloc2中的匿名函数
				 main.alloc2
	-----------+-------------------------------------------------------
	```

traces能列出002中比001中多的那些goroutine的调用栈
这里只有1个调用栈，有20312个goroutine都执行这个调用路径
可以看到alloc2中的匿名函数alloc2.func1调用了写channel的操作，然后阻塞挂起了goroutine

用 list 命令查看 alloc2.func1

```
(pprof) list main.alloc2.func1
Total: 20312
ROUTINE ======================== main.alloc2.func1 in /home/ubuntu/heap/leak_demo.go
         0      20312 (flat, cum)   100% of Total
         .          .     48:        // 分配内存，假用一下
         .          .     49:        buf := make([]byte, 1024*1024*10)
         .          .     50:        _ = len(buf)
         .          .     51:        fmt.Println("alloc done")
         .          .     52:
         .      20312     53:        outCh <- 0  // 看这
         .          .     54:    }()
         .          .     55:}
         .          .     56:
```

配合使用 top list traces 就找到了问题位置
*/