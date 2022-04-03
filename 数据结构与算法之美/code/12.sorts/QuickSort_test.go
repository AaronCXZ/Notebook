package sorts

import (
	"math/rand"
	"testing"
)

func createRandomArray(length int) []int {
	arr := make([]int, length, length)
	for i := 0; i < length; i++ {
		arr[i] = rand.Intn(100)
	}
	return arr
}

func TestQuickSort(t *testing.T) {
	arr := []int{5, 4}
	QuickSort(arr)
	t.Log(arr)

	arr = createRandomArray(10)
	QuickSort(arr)
	t.Log(arr)
}
