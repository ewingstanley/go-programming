package main

//题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。

func slicebytwo(sliceInt *[]int) {
	for i := range *sliceInt {
		(*sliceInt)[i] *=2
	}
}
