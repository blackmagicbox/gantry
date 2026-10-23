package queue

import (
	"sync"
)

type Item struct {
	val  string
	next *Item
	prev *Item
}

type List struct {
	mu         sync.Mutex
	head, tail *Item
	len        int
}

func NewList() *List {
	return &List{
		mu:   sync.Mutex{},
		tail: nil,
		head: nil, len: 0}
}

// Dequeue remove items from the list
func (l *List) Dequeue() (*Item, bool) {
	l.mu.Lock()
	defer func() {
		l.len--
		l.mu.Unlock()
	}()
	return nil, false
}

// Enqueue add items to the List
func (l *List) Enqueue(val string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	item := &Item{val: val}
	item.prev = l.head.prev
	item.next = l.head
}
func (l *List) Peek() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Size() == 0 {
		return ""
	}
	if l.head == nil {
		return ""
	}

	return l.tail.prev.val
}

func (l *List) Size() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.len
}
