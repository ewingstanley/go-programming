package main

import (
	"fmt"
	"sync"
)

//编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

func increment() {
	var (
		counter int
		mutex   sync.Mutex
		wg      sync.WaitGroup
	)

	// 启动10个goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(goroutineID int) { // 传递参数避免闭包陷阱
			defer wg.Done()
			fmt.Printf("Goroutine %d 开始工作\n", goroutineID)
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
			fmt.Printf("Goroutine %d 完成工作\n", goroutineID)
		}(i) // 传递 i 的值
	}

	wg.Wait()
	fmt.Printf("mutex的计数结果: %d\n", counter)
}
