package tree

type BST struct {
	*BinaryTree
	compareFunc func(v, nodeV interface{}) int
}

// NewBST 新建二叉查找树
func NewBST(rootV interface{}, compareFunc func(v, nodeV interface{}) int) *BST {
	if compareFunc == nil {
		return nil
	}
	return &BST{
		BinaryTree:  NewBinaryTree(rootV),
		compareFunc: compareFunc,
	}
}

// Find 查询节点
func (bst *BST) Find(v interface{}) *Node {
	p := bst.root
	for p != nil {
		compareResult := bst.compareFunc(v, p.data)
		// 找到节点
		if compareResult == 0 {
			return p
		} else if compareResult > 0 { // 大于节点查询右子树
			p = p.right
		} else { // 小于节点查询左节点
			p = p.left
		}
	}
	return nil
}

// Insert 插入数据
func (bst *BST) Insert(v interface{}) bool {
	p := bst.root
	for p != nil {
		compareResult := bst.compareFunc(v, p.data)
		// 节点以存在
		if compareResult == 0 {
			return false
		} else if compareResult > 0 { // 插入到右子树
			if p.right == nil { // 没有右子树时，新节点作为右子节点
				p.right = NewNode(v)
				break
			}
			p = p.right
		} else { // 插入到左子树
			if p.left == nil { // 没有左子树时，新节点作为左子节点
				p.left = NewNode(v)
				break
			}
			p = p.left
		}
	}
	return true
}

// Delete 删除节点
func (bst *BST) Delete(v interface{}) bool {
	// 查询到节点的父节点，如果待删除的节点是其唯一子节点，子节点直接取代本节点
	var pp *Node = nil
	p := bst.root
	// 是否删除左子树节点
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
