package storage

import "math"

// SlotKey is a unique identifier for a slot in the SlotMap.
// It encodes both the index of the slot and a generation counter to prevent stale references.
type SlotKey uint64

func newSlotKey(index uint32, gen uint32) SlotKey {
	return SlotKey(uint64(index) | uint64(gen)<<32)
}

// Index returns the index part of the SlotKey.
func (k SlotKey) Index() uint32 {
	return uint32(uint64(k) & (1<<32 - 1))
}

// Gen returns the generation part of the SlotKey.
func (k SlotKey) Gen() uint32 {
	return uint32(uint64(k) >> 32)
}

type slot[T any] struct {
	data T
	gen  uint32
	used bool
}

// SlotMap is a data structure that manages a collection of slots, allowing for efficient insertion, removal, and retrieval of values.
// The map will grow dynamically as needed, and it uses a free list to reuse slots that have been removed.
type SlotMap[T any] struct {
	slots     []slot[T]
	freeSlots []uint32
}

// NewSlotMap creates a new SlotMap with the specified initial capacity.
func NewSlotMap[T any](capacity int) *SlotMap[T] {
	slots := make([]slot[T], 0, capacity)
	freeSlots := make([]uint32, 0, capacity)
	for i := range capacity {
		freeSlots = append(freeSlots, uint32(i))
		slots = append(slots, slot[T]{gen: 0, used: false})
	}
	return &SlotMap[T]{
		slots:     slots,
		freeSlots: freeSlots,
	}
}

// Insert adds a new value to the SlotMap and returns a SlotKey that can be used to retrieve it.
// If the SlotMap is full, it will automatically grow to accommodate more entries.
func (m *SlotMap[T]) Insert(value T) (SlotKey, bool) {
	if len(m.freeSlots) == 0 {
		return m.append(value)
	}

	index := m.freeSlots[len(m.freeSlots)-1]
	m.freeSlots = m.freeSlots[:len(m.freeSlots)-1]

	slot := &m.slots[index]
	slot.data = value
	slot.used = true

	return newSlotKey(uint32(index), slot.gen), true
}

// Remove deletes the value associated with the given SlotKey from the SlotMap.
// It returns true if the value was successfully removed, or false if the key was invalid or already removed.
func (m *SlotMap[T]) Remove(key SlotKey) bool {
	if key.Index() >= uint32(len(m.slots)) {
		return false
	}

	slot := &m.slots[key.Index()]
	if !slot.used || slot.gen != key.Gen() {
		return false
	}

	slot.used = false
	slot.gen++

	m.freeSlots = append(m.freeSlots, key.Index())
	return true
}

// Get retrieves the value associated with the given SlotKey.
// It returns the value and true if the key is valid, or the zero value and false if the key is invalid or removed.
func (m *SlotMap[T]) Get(key SlotKey) (T, bool) {
	if key.Index() >= uint32(len(m.slots)) {
		var zero T
		return zero, false
	}

	slot := &m.slots[key.Index()]
	if !slot.used || slot.gen != key.Gen() {
		var zero T
		return zero, false
	}

	return slot.data, true
}

// Set updates the value associated with the given SlotKey.
// It returns true if the value was successfully updated, or false if the key was invalid or removed.
func (m *SlotMap[T]) Set(key SlotKey, value T) bool {
	if key.Index() >= uint32(len(m.slots)) {
		return false
	}

	slot := &m.slots[key.Index()]
	if !slot.used || slot.gen != key.Gen() {
		return false
	}

	slot.data = value
	return true
}

// Has checks if the given SlotKey is valid and the slot is currently in use.
// It returns true if the key is valid and the slot is used, or false otherwise.
func (m *SlotMap[T]) Has(key SlotKey) bool {
	if key.Index() >= uint32(len(m.slots)) {
		return false
	}

	slot := &m.slots[key.Index()]
	return slot.used && slot.gen == key.Gen()
}

// Len returns the number of currently used slots in the SlotMap.
func (m *SlotMap[T]) Len() int {
	return len(m.slots) - len(m.freeSlots)
}

// Capacity returns the total capacity of the SlotMap, including both used and free slots.
func (m *SlotMap[T]) Capacity() int {
	return cap(m.slots)
}

// Clear removes all entries from the SlotMap and resets it to its initial state.
func (m *SlotMap[T]) Clear() {
	m.freeSlots = m.freeSlots[:0]
	for i := range m.slots {
		m.slots[i].used = false
		m.slots[i].gen++
		m.freeSlots = append(m.freeSlots, uint32(i))
	}
}

// Each iterates over all used slots in the SlotMap and calls the provided function with the SlotKey and value of each slot.
// If the function returns false, the iteration will stop early.
func (m *SlotMap[T]) Each(f func(key SlotKey, value T) bool) {
	for i := range m.slots {
		slot := &m.slots[i]
		if slot.used {
			key := newSlotKey(uint32(i), slot.gen)
			if !f(key, slot.data) {
				return
			}
		}
	}
}

func (m *SlotMap[T]) append(value T) (SlotKey, bool) {
	index := uint32(len(m.slots))
	if index == math.MaxUint32 {
		// Realistically, this should never happen in practice for a typical game, but we'll check to prevent overflow.
		return 0, false
	}
	m.slots = append(m.slots, slot[T]{data: value, gen: 0, used: true})
	return newSlotKey(index, m.slots[index].gen), true
}
