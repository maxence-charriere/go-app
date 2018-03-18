// +build darwin,amd64

package mac

import (
	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
)

// FilePanel implements the app.Element interface.
type FilePanel struct {
	id uuid.UUID

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

	return driver.elements.Add(panel)
}

// ID satistfies the app.Element interface.
func (p *FilePanel) ID() uuid.UUID {
	return p.id
}

func onFilePanelSelect(p *FilePanel, in map[string]interface{}) interface{} {
	if p.onSelect != nil {
		p.onSelect(bridge.Strings(in["Filenames"]))
	}

	driver.elements.Remove(p)
	return nil
}

func handleFilePanel(h func(p *FilePanel, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		elem, err := driver.elements.Element(id)
		if err != nil {
			return nil
		}

		panel := elem.(*FilePanel)
		return h(panel, in)
	}
}

// SaveFilePanel implements the app.Element interface.
type SaveFilePanel struct {
	id uuid.UUID

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

	return driver.elements.Add(panel)
}

// ID satistfies the app.Element interface.
func (p *SaveFilePanel) ID() uuid.UUID {
	return p.id
}

func onSaveFilePanelSelect(p *SaveFilePanel, in map[string]interface{}) interface{} {
	if p.onSelect != nil {
		p.onSelect(in["Filename"].(string))
	}

	driver.elements.Remove(p)
	return nil
}

func handleSaveFilePanel(h func(p *SaveFilePanel, in map[string]interface{}) interface{}) bridge.GoRPCHandler {
	return func(in map[string]interface{}) interface{} {
		id, _ := uuid.Parse(in["ID"].(string))

		elem, err := driver.elements.Element(id)
		if err != nil {
			return nil
		}

		panel := elem.(*SaveFilePanel)
		return h(panel, in)
	}
}
