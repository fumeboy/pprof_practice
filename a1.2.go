package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	// 开启pprof，监听请求
	ip := "0.0.0.0:6060"
	if err := http.ListenAndServe(ip, nil); err != nil {
		fmt.Printf("start pprof failed on %s\n", ip)
	}
}

/*
在没有浏览器的环境下，使用命令go tool pprof url可以获取指定的profile文件
此命令会发起http请求，然后下载数据到本地，之后进入交互式模式，就像gdb一样，可以使用命令查看运行信息

以下是请求的方式：
	下载cpu profile，默认从当前开始收集30s的cpu使用情况，需要等待30s，可以指定时长
		go tool pprof http://localhost:6060/debug/pprof/profile?seconds=120
	如果是 下载heap profile
		go tool pprof http://localhost:6060/debug/pprof/heap
	如果是 下载goroutine profile
		go tool pprof http://localhost:6060/debug/pprof/goroutine
	以此类推


*/