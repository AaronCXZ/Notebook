package queue

import "fmt"

type CircularQueue struct {
	q        []interface{}
	capacity int
	head     int
	tail     int
}

func NewCircularQueue(n int) *CircularQueue {
	if n == 0 {
		return nil
	}
	return &CircularQueue{make([]interface{}, n), n, 0, 0}
}

// IsEmpty 栈空条件：head == tail
func (cq *CircularQueue) IsEmpty() bool {
	if cq.head == cq.tail {
		return true
	}
	return false
}

// IsFull 栈满条件：(tail+1)%capacity == head
func (cq *CircularQueue) IsFull() bool {
	if cq.head == (cq.tail+1)%cq.capacity {
		return true
	}
	return false
}

func (cq *CircularQueue) EnQueue(v interface{}) bool {
	if cq.IsFull() {
		return false
	}
	cq.q[cq.tail] = v
	cq.tail = (cq.tail + 1) % cq.capacity
	return true
}

func (cq *CircularQueue) DeQueue() interface{} {
	if cq.IsEmpty() {
		return nil
	}
	v := cq.q[cq.head]
	cq.head = (cq.head + 1) % cq.capacity
	return v
}

func (cq *CircularQueue) String() string {
	if cq.IsEmpty() {
		return "empty queue"
	}
	result := "head"
	var i = cq.head
	for {
		result += fmt.Sprintf("<-%+v", cq.q[i])
		i = (i + 1) % cq.capacity
		if i == cq.tail {
			break
		}
	}
	result += "<-tail"
	return result
}
