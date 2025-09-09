package solution

import "sort"


func MergeInterval(interval [][]int)[][]int{
	sort.Slice(interval,func (i,j int) bool  {
		return interval[i][0] < interval[j][0]
	})

	// 创建一个空的二维切片准备存放相邻比较合并的区间
	merged := [][]int{interval[0]}
	current := interval[0]

	for i:=1;i<len(interval);i++{
		next := interval[i]
		if current[1]>=next[0]{
			current[1] = max(current[1],next[1])
			merged[len(merged)-1] = current
		}else{
			merged = append(merged,next)
			current = next
		}
	}

	return merged
}

func max(a,b int)int{
	if a>b{
		return a
	}else{
		return b
	}
}