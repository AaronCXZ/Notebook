package main

import "testing"

func TestMaxSubArray(t *testing.T) {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	res := MaxSubArray(nums)
	if res != 6 {
		t.Error("Error")
	}

}
