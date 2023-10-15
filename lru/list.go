package lru

import "time"

type Entry[K comparable, V any] struct {
	//
	next, prev   *Entry[K, V] //  前驱节点和后继节点
	List         *List[K, V]  // 所属的链表
	Key          K            // 键
	Value        V            // 值
	ExpiresAt    time.Time    // 过期时间
	ExpireBucket uint8        // 过期时间的桶
}

// PrevEntry returns the next list element or nil.
func (e *Entry[K, V]) PrevEntry() *Entry[K, V] {
	if p := e.prev; p != nil && p != &e.List.root {
		return p
	}
	return nil
}

type List[K comparable, V any] struct {
	root Entry[K, V]
	len  int
}

// Init returns an initialized list.
func (l *List[K, V]) Init() *List[K, V] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewList Len returns the number of elements of list l.
func NewList[K comparable, V any]() *List[K, V] { return new(List[K, V]).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[K, V]) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *List[K, V]) Front() *Entry[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List[K, V]) Back() *Entry[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit
func (l *List[K, V]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert
func (l *List[K, V]) insert(e, at *Entry[K, V]) *Entry[K, V] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.List = l
	l.len++
	return e
}

// insertValve
func (l *List[K, V]) insertValue(k K, v V, expiresAt time.Time, at *Entry[K, V]) *Entry[K, V] {
	return l.insert(&Entry[K, V]{Value: v, Key: k, ExpiresAt: expiresAt}, at)
}

func (l *List[K, V]) Remove(e *Entry[K, V]) V {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.List = nil
	l.len--
	return e.Value
}

func (l *List[K, V]) move(e, at *Entry[K, V]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

func (l *List[K, V]) PushFront(k K, v V) *Entry[K, V] {
	l.lazyInit()
	return l.insertValue(k, v, time.Time{}, &l.root)
}

func (l *List[K, V]) PushFrontExplorable(k K, v V, expiresAt time.Time) *Entry[K, V] {
	l.lazyInit()
	return l.insertValue(k, v, expiresAt, &l.root)
}

func (l *List[K, V]) MoveToFront(e *Entry[K, V]) {
	if e.List != l || l.root.next == e {
		return
	}
	l.move(e, &l.root)
}
