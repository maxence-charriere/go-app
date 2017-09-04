package bridge

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestPlatformBridge(t *testing.T) {
	var bridge PlatformBridge

	bridge = NewPlatformBridge(func(url string, payload []byte, returnID string) (response []byte, err error) {
		if string(payload) == `"err"` {
			err = errors.New("request error")
			return
		}

		if len(returnID) == 0 {
			response = payload
			return
		}

		retID, err := uuid.Parse(returnID)
		if err != nil {
			t.Fatal(err)
		}

		bridge.Return(retID, payload, nil)
		return
	})

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "request should success",
			test: func(t *testing.T) {
				testPlatformBridgeRequest(t, bridge)
			},
		},
		{
			name: "request with invalid payload should fail",
			test: func(t *testing.T) {
				testPlatformBridgeRequestIvalidPayload(t, bridge)
			},
		},
		{
			name: "request should fail",
			test: func(t *testing.T) {
				testPlatformBridgeRequestFail(t, bridge)
			},
		},
		{
			name: "request with async response should success",
			test: func(t *testing.T) {
				testPlatformBridgeRequestWithAsyncResponse(t, bridge)
			},
		},
		{
			name: "request with async response and invalid payload should fail",
			test: func(t *testing.T) {
				testPlatformBridgeRequestWithAsyncResponseIvalidPayload(t, bridge)
			},
		},
		{
			name: "request with async response should fail",
			test: func(t *testing.T) {
				testPlatformBridgeRequestWithAsyncResponseFail(t, bridge)
			},
		},
		{
			name: "return should panic",
			test: func(t *testing.T) {
				testPlatformBridgeRequestReturnPanic(t, bridge)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testPlatformBridgeRequest(t *testing.T, bridge PlatformBridge) {
	res, err := bridge.Request("", 42)
	if err != nil {
		t.Fatal(err)
	}

	var nb int
	if err = res.Unmarshal(&nb); err != nil {
		t.Fatal(err)
	}
	if nb != 42 {
		t.Fatal("unmarshaled result should be 42:", nb)
	}
}

func testPlatformBridgeRequestIvalidPayload(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.Request("", make(chan int))
	if err == nil {
		t.Fatal("error should not be nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestFail(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.Request("", "err")
	if err == nil {
		t.Fatal("error should not be nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestWithAsyncResponse(t *testing.T, bridge PlatformBridge) {
	res, err := bridge.RequestWithAsyncResponse("", 21)
	if err != nil {
		t.Fatal(err)
	}

	var nb int
	if err = res.Unmarshal(&nb); err != nil {
		t.Fatal(err)
	}
	if nb != 21 {
		t.Fatal("unmarshaled result should be 21:", nb)
	}
}

func testPlatformBridgeRequestWithAsyncResponseIvalidPayload(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.RequestWithAsyncResponse("", make(chan int))
	if err == nil {
		t.Fatal("error should not be nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestWithAsyncResponseFail(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.RequestWithAsyncResponse("", "err")
	if err == nil {
		t.Fatal("error should not be nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestReturnPanic(t *testing.T, bridge PlatformBridge) {
	defer func() { recover() }()
	bridge.Return(uuid.New(), nil, nil)
	t.Fatal("should have panic")
}
