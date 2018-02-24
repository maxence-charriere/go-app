package mac

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
)

// FilePanel implements the app.Element interface.
type FilePanel struct {
	id uuid.UUID

	onSelect func(filenames []string)
}

func newFilePanel(config app.FilePanelConfig) error {
	panel := &FilePanel{
		id:       uuid.New(),
		onSelect: config.OnSelect,
	}

	if _, err := driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/file/panel/new?id=%v", panel.id),
		bridge.NewPayload(config),
	); err != nil {
		return err
	}
	return driver.elements.Add(panel)
}

// ID satistfies the app.Element interface.
func (p *FilePanel) ID() uuid.UUID {
	return p.id
}

func onFilePanelClose(panel *FilePanel, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var filenames []string
	p.Unmarshal(&filenames)

	if len(filenames) != 0 && panel.onSelect != nil {
		panel.onSelect(filenames)
	}

	driver.elements.Remove(panel)
	return nil
}

// SaveFilePanel implements the app.Element interface.
type SaveFilePanel struct {
	id uuid.UUID

	onSelect func(filename string)
}

func newSaveFilePanel(config app.SaveFilePanelConfig) error {
	panel := &SaveFilePanel{
		id:       uuid.New(),
		onSelect: config.OnSelect,
	}

	if _, err := driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/file/savepanel/new?id=%v", panel.id),
		bridge.NewPayload(config),
	); err != nil {
		return err
	}
	return driver.elements.Add(panel)
}

// ID satistfies the app.Element interface.
func (p *SaveFilePanel) ID() uuid.UUID {
	return p.id
}

func onSaveFilePanelClose(panel *SaveFilePanel, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var filename string
	p.Unmarshal(&filename)

	if len(filename) != 0 && panel.onSelect != nil {
		panel.onSelect(filename)
	}

	driver.elements.Remove(panel)
	return nil
}
