package solution

func RemoveDuplicates(nums []int) (int, []int){
	if len(nums) == 0 {
		return 0,[]int{}
	}
	k :=1
	for i:=1;i<len(nums);i++{
		if nums[i] != nums[i-1]{
			nums[k] = nums[i]
			k++
		}
	}
	return k,nums[:k]
}