package queue

import "errors"

type Item struct {
	val  string
	next *Item
}

// Dequeue remove items from the list
func (l *Item) Dequeue() (*Item, error) {
	return nil, errors.New("delist not implemented")
}

// Enqueue add items to the List
func (l *Item) Enqueue(val string) error {
	return errors.New("enqueue not implemented")
}
