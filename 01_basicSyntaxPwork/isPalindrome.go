package main

import "fmt"

// isPalindrome 判断一个整数是否为回文数
// 回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数
func isPalindrome(x int) bool {
	// 负数不是回文数
	if x < 0 {
		return false
	}

	// 单个数字都是回文数
	if x < 10 {
		return true
	}

	// 方法1：转换为字符串进行比较
	str := fmt.Sprintf("%d", x)
	left, right := 0, len(str)-1

	for left < right {
		if str[left] != str[right] {
			return false
		}
		left++
		right--
	}

	return true
}

// isPalindromeReverse 另一种方法：通过数学运算反转数字
func isPalindromeReverse(x int) bool {
	// 负数不是回文数
	if x < 0 {
		return false
	}

	// 单个数字都是回文数
	if x < 10 {
		return true
	}

	original := x
	reversed := 0

	// 反转数字
	for x > 0 {
		reversed = reversed*10 + x%10
		x /= 10
	}

	return original == reversed
}

// 测试函数，可以在其他地方调用来测试回文数功能
func testPalindrome() {
	// 测试用例
	testCases := []int{121, -121, 10, 0, 12321, 12345, 1, 11, 1221}

	fmt.Println("回文数检测结果：")
	fmt.Println("数字\t字符串方法\t数学方法")
	fmt.Println("--------------------------------")

	for _, num := range testCases {
		result1 := isPalindrome(num)
		result2 := isPalindromeReverse(num)
		fmt.Printf("%d\t%t\t\t%t\n", num, result1, result2)
	}
}
