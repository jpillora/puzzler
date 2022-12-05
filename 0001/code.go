package p0001

func twoSum(nums []int, target int) []int {
	for i, a := range nums {
		for j := i + 1; j < len(nums); j++ {
			b := nums[j]
			if (a + b) == target {
				return []int{i, j}
			}
		}
	}
	return nil
}
