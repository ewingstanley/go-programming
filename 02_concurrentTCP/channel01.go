package main

import "sync"

//编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，
// 并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。

func routine01(ch chan int, wg *sync.WaitGroup){
	defer wg.Done()
	defer close(ch)
	for i :=1;i<=10;i++{
		ch <-i
		
	}
}

func routine02(ch chan int, wg *sync.WaitGroup){
	defer wg.Done()
	for i := range ch{
		println(i)
	}
}