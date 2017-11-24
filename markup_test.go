package app_test

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
	"github.com/murlokswarm/app/tests"
)

func TestTag(t *testing.T) {
	tag := app.Tag{
		Type: app.CompoTag,
	}

	if !tag.Is(app.CompoTag) {
		t.Error("tag is not a component tag")
	}
	if tag.Is(app.TextTag) {
		t.Error("tag is not not a text tag")
	}
}

func TestConcurrentMarkup(t *testing.T) {
	tests.TestMarkup(t, func(factory app.Factory) app.Markup {
		return app.NewConcurrentMarkup(html.NewMarkup(factory))
	})
}
