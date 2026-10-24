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
		head: nil,
		len:  0,
	}
}

// Dequeue remove items from the list
func (l *List) Dequeue() (*Item, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.len == 0 {
		return nil, false
	}
	// The oldest item lives at the tail
	item := l.tail
	l.tail = item.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}
	item.prev = nil
	l.len--
	return item, true
}

// Enqueue add items to the List
func (l *List) Enqueue(val string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	item := &Item{
		val:  val,
		prev: nil,
		next: nil,
	}
	if l.len == 0 {
		l.head = item
		l.tail = item
	} else {
		item.next = l.head
		l.head.prev = item
		l.head = item
	}
	l.len++
}
func (l *List) Peek() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.head == nil {
		return ""
	}
	if l.len == 0 {
		return ""
	}

	return l.tail.val
}

func (l *List) Size() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.len
}
