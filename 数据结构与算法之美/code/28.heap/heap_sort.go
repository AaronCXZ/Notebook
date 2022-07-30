package heap

func buildHeap(a []int, n int) {
	for i := n / 2; i >= 1; i-- {
		heapifyUpToDown(a, i, n)
	}
}
