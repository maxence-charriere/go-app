package bridge

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Handler represents the handler that will perform the call.
type Handler func(call string) error

// RPC is a stuct that implements the remote procedure call mechanismes between
// Go and an underlying platform.
type RPC struct {
	Handler Handler

	mutex   sync.RWMutex
	returns map[string]chan rpcReturn
}

// Call calls the given method with the given input and stores the result in
// the value pointed by the output.
// It returns an error if the output is not a pointer.
func (r *RPC) Call(method string, in interface{}, out interface{}) error {
	returnID := uuid.New().String()

	call, err := json.Marshal(call{
		Method:   method,
		Input:    in,
		ReturnID: returnID,
	})
	if err != nil {
		return err
	}

	rpcRetC := make(chan rpcReturn, 1)

	r.mutex.Lock()
	if r.returns == nil {
		r.returns = make(map[string]chan rpcReturn)
	}
	r.returns[returnID] = rpcRetC
	r.mutex.Unlock()

	if err = r.Handler(string(call)); err != nil {
		return err
	}

	rpcRet := <-rpcRetC

	r.mutex.Lock()
	delete(r.returns, returnID)
	close(rpcRetC)
	r.mutex.Unlock()

	if rpcRet.Error != nil {
		return rpcRet.Error
	}

	return json.Unmarshal([]byte(rpcRet.Return), out)
}

// Return returns the given output to the asynchrounous call that waits for the
// given return id.
func (r *RPC) Return(retID string, ret string, errString string) {
	r.mutex.RLock()
	rpcRetC, ok := r.returns[retID]
	r.mutex.RUnlock()

	if !ok {
		panic("no async call for " + retID)
	}

	var err error
	if len(errString) != 0 {
		err = errors.New(errString)
	}

	rpcRetC <- rpcReturn{
		Return: ret,
		Error:  err,
	}
}

type call struct {
	Method   string
	Input    interface{}
	ReturnID string `json:",omitempty"`
}

type rpcReturn struct {
	Return string
	Error  error
}
