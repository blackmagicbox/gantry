package queue

import "testing"

func TestListEnqueue(t *testing.T) {
	l := NewList()
	l.Enqueue("player1")
	if l.Size() != 1 {
		t.Fatalf("expected list size to be 1, got %v", l.Size())
	}
	if v := l.Peek(); v != "player1" {
		t.Fatalf("expected player1 at the front, got %q", v)
	}
}

func TestListPeekReturnsOldest(t *testing.T) {
	l := NewList()
	l.Enqueue("player1")
	l.Enqueue("player2")

	if v := l.Peek(); v != "player1" {
		t.Fatalf("expected Peek to return oldest item player1, got %q", v)
	}
	if l.Size() != 2 {
		t.Fatalf("expected Peek not to change size, got %v", l.Size())
	}
}

func TestListDequeueFIFOOrder(t *testing.T) {
	l := NewList()
	want := []string{"player1", "player2", "player3"}
	for _, p := range want {
		l.Enqueue(p)
	}

	for i, w := range want {
		item, ok := l.Dequeue()
		if !ok {
			t.Fatalf("dequeue %d: expected ok, got false", i)
		}
		if item.val != w {
			t.Fatalf("dequeue %d: expected %q, got %q", i, w, item.val)
		}
		if got, wantSize := l.Size(), len(want)-i-1; got != wantSize {
			t.Fatalf("dequeue %d: expected list size to be %v, got %v", i, wantSize, got)
		}
	}
}

func TestListDequeueEmpty(t *testing.T) {
	l := NewList()
	if _, ok := l.Dequeue(); ok {
		t.Fatalf("expected Dequeue on empty list to return false")
	}
	if v := l.Peek(); v != "" {
		t.Fatalf("expected Peek on empty list to return empty string, got %q", v)
	}
}

func TestListReusableAfterDrain(t *testing.T) {
	l := NewList()
	l.Enqueue("player1")
	l.Dequeue()

	l.Enqueue("player2")
	l.Enqueue("player3")

	item, ok := l.Dequeue()
	if !ok || item.val != "player2" {
		t.Fatalf("expected player2 after draining and refilling, got %v (ok=%v)", item, ok)
	}
}
