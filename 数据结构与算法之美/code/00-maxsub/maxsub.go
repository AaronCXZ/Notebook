package main

func MaxSubArray(nums []int) int {
	res := 0
	temp := 0
	for _, num := range nums {
		temp += num
		if res < temp {
			res = temp
		}
		if temp < 0 {
			temp = 0
		}
	}
	return res
}
