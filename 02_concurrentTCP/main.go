package main

import ("fmt"
"sync")

func main() {
	num := 100
	get_intPoint(&num)
	println(num)

	sliceInt := []int{1, 2, 3, 4, 5}
	slicebytwo(&sliceInt)
	fmt.Println(sliceInt)

	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(2)
	go routine01(ch,&wg)
	go routine02(ch,&wg)
	wg.Wait()


	var wg2 sync.WaitGroup
	ch2 := make(chan int,100)
	wg2.Add(2)
	go routine03(ch2,&wg2)
	go routine04(ch2,&wg2)
	wg2.Wait()


	increment()

	increment02()


	var wg3 sync.WaitGroup
	wg3.Add(2)
	go goroutine05(&wg3)
	go goroutine06(&wg3)
	wg3.Wait()

	durations()

}
