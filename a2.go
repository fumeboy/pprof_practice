// 之前的程序太简单了，如果去获取内存profile，几乎获取不到什么
// 现在这个程序会不断的申请内存

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

// 运行一段时间会 fatal error: runtime: out of memory
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
编译运行 go run ./a2.go
然后执行 go tool pprof http://localhost:6060/debug/pprof/heap
之后会进入 pprof 交互模式
在交互模式中输入 help , 输出如下

	```
	~ go tool pprof http://localhost:6060/debug/pprof/heap
	Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
	Saved profile in /home/fumeboy/pprof/pprof.a2.alloc_objects.alloc_space.inuse_objects.inuse_space.001.pb.gz
	File: a2
	Type: inuse_space
	Time: Feb 20, 2021 at 7:43pm (CST)
	Entering interactive mode (type "help" for commands, "o" for options)
	(pprof) help
	  Commands:
		callgrind        Outputs a graph in callgrind format
		comments         Output all profile comments
		disasm           Output assembly listings annotated with samples
		dot              Outputs a graph in DOT format
		eog              Visualize graph through eog
	```

关于命令，一般只会用到3个，也是最常用的：top、list、traces

# top
	```
	(pprof) top
	Showing nodes accounting for 814.62MB, 100% of 814.62MB total
		  flat  flat%   sum%        cum   cum%
	  814.62MB   100%   100%   814.62MB   100%  main.main
			 0     0%   100%   814.62MB   100%  runtime.main
	```
	按指标大小列出前10个函数，比如内存是按内存占用多少，CPU是按执行时间多少

	top会列出5个统计数据：
		flat: 本函数占用的内存量。
		flat%: 本函数内存占使用中内存总量的百分比。
		sum%: 前面每一行flat百分比的和，比如第2行虽然是100%，其实是 100% + 0%。
		cum: 是累计量，加入main函数调用了函数f，函数f占用的内存量，也会记进来。
		cum%: 是累计量占总量的百分比。

# list
	查看某个函数的代码，以及该函数每行代码的指标信息
	如果函数名不明确，会进行模糊匹配，比如list main会列出main.main和runtime.main

	```
	(pprof) list main.main
	Total: 814.62MB
	ROUTINE ======================== main.main in /home/ubuntu/heap/demo2.go
	  814.62MB   814.62MB (flat, cum)   100% of Total
			 .          .     20:    }()
			 .          .     21:
			 .          .     22:    tick := time.Tick(time.Second / 100)
			 .          .     23:    var buf []byte
			 .          .     24:    for range tick {
	  814.62MB   814.62MB     25:        buf = append(buf, make([]byte, 1024*1024)...)
			 .          .     26:    }
			 .          .     27:}
			 .          .     28:
	```
	可以看到在main.main中的第25行占用了814.62MB内存，左右2个数据分别是flat和cum，含义和top中解释的一样

# traces
	打印所有调用栈，以及调用栈的指标信息。
	```
	(pprof) traces
	File: a2
	Type: inuse_space
	Time: Feb 20, 2021 at 7:43pm (CST)
	-----------+-------------------------------------------------------
		 bytes:  333.19MB
	  333.19MB   main.main
				 runtime.main
	-----------+-------------------------------------------------------
		 bytes:  266.55MB
			 0   main.main
				 runtime.main
	-----------+-------------------------------------------------------
		 bytes:  213.23MB
			 0   main.main
				 runtime.main
	-----------+-------------------------------------------------------
		 bytes:  170.59MB
			 0   main.main
				 runtime.main

	```
*/