package array_linked_list

import "time"

type Node struct {
	Id       int64
	Data     any
	Ts       time.Time
	Prev     int64
	Next     int64
	Sentinal bool
}

type LinkedList struct {
	Head     int64
	Tail     int64
	Size     int64
	MaxSize  int64
	NodeList []Node
	FreeSlot int64
}

func NewLinkedList(maxSize int64) *LinkedList {
	head := Node{Sentinal: true}
	tail := Node{Sentinal: true}
	head.Next = 1
	tail.Prev = 0
	nodeList := make([]Node, int(maxSize))
	nodeList[0] = head
	nodeList[1] = tail

	return &LinkedList{
		Head:     0,
		Tail:     1,
		MaxSize:  maxSize,
		Size:     0,
		NodeList: nodeList,
		FreeSlot: 2,
	}
}

func (l *LinkedList) MoveToFront(idx int64) {
	prevIdx := l.NodeList[idx].Prev
	nextIdx := l.NodeList[idx].Next

	l.NodeList[prevIdx].Next = nextIdx
	l.NodeList[nextIdx].Prev = prevIdx
	l.NodeList[idx].Ts = time.Now()

	l.prepend(idx)
}

func (l *LinkedList) Insert(node *Node) int64 {
	node.Ts = time.Now()
	l.NodeList[l.FreeSlot] = *node
	l.prepend(l.FreeSlot)
	l.FreeSlot++
	l.Size++

	return l.FreeSlot - 1
}

func (l *LinkedList) AddToFullCache(node *Node) (idx int64, oldKey int64) {
	removeIdx := l.NodeList[l.Tail].Prev
	oldKey = l.NodeList[removeIdx].Id
	prevIdx := l.NodeList[removeIdx].Prev
	node.Prev = prevIdx
	node.Next = l.Tail
	l.NodeList[removeIdx] = *node
	l.MoveToFront(removeIdx)
	return removeIdx, oldKey
}

func (l *LinkedList) prepend(idx int64) {
	headIdx := l.Head
	oldFirstIdx := l.NodeList[headIdx].Next

	// Point the new node to Head and the old first node
	l.NodeList[idx].Prev = headIdx
	l.NodeList[idx].Next = oldFirstIdx

	// Update Head and the old first node to point back to the new node
	l.NodeList[headIdx].Next = idx
	l.NodeList[oldFirstIdx].Prev = idx
}
