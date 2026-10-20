package app

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowserAnchorClick(t *testing.T) {
	tests := []struct {
		name                   string
		preventDefault         bool
		canceledBeforeDispatch bool
		redirect               string
		nestedTarget           bool
		wantNavigations        []string
	}{
		{
			name:            "ordinary click navigates",
			wantNavigations: []string{"original"},
		},
		{
			name:           "Go handler cancels navigation",
			preventDefault: true,
		},
		{
			name:            "Go handler redirects without following href",
			preventDefault:  true,
			redirect:        "#custom",
			wantNavigations: []string{"custom"},
		},
		{
			name:                   "already canceled event skips navigation",
			canceledBeforeDispatch: true,
		},
		{
			name:           "Go handler cancels click on nested element",
			preventDefault: true,
			nestedTarget:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine := &engineX{
				ctx:        context.Background(),
				dispatches: make(chan func(), 16),
			}

			var navigations []string
			ctx := Context{
				dispatch: engine.dispatch,
				navigate: func(destination *url.URL, updateHistory bool) {
					require.True(t, updateHistory)
					navigations = append(navigations, destination.Fragment)
				},
			}

			document := Window().Get("document")
			anchor := document.Call("createElement", "a")
			anchor.Set("href", "#original")
			source := A().setJSElement(anchor)
			ctx.sourceElement = source
			handler := nodeManager{}.mountHTMLEventHandler(ctx, source,
				makeEventHandler("click", func(ctx Context, event Event) {
					if test.preventDefault {
						event.PreventDefault()
					}
					if test.redirect != "" {
						ctx.Navigate(test.redirect)
					}
				}))
			t.Cleanup(handler.close)

			originalURL := Window().Get("location").Get("href").String()
			t.Cleanup(func() {
				Window().Get("history").Call("replaceState", nil, "", originalURL)
			})
			document.Get("body").Call("appendChild", anchor)
			t.Cleanup(func() { anchor.Call("remove") })

			previousClick := Window().Get("onclick")
			var browser browser
			browser.handleAnchorClick(ctx)
			t.Cleanup(func() {
				Window().Set("onclick", previousClick)
				browser.anchorClick.Release()
			})

			target := anchor
			if test.nestedTarget {
				target = document.Call("createElement", "span")
				anchor.Call("appendChild", target)
			}
			event := Window().Get("MouseEvent").New("click", map[string]any{
				"bubbles":    true,
				"cancelable": true,
			})
			if test.canceledBeforeDispatch {
				event.Call("preventDefault")
			}
			target.Call("dispatchEvent", event)

			// Process the actual element and window callbacks in queue order.
			// Cancellation must be checked after the Go click handler runs.
			for len(engine.dispatches) != 0 {
				(<-engine.dispatches)()
			}
			require.Equal(t, test.wantNavigations, navigations)
		})
	}
}
