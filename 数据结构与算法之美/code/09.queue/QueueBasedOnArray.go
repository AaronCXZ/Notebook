package queue

import "fmt"

type ArrayQueue struct {
	q        []interface{}
	capacity int
	head     int
	tail     int
}

func NewArrayQueue(n int) *ArrayQueue {
	return &ArrayQueue{
		q:        make([]interface{}, n),
		capacity: n,
		head:     0,
		tail:     0,
	}
}

func (q *ArrayQueue) EnQueue(v interface{}) bool {
	if q.tail == q.capacity {
		return false
	}
	q.q[q.tail] = v
	q.tail++
	return true
}

func (q *ArrayQueue) DeQueue() interface{} {
	if q.head == q.tail {
		return nil
	}
	v := q.q[q.head]
	q.head++
	return v
}

func (q *ArrayQueue) String() string {
	if q.head == q.tail {
		return "empty queue"
	}
	result := "head"
	for i := q.head; i <= q.tail-1; i++ {
		result += fmt.Sprintf("<-%+v", q.q[i])
	}
	result += "<-tail"
	return result
}
