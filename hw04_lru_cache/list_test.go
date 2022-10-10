package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("new list creation", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("list operations", func(t *testing.T) {
		// test PushFront() on empty list
		l := NewList()
		l.PushFront(10)
		require.Equal(t, 1, l.Len())
		require.Equal(t, []interface{}{10}, getListVals(l))

		// test PushBack() on empty list
		l = NewList()
		l.PushBack(10)
		require.Equal(t, 1, l.Len())
		require.Equal(t, []interface{}{10}, getListVals(l))

		// test PushFront(), PushBack() on not empty list
		for i, v := range [...]int{20, 30, 40, 50, 60, 70} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		}
		require.Equal(t, 7, l.Len())
		require.Equal(t, []interface{}{60, 40, 20, 10, 30, 50, 70}, getListVals(l))

		// test Front(), Back()
		require.Equal(t, 60, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		// test Remove() front item
		f := l.Front()
		l.Remove(f)
		require.Equal(t, 6, l.Len())
		require.Equal(t, []interface{}{40, 20, 10, 30, 50, 70}, getListVals(l))

		// test Remove() back item
		b := l.Back()
		l.Remove(b)
		require.Equal(t, 5, l.Len())
		require.Equal(t, []interface{}{40, 20, 10, 30, 50}, getListVals(l))

		// test Remove() middle item
		m := l.Front().Next.Next
		l.Remove(m)
		require.Equal(t, 4, l.Len())
		require.Equal(t, []interface{}{40, 20, 30, 50}, getListVals(l))

		// test MoveToFront front item
		l.MoveToFront(l.Front())
		require.Equal(t, []interface{}{40, 20, 30, 50}, getListVals(l))

		// test MoveToFront back item
		l.MoveToFront(l.Back())
		require.Equal(t, []interface{}{50, 40, 20, 30}, getListVals(l))

		// test MoveToFront middle item
		l.MoveToFront(l.Front().Next)
		require.Equal(t, []interface{}{40, 50, 20, 30}, getListVals(l))
	})
}

func getListVals(l List) []interface{} {
	var vals []interface{}
	li := l.Front()
	for li != nil {
		vals = append(vals, li.Value)
		li = li.Next
	}
	return vals
}
