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
			name: "handle with invalid pattern panics",
			test: func(t *testing.T) { testGoBridgeHandleBadPattern(t, bridge) },
		},
		{
			name: "handle with nil handler panics",
			test: func(t *testing.T) { testGoBridgeHandleNilHandler(t, bridge) },
		},
		{
			name: "send a request",
			test: func(t *testing.T) { testGoBridgeRequest(t, bridge) },
		},
		{
			name: "send a request with bad URL panics",
			test: func(t *testing.T) { testGoBridgeRequestBadURL(t, bridge) },
		},
		{
			name: "send a request with response",
			test: func(t *testing.T) { testGoBridgeRequestWithResponse(t, bridge) },
		},
		{
			name: "send a request with response and a bad URL panics",
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

	t.Fatal("no panic")
}

func testGoBridgeHandleNilHandler(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()

	bridge.Handle("/test", nil)
	t.Fatal("no panic")
}

func testGoBridgeRequest(t *testing.T, bridge GoBridge) {
	bridge.Request("/test", NewPayload(42))
}

func testGoBridgeRequestBadURL(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()

	bridge.Request(":{}K{RKVR<<>>!@#", nil)
	t.Fatal("no panic")
}

func testGoBridgeRequestWithResponse(t *testing.T, bridge GoBridge) {
	res := bridge.RequestWithResponse("/test", NewPayload(21))

	var nb int
	res.Unmarshal(&nb)

	if nb != 21 {
		t.Fatal("response is not 21:", nb)
	}
}

func testGoBridgeRequestWithResponseBadURL(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()
	bridge.RequestWithResponse(":{}K{RKVR<<>>!@#", nil)
	t.Fatal("no panic")
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
	t.Fatal("no panic")
}
