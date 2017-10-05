package app

import "testing"

func TestTag(t *testing.T) {
	tag := Tag{
		Type: CompoTag,
	}

	if !tag.Is(CompoTag) {
		t.Error("tag should be a component tag")
	}
	if tag.Is(TextTag) {
		t.Error("tag should not be a text tag")
	}
}
