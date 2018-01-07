package bridge

import (
	"errors"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func TestPlatformBridge(t *testing.T) {
	var bridge PlatformBridge

	bridge = NewPlatformBridge(func(rawurl string, p Payload, returnID string) (res Payload, err error) {
		u, err := url.Parse(rawurl)
		if err != nil {
			t.Fatal(err)
		}

		if u.Path == "/error" {
			err = errors.New("fake error")
			return
		}

		if len(returnID) == 0 {
			res = p
			return
		}

		bridge.Return(returnID, p, nil)
		return
	})

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "send a request",
			test: func(t *testing.T) {
				testPlatformBridgeRequest(t, bridge)
			},
		},
		{
			name: "request returns an error",
			test: func(t *testing.T) {
				testPlatformBridgeRequestFail(t, bridge)
			},
		},
		{
			name: "send a request with asynchronous response",
			test: func(t *testing.T) {
				testPlatformBridgeRequestWithAsyncResponse(t, bridge)
			},
		},
		{
			name: "request with asynchronous response returns an error",
			test: func(t *testing.T) {
				testPlatformBridgeRequestWithAsyncResponseFail(t, bridge)
			},
		},
		{
			name: "return with invalid id panics",
			test: func(t *testing.T) {
				testPlatformBridgeRequestReturnIvalidID(t, bridge)
			},
		},
		{
			name: "return not set up panics",
			test: func(t *testing.T) {
				testPlatformBridgeRequestReturnUnset(t, bridge)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testPlatformBridgeRequest(t *testing.T, bridge PlatformBridge) {
	res, err := bridge.Request("", NewPayload(42))
	if err != nil {
		t.Fatal(err)
	}

	var nb int
	res.Unmarshal(&nb)

	if nb != 42 {
		t.Fatal("unmarshaled result is not 42:", nb)
	}
}

func testPlatformBridgeRequestFail(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.Request("/error", nil)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestWithAsyncResponse(t *testing.T, bridge PlatformBridge) {
	res, err := bridge.RequestWithAsyncResponse("", NewPayload(21))
	if err != nil {
		t.Fatal(err)
	}

	var nb int
	res.Unmarshal(&nb)

	if nb != 21 {
		t.Fatal("unmarshaled result is not 21:", nb)
	}
}

func testPlatformBridgeRequestWithAsyncResponseFail(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.RequestWithAsyncResponse("/error", nil)
	if err == nil {
		t.Fatal("error is nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestReturnIvalidID(t *testing.T, bridge PlatformBridge) {
	defer func() { recover() }()
	bridge.Return("whoisyourdaddy", nil, nil)
	t.Fatal("no panic")
}

func testPlatformBridgeRequestReturnUnset(t *testing.T, bridge PlatformBridge) {
	defer func() { recover() }()
	bridge.Return(uuid.New().String(), nil, nil)
	t.Fatal("no panic")
}
