package linkedlist

import "fmt"

// ListNode 单链表节点
type ListNode struct {
	next  *ListNode
	value interface{}
}

// LinkedList 单链表
type LinkedList struct {
	head *ListNode
}

func (l *LinkedList) Print() {
	cur := l.head.next
	format := ""
	for cur != nil {
		format += fmt.Sprintf("%+v", cur.value)
		cur = cur.next
		if cur != nil {
			format += "->"
		}
	}
	fmt.Println(format)
}

// Reverse 单链表反转，时间复杂度：O(n)
func (l *LinkedList) Reverse() {
	if l.head == nil || l.head.next == nil || l.head.next.next == nil {
		return
	}
	var pre *ListNode = nil
	cur := l.head.next
	for cur != nil {
		tmp := cur.next
		cur.next = pre
		pre = cur
		cur = tmp
	}
	l.head.next = pre
}

// HasCycle 判断单链表是否有环
func (l *LinkedList) HasCycle() bool {
	if l.head != nil {
		slow := l.head
		fast := l.head
		for fast != nil && fast.next != nil {
			slow = slow.next
			fast = fast.next.next
			if slow == fast {
				return true
			}
		}
	}
	return false
}

// MergeSortedList 两个有序链表合并
func MergeSortedList(l1, l2 *LinkedList) *LinkedList {
	if l1 == nil || l1.head == nil || l1.head.next == nil {
		return l2
	}
	if l2 == nil || l2.head == nil || l2.head.next == nil {
		return l1
	}
	l := &LinkedList{head: &ListNode{}}
	cur := l.head
	curl1 := l1.head.next
	curl2 := l2.head.next
	for curl1 != nil && curl2 != nil {
		if curl1.value.(int) > curl2.value.(int) {
			cur.next = curl2
			curl2 = curl2.next
		} else {
			cur.next = curl1
			curl1 = curl1.next
		}
		cur = cur.next
	}
	if curl1 != nil {
		cur.next = curl1
	} else if curl2 != nil {
		cur.next = curl2
	}
	return l
}

// DeleteBottomN 删除倒数第n和节点
func (l *LinkedList) DeleteBottomN(n int) {
	if n <= 0 || l.head == nil || l.head.next == nil {
		return
	}

	fast := l.head
	for i := 1; i <= n && fast != nil; i++ {
		fast = fast.next
	}
	if fast == nil {
		return
	}
	slow := l.head
	for fast.next != nil {
		slow = slow.next
		fast = fast.next
	}
	slow.next = slow.next.next
}

// FindMiddleNode 获取中间节点
func (l *LinkedList) FindMiddleNode() *ListNode {
	if l.head == nil || l.head.next == nil {
		return nil
	}
	if l.head.next.next == nil {
		return l.head.next
	}
	slow, fast := l.head, l.head
	for fast != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next
	}
	return slow
}
