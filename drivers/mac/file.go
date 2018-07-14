// +build darwin,amd64

package mac

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/bridge"
	"github.com/murlokswarm/app/internal/core"
)

// FilePanel implements the app.Element interface.
type FilePanel struct {
	core.Elem
	id string

	onSelect func(filenames []string)
}

func newFilePanel(c app.FilePanelConfig) error {
	panel := &FilePanel{
		id:       uuid.New(),
		onSelect: c.OnSelect,
	}

	if err := driver.macRPC.Call("files.NewPanel", nil, struct {
		ID                string
		MultipleSelection bool
		IgnoreDirectories bool
		IgnoreFiles       bool
		ShowHiddenFiles   bool
		FileTypes         []string `json:",omitempty"`
	}{
		ID:                panel.ID().String(),
		MultipleSelection: c.MultipleSelection,
		IgnoreDirectories: c.IgnoreDirectories,
		IgnoreFiles:       c.IgnoreFiles,
		ShowHiddenFiles:   c.ShowHiddenFiles,
		FileTypes:         c.FileTypes,
	}); err != nil {
		return err
	}

	driver.elems.Put(panel)
	return nil
}

// ID satistfies the app.Element interface.
func (p *FilePanel) ID() string {
	return p.id
}

func onFilePanelSelect(p *FilePanel, in map[string]interface{}) interface{} {
	if p.onSelect != nil {
		p.onSelect(bridge.Strings(in["Filenames"]))
	}

	driver.elems.Delete(p)
	return nil
}

func handleFilePanel(h func(p *FilePanel, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		e := driver.elems.GetByID(id)
		if e.IsNotSet() {
			return nil
		}

		panel := e.(*FilePanel)
		return h(panel, in)
	}
}

// SaveFilePanel implements the app.Element interface.
type SaveFilePanel struct {
	core.Elem
	id string

	onSelect func(filename string)
}

func newSaveFilePanel(c app.SaveFilePanelConfig) error {
	panel := &SaveFilePanel{
		id:       uuid.New(),
		onSelect: c.OnSelect,
	}

	if err := driver.macRPC.Call("files.NewSavePanel", nil, struct {
		ID              string
		ShowHiddenFiles bool
		FileTypes       []string `json:",omitempty"`
	}{
		ID:              panel.ID().String(),
		ShowHiddenFiles: c.ShowHiddenFiles,
		FileTypes:       c.FileTypes,
	}); err != nil {
		return err
	}

	driver.elems.Put(panel)
	return nil
}

// ID satistfies the app.Element interface.
func (p *SaveFilePanel) ID() string {
	return p.id
}

func onSaveFilePanelSelect(p *SaveFilePanel, in map[string]interface{}) interface{} {
	if p.onSelect != nil {
		p.onSelect(in["Filename"].(string))
	}

	driver.elems.Delete(p)
	return nil
}

func handleSaveFilePanel(h func(p *SaveFilePanel, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		e := driver.elems.GetByID(id)
		if e.IsNotSet() {
			return nil
		}

		panel := e.(*SaveFilePanel)
		return h(panel, in)
	}
}
