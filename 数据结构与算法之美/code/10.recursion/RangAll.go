package recursion

import "fmt"

type RangeType struct {
	value []interface{}
}

func NewRangeArray(n int) *RangeType {
	return &RangeType{
		make([]interface{}, n),
	}
}

func (slice *RangeType) RangeAll(start int) {
	lenght := len(slice.value)
	if start == lenght-1 {
		fmt.Println(slice.value)
	}

	for i := start; i < lenght; i++ {
		if i == start || slice.value[i] == slice.value[start] {
			slice.value[i], slice.value[start] = slice.value[start], slice.value[i]
			slice.RangeAll(start + 1)
			slice.value[i], slice.value[start] = slice.value[start], slice.value[i]
		}
	}
}
