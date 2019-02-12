package app

import (
	"context"
	"net/url"
	"syscall/js"
)

var dom = domEngine{
	AttrTransforms: []attrTransform{jsToGoHandler},
	CompoBuilder:   components,
	Sync:           sync,
	UI:             UI,
}

func render(c Compo) error {
	return dom.Render(c)
}

func run() error {
	rawurl := js.Global().
		Get("location").
		Get("href").
		String()

	url, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	var compo Compo
	if compo, err = components.new(compoNameFromURLString(rawurl)); err != nil {
		return err
	}

	if err = dom.New(compo); err != nil {
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
			return
		}
	}

	return ctx.Err()
}

func sync(changes []change) error {
	return erros.New("not implemented")
}
