package main

import (
	"fmt"
	"time"
)

// 示例1: 基本的channel使用
func basicChannelExample() {
	fmt.Println("=== 基本Channel示例 ===")
	
	// 创建一个无缓冲的int类型channel
	ch := make(chan int)
	
	// 启动一个goroutine发送数据
	go func() {
		fmt.Println("发送数据: 42")
		ch <- 42 // 发送数据到channel
	}()
	
	// 从channel接收数据
	value := <-ch
	fmt.Printf("接收到数据: %d\n\n", value)
}

// 示例2: 有缓冲的channel
func bufferedChannelExample() {
	fmt.Println("=== 有缓冲Channel示例 ===")
	
	// 创建一个容量为3的缓冲channel
	ch := make(chan string, 3)
	
	// 发送数据（不会阻塞，因为有缓冲）
	ch <- "第一条消息"
	ch <- "第二条消息"
	ch <- "第三条消息"
	
	fmt.Printf("Channel长度: %d, 容量: %d\n", len(ch), cap(ch))
	
	// 接收数据
	for i := 0; i < 3; i++ {
		msg := <-ch
		fmt.Printf("接收: %s\n", msg)
	}
	fmt.Println()
}

// 示例3: 使用channel进行goroutine同步
func workerExample() {
	fmt.Println("=== Worker模式示例 ===")
	
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	
	// 启动3个worker goroutine
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	
	// 发送5个任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // 关闭jobs channel，表示没有更多任务
	
	// 收集结果
	for r := 1; r <= 5; r++ {
		result := <-results
		fmt.Printf("结果: %d\n", result)
	}
	fmt.Println()
}

// worker函数
func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d 开始处理任务 %d\n", id, job)
		time.Sleep(time.Second) // 模拟工作
		fmt.Printf("Worker %d 完成任务 %d\n", id, job)
		results <- job * 2 // 发送结果
	}
}

// 示例4: 使用select进行多路复用
func selectExample() {
	fmt.Println("=== Select多路复用示例 ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	// 启动两个goroutine
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自channel1的消息"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自channel2的消息"
	}()
	
	// 使用select等待任一channel有数据
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("收到: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("收到: %s\n", msg2)
		case <-time.After(3 * time.Second):
			fmt.Println("超时!")
		}
	}
	fmt.Println()
}

// 示例5: 使用channel实现生产者-消费者模式
func producerConsumerExample() {
	fmt.Println("=== 生产者-消费者模式示例 ===")
	
	ch := make(chan int, 2)
	done := make(chan bool)
	
	// 生产者
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("生产: %d\n", i)
			ch <- i
			time.Sleep(500 * time.Millisecond)
		}
		close(ch) // 生产完毕，关闭channel
	}()
	
	// 消费者
	go func() {
		for num := range ch { // range会自动处理channel关闭
			fmt.Printf("消费: %d\n", num)
			time.Sleep(1 * time.Second)
		}
		done <- true
	}()
	
	<-done // 等待消费者完成
	fmt.Println()
}

// 示例6: 单向channel示例
func unidirectionalChannelExample() {
	fmt.Println("=== 单向Channel示例 ===")
	
	ch := make(chan int)
	
	// 只能发送的channel
	go sender(ch)
	
	// 只能接收的channel
	receiver(ch)
	fmt.Println()
}

// 只能发送数据的函数
func sender(ch chan<- int) {
	for i := 1; i <= 3; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}
	close(ch)
}

// 只能接收数据的函数
func receiver(ch <-chan int) {
	for value := range ch {
		fmt.Printf("接收: %d\n", value)
	}
}

func main() {
	basicChannelExample()
	bufferedChannelExample()
	workerExample()
	selectExample()
	producerConsumerExample()
	unidirectionalChannelExample()
	
	fmt.Println("所有Channel示例执行完毕!")
}
