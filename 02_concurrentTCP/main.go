package main

import "fmt"

func main() {
	num := 100
	get_intPoint(&num)
	println(num)

	sliceInt := []int{1, 2, 3, 4, 5}
	slicebytwo(&sliceInt)
	fmt.Println(sliceInt)

	println()

}
