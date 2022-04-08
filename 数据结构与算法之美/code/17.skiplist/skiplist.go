package skiplist

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	// 最高层数
	MAX_LEVEL = 16
)

type skiplistNode struct {
	// 跳表保存的值
	v interface{}
	// 用于排序的分值
	score int
	// 层高
	level int
	// 每层前进指针
	forwords []*skiplistNode
}

// newSkipListNode 新建跳表节点
func newSkipListNode(v interface{}, score, levle int) *skiplistNode {
	return &skiplistNode{
		v:        v,
		score:    score,
		level:    levle,
		forwords: make([]*skiplistNode, levle, levle),
	}
}

type SkipList struct {
	head   *skiplistNode // 跳表的头节点
	level  int           // 跳表当前对策层数
	length int           // 跳表的长度
}

func NewSkipList() *SkipList {
	head := newSkipListNode(0, math.MinInt32, MAX_LEVEL)
	return &SkipList{
		head:   head,
		level:  1,
		length: 0,
	}
}

// Length 跳表长度
func (sl *SkipList) Length() int {
	return sl.length
}

// Level 跳表层级
func (sl *SkipList) Level() int {
	return sl.level
}

// Insert 插入数据到跳表中
func (sl *SkipList) Insert(v interface{}, score int) int {
	if v == nil {
		return 1
	}

	//	查找插入的位置
	cur := sl.head
	// 记录每层的路径
	update := [MAX_LEVEL]*skiplistNode{}
	i := MAX_LEVEL - 1
	for ; i >= 0; i-- {
		for cur.forwords[i] != nil {
			// 数据已存在
			if cur.forwords[i].v == v {
				return 2
			}
			// 分值小于节点的分值
			if cur.forwords[i].score > score {
				update[i] = cur
				break
			}
			cur = cur.forwords[i]
		}
		if cur.forwords[i] == nil {
			update[i] = cur
		}
	}

	// 通过随机算法获取该节点层数
	level := 1
	for i := 1; i < MAX_LEVEL; i++ {
		if rand.Int31()%7 == 1 {
			level++
		}
	}

	// 创建一个新地跳表节点
	newNode := newSkipListNode(v, score, level)

	// 原有节点连接
	for i := 0; i <= level-1; i++ {
		next := update[i].forwords[i]
		update[i].forwords[i] = newNode
		newNode.forwords[i] = next
	}

	// 如果当前节点的层数大于之前跳表的层数，更新当前跳表的层数
	if level > sl.level {
		sl.level = level
	}

	// 更新跳表的长度
	sl.length++
	return 0
}

// Find 查找
func (sl *SkipList) Find(v interface{}, score int) *skiplistNode {
	if v == nil || sl.length == 0 {
		return nil
	}

	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forwords[i] != nil {
			if cur.forwords[i].score == score && cur.forwords[i].v == v {
				return cur.forwords[i]
			} else if cur.forwords[i].score > score {
				break
			}
			cur = cur.forwords[i]
		}
	}
	return nil
}

// Delete 删除节点
func (sl *SkipList) Delete(v interface{}, score int) int {
	if v == nil {
		return 1
	}

	cur := sl.head
	update := [MAX_LEVEL]*skiplistNode{}
	for i := sl.level; i >= 0; i-- {
		update[i] = sl.head
		for cur.forwords[i] != nil {
			if cur.forwords[i].score == score && cur.forwords[i].v == v {
				update[i] = cur
				break
			}
			cur = cur.forwords[i]
		}
	}

	cur = update[0].forwords[0]
	for i := cur.level - 1; i >= 0; i-- {
		if update[i] == sl.head && cur.forwords[i] == nil {
			sl.level = i
		}

		if update[i].forwords[i] == nil {
			update[i].forwords[i] = nil
		} else {
			update[i].forwords[i] = update[i].forwords[i].forwords[i]
		}
	}

	sl.length--
	return 0
}

func (sl *SkipList) String() string {
	return fmt.Sprintf("level:%+v, length:%+v", sl.level, sl.length)
}
