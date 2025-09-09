package solution

import "fmt"

// isPalindrome 判断一个整数是否为回文数
// 回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数
func isPalindrome(x int) bool {
	// 负数不是回文数
	if x < 0 {
		return false
	}
	if x < 10 {
		return true
	}

	// 转换为字符串进行比较
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

func TestPalindrome() {
	testCases := []int{121, -121, 10, 0, 12321, 12345, 1, 11, 1221}
	for _, num := range testCases {
		result1 := isPalindrome(num)
		fmt.Printf("%d\t%t\n", num, result1)
	}
}
