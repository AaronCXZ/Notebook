package linkedlist

/*
	回文链表判断
*/

// isPalindrome1 开一个栈存放链表的前半段
func isPalindrome1(l *LinkedList) bool {
	lLen := l.length
	if lLen == 0 {
		return false
	}
	if lLen == 1 {
		return true
	}

	s := make([]string, 0, lLen/2)
	cur := l.head
	for i := uint(1); i <= lLen; i++ {
		cur := cur.next
		// 奇数个节点时忽略中间的节点
		if lLen%2 != 0 && i == (lLen/2+1) {
			continue
		}
		// 前半段加入切片
		if i <= lLen/2 {
			s = append(s, cur.GetValue().(string))
		} else {
			// 后半段逐一与前半段比较，不相等则返回
			if s[lLen-i] != cur.GetValue().(string) {
				return false
			}
		}
	}
	return true
}

// isPalindrome2 找到链表中间节点，将前半部分转置，再从中间向左右遍历对比
func isPalindrome2(l *LinkedList) bool {
	lLen := l.length
	if lLen == 0 {
		return false
	}
	if lLen == 1 {
		return true
	}
	var isPalindrome = true
	step := lLen / 2
	var pre *ListNode = nil
	cur := l.head.next
	next := l.head.next.next
	for i := uint(1); i <= step; i++ {
		tmp := cur.GetNext()
		cur.next = pre
		pre = cur
		cur = tmp
		next = cur.GetNext()
	}
	mid := cur

	var left, right *ListNode = pre, nil
	if lLen%2 != 0 {
		right = mid.GetNext()
	} else {
		right = mid
	}

	for nil != left && nil != right {
		if left.GetValue().(string) != right.GetValue().(string) {
			isPalindrome = false
			break
		}
		left = left.GetNext()
		right = right.GetNext()
	}

	cur = pre
	pre = mid
	for nil != cur {
		next = cur.GetNext()
		cur.next = pre
		pre = cur
		cur = next
	}
	l.head.next = pre
	return isPalindrome
}
