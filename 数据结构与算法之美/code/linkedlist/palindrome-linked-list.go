package linkedlist

func isPalindrome(head *ListNode) bool {
	var (
		slow *ListNode = head
		fast *ListNode = head
		prev *ListNode = nil
		temp *ListNode = nil
	)
	if head == nil || head.next == nil {
		return true
	}
	for fast != nil && fast.next != nil {
		fast = fast.next.next
		temp = slow.next
		slow.next = prev
		prev = slow
		slow = temp
	}
	if fast != nil {
		slow = slow.next
	}

	var (
		l1 *ListNode = prev
		l2 *ListNode = slow
	)
	for l1 != nil && l2 != nil && l1.value == l2.value {
		l1 = l1.next
		l2 = l2.next
	}
	return l1 == nil && l2 == nil
}
