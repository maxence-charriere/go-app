# Getting started

Some text

## Code

Some code

```go
// Some Comments yeah?
type document struct {
	app.Compo

	path        string
	description string
	document    string
	loading     bool
	err         error
}

func newDocument(path string) *document {
	return &document{path: path}
}

func (d *document) Description(t string) *document {
	d.description = t
	return d
}

func (d *document) OnMount(ctx app.Context) {
	d.loading = true
	d.err = nil
	d.Update()

	go d.load(ctx)
}

func (d *document) load(ctx app.Context) {
	var doc string
	var err error

	defer app.Dispatch(func() {
		if err != nil {
			d.err = err
		}

		d.document = doc
		d.loading = false
		d.Update()
	})

	res, err := http.Get(d.path)
	if err != nil {
		err = errors.New("getting document failed").Wrap(err)
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.New("reading document failed").Wrap(err)
		return
	}

	r := bfchroma.NewRenderer(
		bfchroma.WithoutAutodetect(),
		bfchroma.ChromaOptions(html.WithClasses(true)),
	)

	md := string(bf.Run(b, bf.WithRenderer(r)))
	md = strings.ReplaceAll(md, "\t", "    ")
	doc = fmt.Sprintf("<div>%s</div>", md)
}

func (d *document) Render() app.UI {
	return app.Main().
		Class("layout").
		Class("document").
		Body(
			app.Div().Class("header"),
			app.Div().
				Class("content").
				Body(
					app.Section().Body(
						newLoader().
							Description(d.description).
							Err(d.err).
							Loading(d.loading),
						app.Raw(d.document),
					),
				),
		)
}
```
