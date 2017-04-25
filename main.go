package main

import (
	"runtime"
	"sync"

	"github.com/zhangjunfang/webchat/connect"

	//mytcp "github.com/zhangjunfang/im2/tcp"
)

//1.解析ymal文件  待定
//2.获取连接
//3.

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(2)
	go connect.MainService(wg)
	go connect.TickTime(wg)
	wg.Wait()
}
