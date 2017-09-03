package bridge

import "testing"

func TestBridge(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should do a request",
			test: testBridgeDo,
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testBridgeDo(t *testing.T) {
	bridge := &Bridge{
		DoHandler: func(req Request) (res Response, err error) {
			t.Logf("do request %+v", req)
			return
		},
	}

	if _, err := bridge.Do(Request{}); err != nil {
		t.Fatal(err)
	}
}
