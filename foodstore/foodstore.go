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
	"fmt"

	"fooddistro/food"
)

// FoodStore holds up to capacity FoodPacks in ascending sorted order (meat
// 'M' before grain/fruit 'B', then by FoodType) to support binary search.
// HIGHLIGHT: replaces the circularque.CircularQue ring buffer with a sorted slice
type FoodStore struct {
	capacity int
	items    []food.FoodPack
}

// HIGHLIGHT END

// NewFoodStore allocates a FoodStore with room for capacity packets. HIGHLIGHT: was NewCircularQue
func NewFoodStore(capacity int) *FoodStore {
	return &FoodStore{
		capacity: capacity,
		items:    make([]food.FoodPack, 0, capacity),
	}
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
func (fs *FoodStore) AcceptMessage(msg food.FoodPack) error {
	if len(fs.items) >= fs.capacity {
		return errors.New("ERROR - Message rejected - storage is full!")
	}

	idx := 0
	for idx < len(fs.items) && lessKey(fs.items[idx], msg) {
		idx++
	}

	fs.items = append(fs.items, food.FoodPack{})
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

// RetrieveDesired removes and returns the packet found by BinarySearchByType,
// selling the last packet in the list instead if the desired type is absent.
// HIGHLIGHT: replaces RetrieveMessage's plain FIFO dequeue
func (fs *FoodStore) RetrieveDesired(desired food.FoodType) (packet food.FoodPack, ok bool, apology string) {
	if len(fs.items) == 0 {
		return food.FoodPack{}, false, ""
	}

	idx, found := BinarySearchByType(fs.items, desired)
	if !found {
		idx = len(fs.items) - 1
		apology = fmt.Sprintf("Sorry, no food packets of the %s are currently available.", desired)
	}

	packet = fs.items[idx]
	fs.items = append(fs.items[:idx], fs.items[idx+1:]...)
	return packet, true, apology
}

// HIGHLIGHT END

// IsEmpty reports whether the store currently holds no packets.
func (fs *FoodStore) IsEmpty() bool {
	return len(fs.items) == 0
}

// IsFull reports whether the store has reached its capacity.
func (fs *FoodStore) IsFull() bool {
	return len(fs.items) >= fs.capacity
}
