package solution

func PlusOne(nums []int)([]int){
	// 思路：[1,2,3],[9,9,9]
	// 看长度多少，从最后一位去判断：如果大于9就变为0，进位+1接着判断以此类推，如果都是9那就直接返回一个新切片
	n := len(nums)

	for i :=n-1; i >=0; i--{
		if nums[i] < 9 {
			nums[i]++
			return nums
		}
		nums[i] = 0
	}

	result := make([]int, n+1)
	result[0] = 1

	return result
}