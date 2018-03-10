package bridge

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// AsyncInput is the interface that describes an asynchronous call input.
type AsyncInput interface {
	// Async reports whether an call is asynchrous.
	Async() bool
}

// CallHandler represents the handler that will perform the call.
type CallHandler func(call string) (ret string, err error)

// RPC is a stuct that implements the remote procedure call mechanismes between
// Go and an underlying platform.
type RPC struct {
	CallHandler CallHandler

	mutex   sync.RWMutex
	returns map[string]chan asyncReturn
}

// Call calls the given method with the given input and stores the result in
// the value pointed by the output.
// It returns an error if the output is not a pointer.
func (r *RPC) Call(method string, in interface{}, out interface{}) error {
	var returnID string
	if asyncIn, ok := in.(AsyncInput); ok && asyncIn.Async() {
		returnID = uuid.New().String()
	}

	call, err := json.Marshal(call{
		Method:   method,
		Input:    in,
		ReturnID: returnID,
	})
	if err != nil {
		return err
	}

	var ret string
	if ret, err = r.CallHandler(string(call)); err != nil {
		return err
	}

	if len(returnID) != 0 {
		retC := make(chan asyncReturn)

		r.mutex.Lock()
		if r.returns == nil {
			r.returns = make(map[string]chan asyncReturn)
		}
		r.returns[returnID] = retC
		r.mutex.Unlock()

		asyncRet := <-retC

		r.mutex.Lock()
		close(retC)
		delete(r.returns, returnID)
		r.mutex.Unlock()

		if ret, err = asyncRet.Return, asyncRet.Error; err != nil {
			return err
		}
	}

	return json.Unmarshal([]byte(ret), out)
}

// Return returns the given output to the asynchrounous call that waits for the
// given return id.
func (r *RPC) Return(retID string, ret string, err string) error {
	r.mutex.RLock()
	retC, ok := r.returns[retID]
	r.mutex.RUnlock()

	if !ok {
		return errors.New("no async call for " + retID)
	}

	retC <- asyncReturn{
		Return: ret,
		Error:  errors.New(err),
	}
	return nil
}

type call struct {
	Method   string
	Input    interface{}
	ReturnID string `json:",omitempty"`
}

type asyncReturn struct {
	Return string
	Error  error
}
