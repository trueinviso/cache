package linked_list

import (
	"time"
)

type Node struct {
	Id       int64
	Data     any
	Ts       time.Time
	Prev     *Node
	Next     *Node
	Sentinal bool
}

type LinkedList struct {
	Head    *Node
	Tail    *Node
	Size    int64
	MaxSize int64
}

func NewLinkedList(maxSize int64) *LinkedList {
	head := Node{Sentinal: true}
	Tail := Node{Sentinal: true}
	head.Next = &Tail
	Tail.Prev = &head

	return &LinkedList{
		Head:    &head,
		Tail:    &Tail,
		MaxSize: maxSize,
		Size:    0,
	}
}

func (l *LinkedList) MoveToFront(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	node.Ts = time.Now()

	l.prepend(node)
}

func (l *LinkedList) Insert(node *Node) {
	node.Ts = time.Now()
	l.prepend(node)
	l.Size++
}

func (l *LinkedList) Delete(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	node.Next = nil
	node.Prev = nil
	l.Size--
}

func (l *LinkedList) Pop() int64 {
	if l.Size == 0 {
		return 0
	}

	prev := l.Tail.Prev
	prev.Prev.Next = l.Tail
	l.Tail.Prev = prev.Prev
	l.Size--

	return prev.Id
}

func (l *LinkedList) prepend(node *Node) {
	node.Prev = l.Head
	l.Head.Next.Prev = node
	node.Next = l.Head.Next
	l.Head.Next = node
}
