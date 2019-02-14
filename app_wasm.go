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

func render(c Compo) error {
	return dom.Render(c)
}

func run() error {
	initEmit()

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
			return nil
		}
	}
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
		Set("emit", js.NewCallback(emit))
}

func emit(args []js.Value) {
	var m mapping
	if err := json.Unmarshal([]byte(args[0].String()), &m); err != nil {
		Logf("go callback failed: %s", err)
		return
	}

	c, err := dom.CompoByID(m.CompoID)
	if err != nil {
		Logf("go callback failed: %s", err)
		return
	}

	var f func()
	if f, err = m.Map(c); err != nil {
		Logf("go callback failed: %s", err)
		return
	}

	if f != nil {
		f()
		return
	}

	Render(c)
}
