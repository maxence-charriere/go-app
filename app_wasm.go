package app

import (
	"context"
	"encoding/json"
	"net/url"
	"syscall/js"
)

var dom = domEngine{
	AttrTransforms: []attrTransform{jsToGoHandler},
	CompoBuilder:   components,
	Sync:           syncDom,
	UI:             UI,
}

func navigate(url string) {
	js.Global().Get("location").Set("href", url)
}

func reload() {
	js.Global().Get("location").Call("reload")
}

func render(c Compo) error {
	// return dom.Render(c)
	return maestre.Render(c)
}

func run() error {
	initEmit()

	url, err := getURL()
	if err != nil {
		return err
	}

	compo, err := maestre.New(compoNameFromURL(url))
	if err != nil {
		return err
	}

	if err := maestre.NewBody(compo); err != nil {
		return err
	}

	// compo, err := components.new(compoNameFromURL(url))
	// if err != nil {
	// 	return err
	// }

	// if err = dom.New(compo); err != nil {
	// 	return err
	// }

	// if nav, ok := compo.(Navigable); ok {
	// 	UI(func() {
	// 		nav.OnNavigate(url)
	// 	})
	// }

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
	// if !components.isRegistered(compoNameFromURL(url)) {
	// 	url.Path = NotFoundPath
	// }
	return url, nil
}

func syncDom(changes []change) error {
	jsChanges := make([]interface{}, len(changes))

	for i, c := range changes {
		jsChange := make(map[string]interface{}, 10)

		setValue := func(k, v string) {
			if v != "" {
				jsChange[k] = v
			}
		}

		jsChange["Action"] = int(c.Action)
		jsChange["NodeID"] = c.NodeID

		setValue("CompoID", c.CompoID)
		setValue("Type", c.Type)
		setValue("Namespace", c.Namespace)
		setValue("Key", c.Key)
		setValue("Value", c.Value)
		setValue("ChildID", c.ChildID)
		setValue("NewChildID", c.NewChildID)

		if c.IsCompo {
			jsChange["IsCompo"] = c.IsCompo
		}

		jsChanges[i] = jsChange
	}

	js.Global().Call("render", jsChanges)
	return nil
}

func initEmit() {
	js.Global().
		Get("goapp").
		Set("emit", js.FuncOf(emit))
}

func emit(this js.Value, args []js.Value) interface{} {
	var m mapping
	if err := json.Unmarshal([]byte(args[0].String()), &m); err != nil {
		Logf("go callback failed: %s", err)
		return nil
	}

	c, err := dom.CompoByID(m.CompoID)
	if err != nil {
		Logf("go callback failed: %s", err)
		return nil
	}

	var f func()
	if f, err = m.Map(c); err != nil {
		Logf("go callback failed: %s", err)
		return nil
	}

	if f != nil {
		f()
		return nil
	}

	Render(c)
	return nil
}
