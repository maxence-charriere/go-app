package app

import (
	"context"
	"net/url"
	"syscall/js"

	"github.com/maxence-charriere/app/internal/maestro"
	"github.com/maxence-charriere/app/pkg/log"
)

var (
	dom *maestro.Dom
)

func init() {
	dom = maestro.NewMaestro(components, UI)
	log.DefaultColor = ""
	log.InfoColor = ""
	log.ErrorColor = ""
	log.WarnColor = ""
	log.DebugColor = ""
}

func navigate(url string) {
	js.Global().Get("location").Set("href", url)
}

// Reload reloads the current page.
func Reload(s, e js.Value) {
	js.Global().Get("location").Call("reload")

}

// NewContextMenu displays a context menu filled with the given menu items.
//
// Context menu requires an app.contextmenu component in the loaded page.
// 	func (c *Compo) Render() string {
// 		return `
// 	<div>
// 		<!-- ... -->
// 		<app.contextmenu>
// 	</div>
// 		`
// 	}
func NewContextMenu(items ...MenuItem) {
}

func render(c Compo) error {
	return dom.Render(c)
}

func run() error {
	url, err := getURL()
	if err != nil {
		return err
	}

	compo, err := components.New(compoNameFromURL(url))
	if err != nil {
		return err
	}

	if err := dom.NewBody(compo); err != nil {
		return err
	}

	if nav, ok := compo.(Navigable); ok {
		UI(func() {
			nav.OnNavigate(url)
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case f := <-ui:
			f()

		case <-ctx.Done():
			return nil
		}
	}
}

func getURL() (*url.URL, error) {
	rawurl := js.Global().
		Get("location").
		Get("href").
		String()

	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if url.Path == "" || url.Path == "/" {
		url.Path = DefaultPath
	}
	if !components.IsImported(compoNameFromURL(url)) {
		url.Path = NotFoundPath
	}
	return url, nil
}
