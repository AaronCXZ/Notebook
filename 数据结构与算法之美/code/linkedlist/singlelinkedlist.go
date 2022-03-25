package linkedlist

import "fmt"

/*
	单链表的基本操作
*/

type ListNode struct {
	next  *ListNode
	value interface{}
}

type LinkedList struct {
	head   *ListNode
	length uint
}

func NewListNode(v interface{}) *ListNode {
	return &ListNode{nil, v}
}

func (n *ListNode) GetNext() *ListNode {
	return n.next
}

func (n *ListNode) GetValue() interface{} {
	return n.value
}

func NewLinkedList() *LinkedList {
	return &LinkedList{NewListNode(0), 0}
}

// InsertAfter 在某个节点后面插入节点
func (l *LinkedList) InsertAfter(p *ListNode, v interface{}) bool {
	if nil == p {
		return false
	}
	newNode := NewListNode(v)
	newNode.next = p.next
	p.next = newNode
	l.length++
	return true
}

// InsertBefore 在某个节点前面插入节点
func (l *LinkedList) InsertBefore(p *ListNode, v interface{}) bool {
	if nil == p || p == l.head {
		return false
	}
	// 从头开始历遍链表，找到p节点和p节点的上一个节点
	pre := l.head
	cur := l.head.next
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.next
	}
	// 插入节点
	newNode := NewListNode(v)
	newNode.next = cur
	pre.next = newNode
	l.length++
	return true
}

// InsertToHead 在链表的头部插入节点
func (l *LinkedList) InsertToHead(v interface{}) bool {
	return l.InsertAfter(l.head, v)
}

// InsertToTail 在链表的尾部插入节点
func (l *LinkedList) InsertToTail(v interface{}) bool {
	cur := l.head
	for nil != cur.next {
		cur = cur.next
	}
	return l.InsertAfter(cur, v)
}

// FindByIndex 通过索引查找节点
func (l *LinkedList) FindByIndex(index uint) *ListNode {
	if index >= l.length {
		return nil
	}
	cur := l.head.next
	var i uint = 0
	for ; i < index; i++ {
		cur = cur.next
	}
	return cur
}

// DeleteNode 删除传入的节点
func (l *LinkedList) DeleteNode(p *ListNode) bool {
	if nil == p {
		return false
	}
	cur := l.head.next
	pre := l.head
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.next
	}
	if nil == cur {
		return false
	}
	pre.next = p.next
	p = nil
	l.length--
	return true
}

func (l *LinkedList) Print() {
	cur := l.head.next
	format := ""
	for cur != nil {
		format += fmt.Sprintf("%+v", cur.GetValue())
		cur = cur.next
		if nil != cur {
			format += "->"
		}
	}
	fmt.Println(format)
}
