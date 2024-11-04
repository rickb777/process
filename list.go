package process

import "sync"

// list is a simple list safe for concurrent use. The stored items
// must also be safe for concurrent use.
type list[T any] struct {
	list []T
	µ    sync.RWMutex
}

// newList creates a list.
func newList[T any](vals ...T) *list[T] {
	return &list[T]{list: vals}
}

func (list *list[T]) Add(val T) {
	list.µ.Lock()
	defer list.µ.Unlock()

	list.list = append(list.list, val)
}

// Clear empties the list, returning its prior contents as a slice.
func (list *list[T]) Clear() []T {
	list.µ.Lock()
	defer list.µ.Unlock()

	if len(list.list) == 0 {
		return nil
	}

	result := make([]T, len(list.list))
	for i, v := range list.list {
		result[i] = v
	}

	list.list = nil

	return result
}
