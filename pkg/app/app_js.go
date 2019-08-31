package app

import (
	"context"
	"net/url"
	"reflect"
	"syscall/js"

	"github.com/maxence-charriere/app/internal/maestro"
	"github.com/maxence-charriere/app/pkg/log"
)

var (
	dom  *maestro.Dom
	msgs *messenger
)

func init() {
	dom = maestro.NewMaestro(components, UI)
	msgs = &messenger{}

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

// Bind creates a binding between a message and the given component.
func Bind(msg string, c Compo) *Binding {
	b, close := msgs.bind(msg, c)

	if err := dom.SetBindingClose(c, close); err != nil {
		log.Error("creating a binding failed").
			T("reason", err).
			T("component", reflect.TypeOf(c)).
			T("message", msg).
			Panic()
	}

	return b
}

// Emit emits a message that triggers the associated bindings.
func Emit(ctx context.Context, msg string, args ...interface{}) {
	go msgs.emit(ctx, msg, args...)
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
