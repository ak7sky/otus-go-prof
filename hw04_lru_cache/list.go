package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	i := &ListItem{Value: v}
	if l.len == 0 {
		l.front, l.back = i, i
	} else {
		l.front.Prev = i
		i.Next = l.front
		l.front = i
	}
	l.len++
	return i
}

func (l *list) PushBack(v interface{}) *ListItem {
	i := &ListItem{Value: v}
	if l.len == 0 {
		l.front, l.back = i, i
	} else {
		l.back.Next = i
		i.Prev = l.back
		l.back = i
	}
	l.len++
	return i
}

// Remove existed(!!!) ListItem from list.
func (l *list) Remove(i *ListItem) {
	switch {
	case l.len == 1:
		l.front, l.back = nil, nil
	case l.len > 1 && l.Front() == i:
		l.front = i.Next
		i.Next, i.Next.Prev = nil, nil
	case l.len > 1 && l.Back() == i:
		l.back = i.Prev
		i.Prev, i.Prev.Next = nil, nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Prev, i.Next = nil, nil
	}
	l.len--
}

// MoveToFront existed(!!!) ListItem.
func (l *list) MoveToFront(i *ListItem) {
	if l.front != i {
		if l.back == i {
			i.Prev.Next = nil
			l.back = i.Prev
		} else {
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		}
		i.Prev = nil
		i.Next = l.front
		l.front.Prev = i
		l.front = i
	}
}

func NewList() List {
	return new(list)
}
