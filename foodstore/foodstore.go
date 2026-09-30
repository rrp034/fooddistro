// HIGHLIGHT: in file foodstore/foodstore.go; replaces circularque/circularque.go.
// Storage is now a linear list (sorted slice) instead of a circular queue,
// and retrieval uses a binary search instead of sequential FIFO removal.
//
// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, foodstore, and sales packages.

package foodstore

import (
	"errors"

	"fooddistro/food"
)

// FoodStore holds up to capacity values in ascending order using the supplied
// comparison function.
// HIGHLIGHT: replaces the circularque.CircularQue ring buffer with a generic sorted slice
type FoodStore[T any] struct {
	capacity int
	items    []T
	less     func(T, T) bool
}

// HIGHLIGHT END

// NewFoodStore allocates a generic sorted store with room for capacity values.
// HIGHLIGHT: was NewCircularQue and now accepts any value type
func NewFoodStore[T any](capacity int, less func(T, T) bool) *FoodStore[T] {
	return &FoodStore[T]{
		capacity: capacity,
		items:    make([]T, 0, capacity),
		less:     less,
	}
}

// HIGHLIGHT END

// NewFoodPackStore creates a sorted store using the FoodPack ordering required
// by the food distribution system.
// HIGHLIGHT: provides the domain-specific adapter for the generic store
func NewFoodPackStore(capacity int) *FoodStore[food.FoodPack] {
	return NewFoodStore(capacity, lessKey)
}

// HIGHLIGHT END

// lessKey reports whether food pack a sorts strictly before food pack b.
// HIGHLIGHT: new helper needed to keep the list in sorted order
func lessKey(a, b food.FoodPack) bool {
	if a.FoodShipment != b.FoodShipment {
		return a.FoodShipment == 'M' // meat ('M') always precedes grain/fruit ('B')
	}
	return a.FoodType < b.FoodType
}

// HIGHLIGHT END

// AcceptMessage inserts msg into its sorted position (no sort package used).
// HIGHLIGHT: was a plain ring-buffer insert; now a sorted insert into a slice
func (fs *FoodStore[T]) AcceptMessage(msg T) error {
	if len(fs.items) >= fs.capacity {
		return errors.New("ERROR - Message rejected - storage is full!")
	}

	idx := 0
	for idx < len(fs.items) && fs.less(fs.items[idx], msg) {
		idx++
	}

	var zero T
	fs.items = append(fs.items, zero)
	copy(fs.items[idx+1:], fs.items[idx:])
	fs.items[idx] = msg
	return nil
}

// HIGHLIGHT END

// BinarySearchByType locates the desired FoodType in the sorted items slice.
// HIGHLIGHT: new function, the required binary search for the "B" grading option
func BinarySearchByType(items []food.FoodPack, desired food.FoodType) (int, bool) {
	desiredShipment := byte('M')
	if food.IsGrainVegetable(desired) {
		desiredShipment = 'B'
	}
	target := food.FoodPack{FoodType: desired, FoodShipment: desiredShipment}

	lo, hi := 0, len(items)
	for lo < hi {
		mid := (lo + hi) / 2
		if lessKey(items[mid], target) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	if lo < len(items) && items[lo].FoodType == desired {
		return lo, true
	}
	return lo, false
}

// HIGHLIGHT END

// RemoveAt removes and returns the value at index while preserving slice order.
// HIGHLIGHT: generic removal replaces the circular queue dequeue operation
func (fs *FoodStore[T]) RemoveAt(index int) (value T, ok bool) {
	if index < 0 || index >= len(fs.items) {
		return value, false
	}

	value = fs.items[index]
	copy(fs.items[index:], fs.items[index+1:])
	fs.items = fs.items[:len(fs.items)-1]
	return value, true
}

// HIGHLIGHT END

// Items returns the values currently held by the store in sorted order.
func (fs *FoodStore[T]) Items() []T {
	return fs.items
}

// IsEmpty reports whether the store currently holds no values.
func (fs *FoodStore[T]) IsEmpty() bool {
	return len(fs.items) == 0
}

// IsFull reports whether the store has reached its capacity.
func (fs *FoodStore[T]) IsFull() bool {
	return len(fs.items) >= fs.capacity
}
