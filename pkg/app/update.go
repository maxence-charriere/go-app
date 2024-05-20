package app

// updateManager manages scheduled updates for UI components. It ensures that
// components are updated in an order respecting their depth in the UI
// hierarchy.
type updateManager struct {
	pending []map[Composer]int
}

// Add queues a component for an update and increments its associated counter by
// the given value. The component will be marked for update if its counter
// becomes greater than 0.
func (m *updateManager) Add(c Composer, v int) {
	depth := int(c.depth())
	if len(m.pending) <= depth {
		size := max(depth+1, 100)
		pending := make([]map[Composer]int, size)
		copy(pending, m.pending)
		m.pending = pending
	}

	updates := m.pending[depth]
	if updates == nil {
		updates = make(map[Composer]int)
		m.pending[depth] = updates
	}
	updates[c] += v
}

// Done removes the given component from the update queue, marking it as updated.
func (m *updateManager) Done(v Composer) {
	depth := v.depth()
	if len(m.pending) <= int(depth) {
		return
	}
	updates := m.pending[depth]
	delete(updates, v)
}

// UpdateForEach iterates over all components queued for updates via the Add
// method, executing a specified function on each component with an associated
// counter greater than 0. After the function is invoked on a component, its
// entry is removed from the queue.
//
// This method ensures actions are taken only on components ready for an update
// (counter > 0) and maintains the queue's cleanliness by removing components
// once processed.
func (m *updateManager) UpdateForEach(do func(Composer)) {
	for _, updates := range m.pending {
		for compo, counter := range updates {
			if counter > 0 {
				do(compo)
			}
			delete(updates, compo)
		}
	}
}
