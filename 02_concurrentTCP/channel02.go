package main

import "sync"

//实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。

func routine03(ch chan int, wg *sync.WaitGroup){
	defer wg.Done()
	defer close(ch)
	for i :=1;i<=100;i++{
		ch <-i
		
	}
}

func routine04(ch chan int, wg *sync.WaitGroup){
	defer wg.Done()
	for i := range ch{
		println(i)
	}
}