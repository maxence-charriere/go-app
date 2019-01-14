// +build darwin,amd64

package mac

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// FilePanel implements the app.Elem interface.
type FilePanel struct {
	core.Elem

	id string

	onSelect func(filenames []string)
}

func newFilePanel(c app.FilePanelConfig) *FilePanel {
	p := &FilePanel{
		id:       uuid.New().String(),
		onSelect: c.OnSelect,
	}

	if err := driver.Platform.Call("files.NewPanel", nil, struct {
		ID                string
		MultipleSelection bool
		IgnoreDirectories bool
		IgnoreFiles       bool
		ShowHiddenFiles   bool
		FileTypes         []string `json:",omitempty"`
	}{
		ID:                p.id,
		MultipleSelection: c.MultipleSelection,
		IgnoreDirectories: c.IgnoreDirectories,
		IgnoreFiles:       c.IgnoreFiles,
		ShowHiddenFiles:   c.ShowHiddenFiles,
		FileTypes:         c.FileTypes,
	}); err != nil {
		p.SetErr(err)
		return p
	}

	driver.Elems.Put(p)
	return p
}

// ID satistfies the app.Elem interface.
func (p *FilePanel) ID() string {
	return p.id
}

func onFilePanelSelect(p *FilePanel, in map[string]interface{}) {
	if p.onSelect != nil {
		p.onSelect(core.ConvertToStringSlice(in["Filenames"]))
	}

	driver.Elems.Delete(p)
}

func handleFilePanel(h func(p *FilePanel, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)

		e := driver.Elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return
		}

		p := e.(*FilePanel)
		h(p, in)
	}
}

// SaveFilePanel implements the app.Elem interface.
type SaveFilePanel struct {
	core.Elem

	id string

	onSelect func(filename string)
}

func newSaveFilePanel(c app.SaveFilePanelConfig) *SaveFilePanel {
	p := &SaveFilePanel{
		id: uuid.New().String(),

		onSelect: c.OnSelect,
	}

	if err := driver.Platform.Call("files.NewSavePanel", nil, struct {
		ID              string
		ShowHiddenFiles bool
		FileTypes       []string `json:",omitempty"`
	}{
		ID:              p.id,
		ShowHiddenFiles: c.ShowHiddenFiles,
		FileTypes:       c.FileTypes,
	}); err != nil {
		p.SetErr(err)
		return p
	}

	driver.Elems.Put(p)
	return p
}

// ID satistfies the app.Elem interface.
func (p *SaveFilePanel) ID() string {
	return p.id
}

func onSaveFilePanelSelect(p *SaveFilePanel, in map[string]interface{}) interface{} {
	if p.onSelect != nil {
		p.onSelect(in["Filename"].(string))
	}

	driver.Elems.Delete(p)
	return nil
}

func handleSaveFilePanel(h func(p *SaveFilePanel, in map[string]interface{}) interface{}) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)

		e := driver.Elems.GetByID(id)
		if e.Err() == app.ErrElemNotSet {
			return
		}

		p := e.(*SaveFilePanel)
		h(p, in)
	}
}

// SharePanel implements the app.Elem interface.
type SharePanel struct {
	core.Elem

	id string
}

func newSharePanel(v interface{}) *SharePanel {
	p := &SharePanel{
		id: uuid.New().String(),
	}

	in := struct {
		Share string
		Type  string
	}{
		Share: fmt.Sprint(v),
	}

	switch v.(type) {
	case url.URL, *url.URL:
		in.Type = "url"

	default:
		in.Type = "string"
	}

	err := driver.Platform.Call("driver.Share", nil, in)
	p.SetErr(err)

	return p
}

// ID satisfies the app.Elem interface.
func (p *SharePanel) ID() string {
	return p.id
}
