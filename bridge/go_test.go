package bridge

import (
	"net/url"
	"testing"
)

func TestGoBridge(t *testing.T) {
	uichan := make(chan func(), 256)
	defer close(uichan)

	go func() {
		for f := range uichan {
			f()
		}
	}()

	bridge := NewGoBridge(uichan)

	bridge.Handle("/test", func(url *url.URL, p Payload) (res Payload) {
		res = p
		return
	})

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "handle with invalid pattern should fail",
			test: func(t *testing.T) { testGoBridgeHandleBadPattern(t, bridge) },
		},
		{
			name: "handle with nil handler should fail",
			test: func(t *testing.T) { testGoBridgeHandleNilHandler(t, bridge) },
		},
		{
			name: "request should success",
			test: func(t *testing.T) { testGoBridgeRequest(t, bridge) },
		},
		{
			name: "request with bad URL should fail",
			test: func(t *testing.T) { testGoBridgeRequestBadURL(t, bridge) },
		},
		{
			name: "Request with response should success",
			test: func(t *testing.T) { testGoBridgeRequestWithResponse(t, bridge) },
		},
		{
			name: "request with response and bad URL should fail",
			test: func(t *testing.T) { testGoBridgeRequestWithResponseBadURL(t, bridge) },
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testGoBridgeHandleBadPattern(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()

	bridge.Handle("badpattern", func(url *url.URL, p Payload) (res Payload) {
		return
	})

	t.Fatal("should have panic")
}

func testGoBridgeHandleNilHandler(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()

	bridge.Handle("/test", nil)
	t.Fatal("should have panic")
}

func testGoBridgeRequest(t *testing.T, bridge GoBridge) {
	bridge.Request("/test", NewPayload(42))
}

func testGoBridgeRequestBadURL(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()

	bridge.Request(":{}K{RKVR<<>>!@#", nil)
	t.Fatal("should have panic")
}

func testGoBridgeRequestWithResponse(t *testing.T, bridge GoBridge) {
	res := bridge.RequestWithResponse("/test", NewPayload(21))

	var nb int
	res.Unmarshal(&nb)

	if nb != 21 {
		t.Fatal("response should be 21:", nb)
	}
}

func testGoBridgeRequestWithResponseBadURL(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()
	bridge.RequestWithResponse(":{}K{RKVR<<>>!@#", nil)
	t.Fatal("should have panic")
}

func TestGoBridgeHandleSubpath(t *testing.T) {
	handler := func(url *url.URL, p Payload) (res Payload) {
		return
	}

	b := newGoBridge(make(chan func(), 42))
	b.Handle("/test/foo", handler)

	u, _ := url.Parse("/test/foo")
	b.handle(u, nil)

	u, _ = url.Parse("/test/foo/bar")
	b.handle(u, nil)
}

func TestGoBridgeHandleNotHandled(t *testing.T) {
	defer func() { recover() }()

	u, _ := url.Parse("/nothandled")
	b := newGoBridge(make(chan func(), 42))
	b.handle(u, nil)
	t.Fatal("should have panic")
}
