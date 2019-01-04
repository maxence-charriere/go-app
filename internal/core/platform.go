package core

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/google/uuid"
)

// PlatformHandler represents the handler that will perform underlying platform
// calls.
type PlatformHandler func(call string) error

// Platform is a struct that provides remote procedure calls from Go to the
// underlying platform.
type Platform struct {
	// The function that performs plaform calls.
	Handler PlatformHandler

	mutex   sync.RWMutex
	returns map[string]chan platformReturn
}

// Call calls the given method with the given input and stores the result in the
// value pointed by the output. It returns an error if the output is not a
// pointer.
func (p *Platform) Call(method string, out interface{}, in interface{}) error {
	returnID := uuid.New().String()

	call, err := json.Marshal(platformCall{
		Method:   method,
		In:       in,
		ReturnID: returnID,
	})
	if err != nil {
		return err
	}

	returnChan := make(chan platformReturn, 1)
	defer close(returnChan)

	p.mutex.Lock()
	if p.returns == nil {
		p.returns = make(map[string]chan platformReturn)
	}

	p.returns[returnID] = returnChan
	p.mutex.Unlock()

	if err := p.Handler(string(call)); err != nil {
		return err
	}

	ret := <-returnChan

	p.mutex.Lock()
	delete(p.returns, returnID)
	p.mutex.Unlock()

	if ret.Err != nil {
		return ret.Err
	}

	if out == nil || len(ret.Out) == 0 {
		return nil
	}

	return json.Unmarshal([]byte(ret.Out), out)
}

// Return returns the given output to the call that waits for the given return
// id.
func (p *Platform) Return(returnID string, out string, err string) {
	p.mutex.RLock()
	returnChan, ok := p.returns[returnID]
	p.mutex.RUnlock()

	if !ok {
		panic(errors.New("no return for " + returnID))
	}

	if len(err) != 0 {
		returnChan <- platformReturn{Err: errors.New(err)}
		return
	}

	returnChan <- platformReturn{Out: out}
}

type platformCall struct {
	Method   string
	In       interface{} `json:",omitempty"`
	ReturnID string
}

type platformReturn struct {
	Out string
	Err error
}
