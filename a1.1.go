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
输入网址 http://127.0.0.1:6060/debug/pprof/
打开pprof主页，从上到下依次是5类profile信息：

block：goroutine的阻塞信息，本例就截取自一个goroutine阻塞的demo，但block为0，没掌握block的用法
goroutine：所有goroutine的信息，下面的full goroutine stack dump是输出所有goroutine的调用栈，是goroutine的debug=2，后面会详细介绍。
heap：堆内存的信息
mutex：锁的信息
threadcreate：线程信息

我们主要关注goroutine和heap，这两个都会打印调用栈信息，goroutine里面还会包含goroutine的数量信息，heap则是内存分配信息
*/