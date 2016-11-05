package app

import "testing"

func TestResourceLocationJoin(t *testing.T) {
	l := ResourceLocation("resources")

	if j := l.Join("css"); j != "resources/css" {
		t.Error("j should be resources/css:", j)
	}
}

func TestResources(t *testing.T) {
	t.Log(Resources())
}
