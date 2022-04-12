package tree

type BST struct {
	*BinaryTree
	compareFunc func(v, nodeV interface{}) int
}

func NewBST(rootV interface{}, compareFunc func(v, nodeV interface{}) int) *BST {
	if compareFunc == nil {
		return nil
	}
	return &BST{
		BinaryTree:  NewBinaryTree(rootV),
		compareFunc: compareFunc,
	}
}

func (bst *BST) Find(v interface{}) *Node {
	p := bst.root
	for p != nil {
		compareResult := bst.compareFunc(v, p.data)
		if compareResult == 0 {
			return p
		} else if compareResult > 0 {
			p = p.right
		} else {
			p = p.left
		}
	}
	return nil
}

func (bst *BST) Insert(v interface{}) bool {
	p := bst.root
	for p != nil {
		compareResult := bst.compareFunc(v, p.data)
		if compareResult == 0 {
			return false
		} else if compareResult > 0 {
			if p.right == nil {
				p.right = NewNode(v)
				break
			}
			p = p.right
		} else {
			if p.left == nil {
				p.left = NewNode(v)
				break
			}
			p = p.left
		}
	}
	return true
}

func (bst *BST) Delete(v interface{}) bool {
	var pp *Node = nil
	p := bst.root
	deleteLeft := false
	for p != nil {
		compareResult := bst.compareFunc(v, p.data)
		if compareResult > 0 {
			pp = p
			p = p.right
			deleteLeft = false
		} else if compareResult < 0 {
			pp = p
			p = p.left
			deleteLeft = true
		} else {
			break
		}
	}

	if p == nil {
		return false
	} else if p.left == nil && p.right == nil {
		if pp != nil {
			if deleteLeft {
				pp.left = nil
			} else {
				pp.right = nil
			}
		} else {
			bst.root = nil
		}
	} else if p.right != nil {
		pq := p
		q := p.right
		fromRight := true
		for q.left != nil {
			pq = q
			q = q.left
			fromRight = false
		}
		if fromRight {
			pq.right = nil
		} else {
			pq.left = nil
		}
		q.left = p.left
		q.right = p.right
		if pp == nil {
			bst.root = q
		} else {
			if deleteLeft {
				pq.left = q
			} else {
				pq.right = q
			}
		}
	} else {
		if pp != nil {
			if deleteLeft {
				pp.left = p.left
			} else {
				pp.right = p.left
			}
		} else {
			bst.root = p.left
		}
	}
	return true
}

func (bst *BST) Min() *Node {
	p := bst.root
	for p.left != nil {
		p = p.left
	}
	return p
}

func (bst *BST) Max() *Node {
	p := bst.root
	for p.right != nil {
		p = p.right
	}
	return p
}
