package stack

import "fmt"

// ArrayStack 基于数组实现的栈
type ArrayStack struct {
	// 数据
	data []interface{}
	// 栈顶指针
	top int
}

func NewArrayStack() *ArrayStack {
	return &ArrayStack{
		data: make([]interface{}, 0, 32),
		top:  -1,
	}
}

func (s *ArrayStack) Push(v interface{}) {
	if s.top < 0 {
		s.top = 0
	} else {
		s.top += 1
	}

	if s.top > len(s.data)-1 {
		s.data = append(s.data, v)
	} else {
		s.data[s.top] = v
	}
}

func (s *ArrayStack) Pop() interface{} {
	if s.IsEmpty() {
		return nil
	}
	v := s.data[s.top]
	s.top -= 1
	return v
}

func (s *ArrayStack) IsEmpty() bool {
	if s.top < 0 {
		return true
	}
	return false
}

func (s *ArrayStack) Top() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return s.data[s.top]
}

func (s *ArrayStack) Flush() {
	s.top = -1
}

func (s *ArrayStack) Print() {
	if s.IsEmpty() {
		fmt.Println("empty stack")
	} else {
		for i := s.top; i >= 0; i-- {
			fmt.Println(s.data[i])
		}
	}
}
