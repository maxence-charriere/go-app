package bridge

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/pkg/errors"
)

// PlatformRPCHandler represents the handler that will perform platform calls.
type PlatformRPCHandler func(call string) error

// PlatformRPC is a struct that implements the remote procedure call from  Go to an
// underlying platform.
type PlatformRPC struct {
	Handler PlatformRPCHandler

	mutex   sync.RWMutex
	returns map[string]chan rpcReturn
}

// Call calls the given method with the given input and stores the result in
// the value pointed by the output.
// It returns an error if the output is not a pointer.
func (r *PlatformRPC) Call(method string, out interface{}, in interface{}) error {
	returnID := uuid.New().String()

	call, err := json.Marshal(PlatformCall{
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

	if len(rpcRet.Output) != 0 {
		return json.Unmarshal([]byte(rpcRet.Output), out)
	}
	return nil
}

// Return returns the given output to the call that waits for the given return
// id.
func (r *PlatformRPC) Return(retID string, out string, errString string) {
	r.mutex.RLock()
	rpcRetC, ok := r.returns[retID]
	r.mutex.RUnlock()

	if !ok {
		app.Panic("no async call for " + retID)
	}

	var err error
	if len(errString) != 0 {
		err = errors.New(errString)
	}

	rpcRetC <- rpcReturn{
		Output: out,
		Error:  err,
	}
}

// PlatformCall represents a platform call.
type PlatformCall struct {
	Method   string
	Input    interface{} `json:",omitempty"`
	ReturnID string
}

type rpcReturn struct {
	Output string
	Error  error
}
