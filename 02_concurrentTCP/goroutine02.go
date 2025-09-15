package main

import (
	"fmt"
	"sync"
	"time"
)

func durations() {
	tasks := []func(){
		func() { time.Sleep(1 * time.Second); fmt.Println("任务1完成") },
		func() { time.Sleep(2 * time.Second); fmt.Println("任务2完成") },
		func() { time.Sleep(500 * time.Millisecond); fmt.Println("任务3完成") },
	}

	var wg sync.WaitGroup
	durations := make([]time.Duration, len(tasks))

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, job func()) {
			defer wg.Done()
			start := time.Now()
			job()
			durations[index] = time.Since(start)
		}(i, task)
	}

	wg.Wait()

	fmt.Println("\n任务执行时间统计:")
	for i, d := range durations {
		fmt.Printf("任务%d: %v\n", i+1, d)
	}
}