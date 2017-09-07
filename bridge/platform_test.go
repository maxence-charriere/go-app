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
			name: "request should success",
			test: func(t *testing.T) {
				testPlatformBridgeRequest(t, bridge)
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
			name: "request with async response should fail",
			test: func(t *testing.T) {
				testPlatformBridgeRequestWithAsyncResponseFail(t, bridge)
			},
		},
		{
			name: "return with ivalid id should panic",
			test: func(t *testing.T) {
				testPlatformBridgeRequestReturnIvalidID(t, bridge)
			},
		},
		{
			name: "unset return should panic",
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
		t.Fatal("unmarshaled result should be 42:", nb)
	}
}

func testPlatformBridgeRequestFail(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.Request("/error", nil)
	if err == nil {
		t.Fatal("error should not be nil")
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
		t.Fatal("unmarshaled result should be 21:", nb)
	}
}

func testPlatformBridgeRequestWithAsyncResponseFail(t *testing.T, bridge PlatformBridge) {
	_, err := bridge.RequestWithAsyncResponse("/error", nil)
	if err == nil {
		t.Fatal("error should not be nil")
	}
	t.Log(err)
}

func testPlatformBridgeRequestReturnIvalidID(t *testing.T, bridge PlatformBridge) {
	defer func() { recover() }()
	bridge.Return("whoisyourdaddy", nil, nil)
	t.Fatal("should have panic")
}

func testPlatformBridgeRequestReturnUnset(t *testing.T, bridge PlatformBridge) {
	defer func() { recover() }()
	bridge.Return(uuid.New().String(), nil, nil)
	t.Fatal("should have panic")
}
