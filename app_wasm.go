package app

import (
	"context"
	"fmt"
	"net/url"
	"syscall/js"
)

var dom = domEngine{
	AttrTransforms: []attrTransform{jsToGoHandler},
	CompoBuilder:   components,
	Sync:           syncDom,
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

	go func() {

		for {
			select {
			case f := <-ui:
				f()

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func syncDom(changes []change) error {
	for _, c := range changes {
		fmt.Printf("%+v\n", c)
	}

	return nil
}
