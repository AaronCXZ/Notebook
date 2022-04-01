package sorts

// BubbleSort 冒泡排序，a是数组，n是数组的大小
func BubbleSort(a []int, n int) {
	if n <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		// 退出标志
		flag := false
		for j := 0; j < n-i-1; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
				// 此次冒泡有数据交换
				flag = true
			}
		}
		// 如果没有数据交换，退出排序
		if !flag {
			break
		}
	}
}

// InsertionSort 插入排序，a是数组，n是数组的大小
func InsertionSort(a []int, n int) {
	if n <= 1 {
		return
	}
	for i := 1; i < n; i++ {
		value := a[i]
		j := i - 1
		for ; j >= 0; j-- {
			if a[j] > value {
				a[j+1] = a[j]
			} else {
				break
			}
		}
		a[j+1] = value
	}
}

func SelectionSort(a []int, n int) {
	if n <= 1 {
		return
	}

	for i := 0; i < n; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if a[j] < a[minIndex] {
				minIndex = j
			}
		}
		a[i], a[minIndex] = a[minIndex], a[i]
	}
}
