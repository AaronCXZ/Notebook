package array

import (
	"errors"
	"fmt"
)

/**
数组的插入。删除，按照下标随即访问操作
数组中的数据时int类型
*/

type Array struct {
	data   []int
	length uint
}

func NewArray(capacity uint) *Array {
	if capacity == 0 {
		return nil
	}
	return &Array{
		data:   make([]int, capacity, capacity),
		length: 0,
	}
}

func (a *Array) Len() uint {
	return a.length
}

// isIndexOutOfRange 判断索引是否越界
func (a *Array) isIndexOutOfRange(index uint) bool {
	if index >= uint(cap(a.data)) {
		return true
	}
	return false
}

// Find 根据索引查找数组
func (a *Array) Find(index uint) (int, error) {
	if a.isIndexOutOfRange(index) {
		return 0, errors.New("越界")
	}
	return a.data[index], nil
}

// Insert 插入数据到索引index处
func (a *Array) Insert(index uint, v int) error {
	if a.Len() == uint(cap(a.data)) {
		return errors.New("数组已满")
	}
	if index != a.length && a.isIndexOutOfRange(index) {
		return errors.New("越界")
	}
	// index之后的数据依次向后移动一位
	for i := a.length; i > index; i-- {
		a.data[i] = a.data[i-1]
	}
	a.data[index] = v
	a.length++
	return nil
}

func (a *Array) InsertToTail(v int) error {
	return a.Insert(a.Len(), v)
}

// Delete 删除索引index的值
func (a *Array) Delete(index uint) (int, error) {
	if a.isIndexOutOfRange(index) {
		return 0, errors.New("越界")
	}
	v := a.data[index]
	for i := index; i < a.Len()-1; i++ {
		a.data[i] = a.data[i+1]
	}
	a.length--
	return v, nil
}

func (a *Array) Print() {
	var format string
	for i := uint(0); i < a.Len(); i++ {
		format += fmt.Sprintf("|%+v", a.data[i])
	}
	fmt.Println(format)
}
