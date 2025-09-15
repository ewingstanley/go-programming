package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func increment02(){
	var (
		counter int64
		wg      sync.WaitGroup
	)

	// 启动10个goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(goroutineID int) { // 传递参数避免闭包陷阱
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}(i) // 传递 i 的值
	}

	wg.Wait()

	finalCount := atomic.LoadInt64(&counter)
	fmt.Printf("atomic的计数结果: %d\n", finalCount)
}