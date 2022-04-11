package lru

const (
	hostbit = uint64(^uint(0)) == ^uint64(0)
	LENGTH  = 100
)

type lruNode struct {
	prev *lruNode // 前驱节点
	next *lruNode // 后继节点

	key   int
	value int

	hnext *lruNode // 拉链
}

type LRUCache struct {
	node []lruNode // hash数据，hash(key)的值尾下标，对应的数据为lruNode

	head *lruNode // lru 头节点
	tail *lruNode // lru 尾节点

	capacity int // 缓存大小
	used     int // 易使用缓存大小
}

// hash hash函数
func hash(key int) int {
	if hostbit {
		return (key ^ (key >> 32)) & (LENGTH - 1)
	}
	return (key ^ (key >> 16)) & (LENGTH - 1)
}

// Constructor 构造函数
func Constructor(capacity int) LRUCache {
	return LRUCache{
		node:     make([]lruNode, LENGTH),
		head:     nil,
		tail:     nil,
		capacity: capacity,
		used:     0,
	}
}

// addNode 添加节点
func (l *LRUCache) addNode(key int, value int) {
	// 根据数据新建节点
	newNode := &lruNode{
		key:   key,
		value: value,
	}

	// 从数据查询key是否存在
	tmp := &l.node[hash(key)]
	// 将新节点的hnext插入到tmp.hnext后面
	newNode.hnext = tmp.hnext
	tmp.hnext = newNode
	l.used++

	// 如果尾节点为空，新节点即为头尾节点
	if l.tail == nil {
		l.tail, l.head = newNode, newNode
		return
	}
	// 将新节点插入双向链表的尾节点
	l.tail.next = newNode
	newNode.prev = l.tail
	l.tail = newNode
}

func (l *LRUCache) delNode() {
	// 头结点为空则缓存没有数据
	if l.head == nil {
		return
	}
	// 从数组查询到头结点对应的指针，然后从拉链查询到节点
	prev := &l.node[hash(l.head.key)]
	tmp := prev.hnext

	// 从拉链删除头结点
	for tmp != nil && tmp.key != l.head.key {
		prev = tmp
		tmp = tmp.hnext
	}
	if tmp == nil {
		return
	}

	// 从双向链表删除头结点
	prev.hnext = tmp.hnext
	l.head = l.head.next
	l.head.prev = nil
	l.used--
}

func (l *LRUCache) searchNode(key int) *lruNode {
	// 缓存为空
	if l.tail == nil {
		return nil
	}

	// 从拉链查找数据
	tmp := l.node[hash(key)].hnext
	for tmp != nil {
		// hash冲突处理
		if tmp.key == key {
			return tmp
		}
		tmp = tmp.hnext
	}
	return nil
}

func (l *LRUCache) moveToTail(node *lruNode) {
	// 是否尾尾节点
	if l.tail == node {
		return
	}
	// 是否尾头结点，如果是头结点则删除头结点，否则将对应的节点删除
	if l.head == node {
		l.head = node.next
		l.head.prev = nil
	} else {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	// 将节点插入到尾节点后边
	node.next = nil
	l.tail.next = node
	node.prev = l.tail

	l.tail = node
}

// Get 查询
func (l *LRUCache) Get(key int) int {
	// 如果尾节点为空，表明没有数据
	if l.tail == nil {
		return -1
	}
	// 根据keu查询到节点，如果存在将节点移动到尾节点
	if tmp := l.searchNode(key); tmp != nil {
		l.moveToTail(tmp)
		return tmp.value
	}
	// 没有对应的数据
	return -1
}

// Put 插入数据
func (l *LRUCache) Put(key int, value int) {
	// 根据keu查询数据在不在缓存，在则移动到尾节点
	if tmp := l.searchNode(key); tmp != nil {
		tmp.value = value
		l.moveToTail(tmp)
		return
	}
	// 数据不在缓存，则将数据插入缓存
	l.addNode(key, value)

	// 查询缓存是否超出容量
	if l.used > l.capacity {
		l.delNode()
	}
}
