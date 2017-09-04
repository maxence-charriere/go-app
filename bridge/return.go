package bridge

import (
	"sync"

	"github.com/google/uuid"
)

type returnPayload struct {
	response Payload
	err      error
}

type returnRegistry struct {
	mutex   sync.Mutex
	returns map[uuid.UUID]chan returnPayload
}

func newReturnRegistry() *returnRegistry {
	return &returnRegistry{
		returns: make(map[uuid.UUID]chan returnPayload),
	}
}

func (r *returnRegistry) Set(id uuid.UUID, retchan chan returnPayload) {
	r.mutex.Lock()
	r.returns[id] = retchan
	r.mutex.Unlock()
}

func (r *returnRegistry) Get(id uuid.UUID) (retchan chan returnPayload, ok bool) {
	r.mutex.Lock()
	retchan, ok = r.returns[id]
	r.mutex.Unlock()
	return
}

func (r *returnRegistry) Delete(id uuid.UUID) {
	r.mutex.Lock()
	delete(r.returns, id)
	r.mutex.Unlock()
}
