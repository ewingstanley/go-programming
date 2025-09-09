package solution

import "fmt"


// 给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现两次。
// 找出那个只出现了一次的元素。
func singleNumber(nums []int) int {
	var numMap = make(map[int]int)

	// 通过map集合获取每个元素出现次数
	for _, num := range nums {
		numMap[num]++
	}

	//遍历map，找出次数（val）为1得元素，即索引。
	for index := range numMap {
		if numMap[index] == 1 {
			return index
		}
	}
	return -1
}



func Get_single_number() {
	var numList = []int{1,4,4}

	result := singleNumber(numList)
	
	fmt.Println(result)
}



