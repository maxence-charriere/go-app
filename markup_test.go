package app_test

import (
	"reflect"
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/html"
	"github.com/murlokswarm/app/internal/tests"
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
	tests.TestMarkup(t, func(factory *app.Factory) app.Markup {
		return app.ConcurrentMarkup(html.NewMarkup(factory))
	})
}

func TestParseMappingTarget(t *testing.T) {
	tests := []struct {
		scenario         string
		target           string
		expectedPipeline []string
		shouldErr        bool
	}{
		{
			scenario:         "parses target",
			target:           "Hello",
			expectedPipeline: []string{"Hello"},
		},
		{
			scenario:         "parses target with multiple elements",
			target:           "Hello.World",
			expectedPipeline: []string{"Hello", "World"},
		},
		{
			scenario:  "parses empyt target returns an error",
			shouldErr: true,
		},
		{
			scenario:  "parses target with empty element returns an error",
			target:    ".Hello.World",
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			pipeline, err := app.ParseMappingTarget(test.target)

			if test.shouldErr {
				if err == nil {
					t.Fatal("error is nil")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(pipeline, test.expectedPipeline) {
				t.Errorf("%v != %v", pipeline, test.expectedPipeline)
			}
		})
	}
}
