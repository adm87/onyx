package storage

type SlotKey struct {
	index uint
	gen   uint64
}

func (k SlotKey) Index() uint {
	return k.index
}

type slot[T any] struct {
	data T
	gen  uint64
	used bool
}

type SlotMap[T any] struct {
	slots     []slot[T]
	freeSlots []int
}

func NewSlotMap[T any](capacity int) *SlotMap[T] {
	slots := make([]slot[T], 0, capacity)
	freeSlots := make([]int, 0, capacity)
	for i := range capacity {
		freeSlots = append(freeSlots, i)
		slots = append(slots, slot[T]{gen: 0, used: false})
	}
	return &SlotMap[T]{
		slots:     slots,
		freeSlots: freeSlots,
	}
}

func (m *SlotMap[T]) Insert(value T) (SlotKey, bool) {
	if len(m.freeSlots) == 0 {
		return m.append(value)
	}

	index := m.freeSlots[len(m.freeSlots)-1]
	m.freeSlots = m.freeSlots[:len(m.freeSlots)-1]

	slot := &m.slots[index]
	slot.data = value
	slot.used = true

	return SlotKey{index: uint(index), gen: slot.gen}, true
}

func (m *SlotMap[T]) Remove(key SlotKey) bool {
	if key.index >= uint(len(m.slots)) {
		return false
	}

	slot := &m.slots[key.index]
	if !slot.used || slot.gen != key.gen {
		return false
	}

	slot.used = false
	slot.gen++

	m.freeSlots = append(m.freeSlots, int(key.index))
	return true
}

func (m *SlotMap[T]) Get(key SlotKey) (T, bool) {
	if key.index >= uint(len(m.slots)) {
		var zero T
		return zero, false
	}

	slot := &m.slots[key.index]
	if !slot.used || slot.gen != key.gen {
		var zero T
		return zero, false
	}

	return slot.data, true
}

func (m *SlotMap[T]) Set(key SlotKey, value T) bool {
	if key.index >= uint(len(m.slots)) {
		return false
	}

	slot := &m.slots[key.index]
	if !slot.used || slot.gen != key.gen {
		return false
	}

	slot.data = value
	return true
}

func (m *SlotMap[T]) Has(key SlotKey) bool {
	if key.index >= uint(len(m.slots)) {
		return false
	}

	slot := &m.slots[key.index]
	return slot.used && slot.gen == key.gen
}

func (m *SlotMap[T]) Len() int {
	return len(m.slots) - len(m.freeSlots)
}

func (m *SlotMap[T]) Capacity() int {
	return cap(m.slots)
}

func (m *SlotMap[T]) Clear() {
	m.freeSlots = m.freeSlots[:0]
	for i := range m.slots {
		m.slots[i].used = false
		m.slots[i].gen++
		m.freeSlots = append(m.freeSlots, i)
	}
}

func (m *SlotMap[T]) Each(f func(key SlotKey, value T) bool) {
	for i := range m.slots {
		slot := &m.slots[i]
		if slot.used {
			key := SlotKey{index: uint(i), gen: slot.gen}
			if !f(key, slot.data) {
				return
			}
		}
	}
}

func (m *SlotMap[T]) append(value T) (SlotKey, bool) {
	index := len(m.slots)
	m.slots = append(m.slots, slot[T]{data: value, gen: 0, used: true})
	return SlotKey{index: uint(index), gen: m.slots[index].gen}, true
}
