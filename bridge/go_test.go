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

	bridge := NewGoBridge(func(u *url.URL, p Payload) (res Payload) {
		res = p
		return
	}, uichan)

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "request should success",
			test: func(t *testing.T) { testGoBridgeRequest(t, bridge) },
		},
		{
			name: "request with bad URL should panic",
			test: func(t *testing.T) { testGoBridgeRequestBadURL(t, bridge) },
		},
		{
			name: "Request with response should success",
			test: func(t *testing.T) { testGoBridgeRequestWithResponse(t, bridge) },
		},
		{
			name: "request with response and bad URL should panic",
			test: func(t *testing.T) { testGoBridgeRequestWithResponseBadURL(t, bridge) },
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testGoBridgeRequest(t *testing.T, bridge GoBridge) {
	bridge.Request("/", NewPayload(42))
}

func testGoBridgeRequestBadURL(t *testing.T, bridge GoBridge) {
	defer func() { recover() }()
	bridge.Request(":{}K{RKVR<<>>!@#", nil)
	t.Fatal("should have panic")
}

func testGoBridgeRequestWithResponse(t *testing.T, bridge GoBridge) {
	res := bridge.RequestWithResponse("/", NewPayload(21))

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
