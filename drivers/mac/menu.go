// +build darwin,amd64

package mac

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/murlokswarm/app/html"
)

// Menu implements the app.Menu interface.
type Menu struct {
	id        uuid.UUID
	markup    app.Markup
	lastFocus time.Time
	component app.Component

	onClose func()
}

func newMenu(config app.MenuConfig) (m *Menu, err error) {
	m = &Menu{
		id:        uuid.New(),
		markup:    html.NewMarkup(driver.factory),
		lastFocus: time.Now(),

		onClose: config.OnClose,
	}

	if _, err = driver.macos.Request(
		fmt.Sprintf("/menu/new?id=%s", m.id),
		nil,
	); err != nil {
		return
	}

	if err = driver.elements.Add(m); err != nil {
		return
	}

	if len(config.DefaultURL) != 0 {
		err = m.Load(config.DefaultURL)
	}
	return
}

// ID satisfies the app.Element interface.
func (m *Menu) ID() uuid.UUID {
	return m.id
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(rawurl string, v ...interface{}) error {
	var compoName string
	var compo app.Component
	var root app.Tag

	rawurl = fmt.Sprintf(rawurl, v...)
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	compoName = app.ComponentNameFromURL(u)
	if compo, err = driver.factory.NewComponent(compoName); err != nil {
		return err
	}

	if m.component != nil {
		m.markup.Dismount(m.component)
	}

	if _, err = m.markup.Mount(compo); err != nil {
		return err
	}
	m.component = compo

	if navigable, ok := compo.(app.Navigable); ok {
		navigable.OnNavigate(u)
	}

	if root, err = m.markup.Root(compo); err != nil {
		return err
	}

	if root, err = m.markup.FullRoot(root); err != nil {
		return err
	}

	_, err = driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/menu/load?id=%s", m.id),
		bridge.NewPayload(root),
	)
	return err
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(compo app.Component) bool {
	return m.markup.Contains(compo)
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(compo app.Component) error {
	syncs, err := m.markup.Update(compo)
	if err != nil {
		return err
	}

	for _, sync := range syncs {
		if sync.Replace {
			err = m.render(sync)
		} else {
			err = m.renderAttributes(compo, sync)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Menu) render(sync app.TagSync) error {
	tag, err := m.markup.FullRoot(sync.Tag)
	if err != nil {
		return err
	}

	_, err = driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/menu/render?id=%s", m.id),
		bridge.NewPayload(tag),
	)
	return err
}

func (m *Menu) renderAttributes(compo app.Component, sync app.TagSync) error {
	root, err := m.markup.Root(compo)
	if err != nil {
		return err
	}

	tag := sync.Tag
	if root.ID != tag.ID {
		// Ensure that objc will not do extra initializations.
		tag.Children = nil
	}

	_, err = driver.macos.RequestWithAsyncResponse(
		fmt.Sprintf("/menu/render/attributes?id=%s", m.id),
		bridge.NewPayload(tag),
	)
	return err
}

// LastFocus satisfies the app.Menu interface.
func (m *Menu) LastFocus() time.Time {
	return m.lastFocus
}

func onMenuClose(m *Menu, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	// When a context menu button is clicked, the onclick event is called
	// after the onclose one.
	// We have to delay the time we remove the element otherwise the onclick
	// event cannot be called.
	go func() {
		time.Sleep(7 * time.Millisecond)

		app.CallOnUIGoroutine(func() {
			if m.onClose != nil {
				m.onClose()
			}

			driver.elements.Remove(m)
		})
	}()
	return
}

func onMenuCallback(m *Menu, u *url.URL, p bridge.Payload) (res bridge.Payload) {
	var mapping app.Mapping
	p.Unmarshal(&mapping)

	function, err := m.markup.Map(mapping)
	if err != nil {
		app.DefaultLogger.Error(err)
		return
	}

	if function != nil {
		function()
		return
	}

	var compo app.Component
	if compo, err = m.markup.Component(mapping.CompoID); err != nil {
		app.DefaultLogger.Error(err)
		return
	}

	if err = m.Render(compo); err != nil {
		app.DefaultLogger.Error(err)
	}
	return
}
