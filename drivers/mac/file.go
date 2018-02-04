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

func newFilePanel(config app.FilePanelConfig) (panel *FilePanel, err error) {
	panel = &FilePanel{
		id:       uuid.New(),
		onSelect: config.OnSelect,
	}

	if _, err = driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/file/panel/new?id=%v", panel.id),
		bridge.NewPayload(config),
	); err != nil {
		return nil, err
	}

	err = driver.elements.Add(panel)
	return panel, err
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
