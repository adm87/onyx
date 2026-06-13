package slotmap

// SlotMap is a data structure that provides efficient storage and retrieval of values using unique IDs.
type SlotMap[T any] struct {
	data  []T      // indexable data storage
	slots []uint64 // indexable slot storage, contains packed index and generation
	free  []int    // stack of free slot indices
}

func unpack(id uint64) (index int, generation uint32) {
	return int(id & 0xFFFFFFFF), uint32(id >> 32)
}

func pack(index int, generation uint32) uint64 {
	return (uint64(generation) << 32) | uint64(index)
}

// New creates a new SlotMap with an optional initial seed size. The seed determines the initial capacity of the SlotMap.
func New[T any](seed int) *SlotMap[T] {
	seed = max(seed, 0)

	slotMap := &SlotMap[T]{
		data:  make([]T, seed),
		slots: make([]uint64, seed),
		free:  make([]int, 0, seed),
	}

	for i := seed - 1; i >= 0; i-- {
		slotMap.slots[i] = pack(i, 1)
		slotMap.free = append(slotMap.free, i)
	}

	return slotMap
}

// Insert adds a new value to the SlotMap and returns a unique ID for that value.
// It reuses free slots if available, otherwise it appends a new slot.
func (s *SlotMap[T]) Insert(value T) uint64 {
	if len(s.free) == 0 {
		return s.append(value)
	}

	idx := s.free[len(s.free)-1]
	s.free = s.free[:len(s.free)-1]

	id := s.slots[idx]
	s.data[idx] = value

	return id
}

// Get retrieves the value associated with the given ID.
// It returns the value and a boolean indicating whether the retrieval was successful.
func (s *SlotMap[T]) Get(id uint64) (T, bool) {
	reqIdx, reqGen := unpack(id)

	if reqIdx >= len(s.data) {
		var zero T
		return zero, false
	}

	idx, gen := unpack(s.slots[reqIdx])

	if idx != reqIdx || gen != reqGen {
		var zero T
		return zero, false
	}

	return s.data[idx], true
}

func (s *SlotMap[T]) Set(id uint64, value T) (oldValue T, ok bool) {
	idx, gen := unpack(id)

	if idx >= len(s.data) {
		return
	}

	currIdx, currGen := unpack(s.slots[idx])

	if currIdx != idx || currGen != gen {
		return
	}

	oldValue = s.data[idx]
	s.data[idx] = value

	return oldValue, true
}

// Delete removes the value associated with the given ID from the SlotMap.
// It returns the deleted value and a boolean indicating whether the deletion was successful.
func (s *SlotMap[T]) Delete(id uint64) (T, bool) {
	reqIdx, reqGen := unpack(id)

	if reqIdx == 0 {
		var zero T
		return zero, false
	}

	if reqIdx >= len(s.data) {
		var zero T
		return zero, false
	}

	idx, gen := unpack(s.slots[reqIdx])

	if idx != reqIdx || gen != reqGen {
		var zero T
		return zero, false
	}

	s.slots[reqIdx] = pack(idx, gen+1)
	s.free = append(s.free, int(reqIdx))

	var zero T

	data := s.data[idx]
	s.data[idx] = zero

	return data, true
}

func (s *SlotMap[T]) append(value T) uint64 {
	id := pack(len(s.data), 1)
	s.data = append(s.data, value)
	s.slots = append(s.slots, id)
	return id
}
