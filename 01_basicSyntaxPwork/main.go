package main

import (
	"01_basicSyntaxPwork/solution"
	"fmt"
)

func main() {
	fmt.Println("solution1:")
	solution.Get_single_number()

	fmt.Println("solution2:")
	solution.TestPalindrome()

	fmt.Println("solution3:")
	res := solution.IsNotationValid("({)}")
	fmt.Println(res)

	fmt.Println("solution4:")
	res2 := solution.LongestCommonPrefix([]string{"flower", "flow", "flight"})
	fmt.Println(res2)

	fmt.Println("solution5:")
	res3 := solution.PlusOne([]int{1, 2, 3})
	fmt.Println(res3)

	fmt.Println("solution6:")
	res4, res5 := solution.RemoveDuplicates([]int{1, 1, 2, 2, 3, 3, 3, 4, 5})
	fmt.Println(res4, res5)

	fmt.Println("solution7:")
	res6 := solution.MergeInterval([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}})
	fmt.Println(res6)

	fmt.Println("Solution8:")
	res7 := solution.TwoSum([]int{2, 7, 11, 15}, 9)
	fmt.Println(res7)
}
