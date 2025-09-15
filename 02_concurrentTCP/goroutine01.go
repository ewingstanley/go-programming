package main

import "sync"

func goroutine05(wg *sync.WaitGroup){
	defer wg.Done()
	for i := 1; i <= 10; i+=2 {
		if i/2 !=0 {
			println("goroutine01奇数:", i)
		}
		
	}
}

func goroutine06(wg *sync.WaitGroup){
	defer wg.Done()
	for i := 2; i <= 10; i+=2 {
		println("goroutine02偶数:", i)

		
	}
}
