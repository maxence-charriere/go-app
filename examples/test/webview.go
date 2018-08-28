package test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
)

// Webview is a component to test webview based elements.
// It implements the app.Compo interface.
type Webview struct {
	Title       string
	Page        int
	SquareColor string
	Number      int
	CanPrevious bool
	CanNext     bool
}

// Render statisfies the app.Compo interface.
func (c *Webview) Render() string {
	return `
<div class="root" oncontextmenu="OnContextMenu">
	<h1>Test Window</h1>
	<p>
		Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod 
		tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
		quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
		consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
		cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
		non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
	</p>
	
	<ul>
		<li><a href="test.Webview?page=42">To page 42</a></li>
		<li><a href="unknown?page=42">Unknown compopent</a></li>
		<li><a href="http://theverge.com">external hyperlink</a></li>
		<li><button onclick="OnNextPage">Next Page</button></li>
		<li><button onclick="OnLink">External link</button></li>
		<li><button onclick="NotMapped">Not mapped</button></li>
		<li>
			<button onclick="OnChangeSquareColor">Render Attributes: change square color</button>
			<div class="square {{.SquareColor}}"></div>
		</li>

		<li>
			<button onclick="OnChangeNumber">Render: change number</button>
			<div>{{.Number}}</div>
		</li>

		<li>
			<button onclick="OnShare">Share</button>
			<button onclick="OnShareURL">Share URL</button>
		</li>

		<li>
			<button onclick="OnDirPanel">Dir panel</button>
			<button onclick="OnFilePanel">File panel</button>
			<button onclick="OnMultipleFilePanel">Multiple file panel</button>
			<button onclick="OnSavePanel">Save panel</button>
		</li>

		<li>
			<button onclick="OnNotification">Notification</button>
			<button onclick="OnNotificationWithReply">Notification with reply</button>
		</li>

		<li>
			<button {{if not .CanPrevious}}disabled{{end}} onclick="OnPrevious">Previous</button>
			<button onclick="OnReload">Reload</button>
			<button {{if not .CanNext}}disabled{{end}} onclick="OnNext">Next</button>
		</li>

		<li>
			<div class="square dragdrop blue" 
				 draggable="true"
				 data-drag="the blue square on the left"
				 ondragstart="OnDragStart">
				Drag me
			</div>
			<div class="square dragdrop" 
				 ondrop="OnDrop" 
				 ondragover="js:event.preventDefault()">
				Drop something here
			</div>
		</li>
	</ul>
	
	<p>Page: {{.Page}}</p>
</div>
	`
}

// OnContextMenu is the function that is called when the context menu is
// requested.
func (c *Webview) OnContextMenu() {
	app.NewContextMenu(app.MenuConfig{
		URL: "/test.Menu",
		OnClose: func() {
			app.Logf("context menu is closed")
		},
	})
}

// OnNavigate is the function that is called when a component is navigated.
func (c *Webview) OnNavigate(u *url.URL) {
	if pagevals := u.Query()["page"]; len(pagevals) != 0 {
		c.Page, _ = strconv.Atoi(pagevals[0])
	}

	if c.Page == 0 {
		c.Page = 1
	}

	app.ElemByCompo(c).WhenNavigator(func(n app.Navigator) {
		c.CanPrevious = n.CanPrevious()
		c.CanNext = n.CanNext()
	})

	app.Render(c)
}

// PageConfig return allow to set page information like title or meta when the
// component is mounted as the root component.
func (c *Webview) PageConfig() app.HTMLConfig {
	return app.HTMLConfig{
		Title: fmt.Sprintf("Test component %v", c.Page),
	}
}

// OnNextPage is the function to be called when the Next page button is clicked.
func (c *Webview) OnNextPage() {
	page := c.Page
	page++

	app.ElemByCompo(c).WhenNavigator(func(n app.Navigator) {
		n.Load("/test.Webview?page=%v", page)
	})
}

// OnLink is the function to be called when the External link button is clicked.
func (c *Webview) OnLink() {
	app.ElemByCompo(c).WhenNavigator(func(n app.Navigator) {
		n.Load("http://www.github.com")
	})
}

// OnChangeSquareColor is the function to be called when the change color button
// is clicked.
func (c *Webview) OnChangeSquareColor() {
	switch c.SquareColor {
	case "blue":
		c.SquareColor = "pink"
	case "pink":
		c.SquareColor = ""
	default:
		c.SquareColor = "blue"
	}
	app.Render(c)
}

// OnChangeNumber is the function that is called when the change number button is
// clicked.
func (c *Webview) OnChangeNumber() {
	c.Number = rand.Int()
	app.Render(c)
}

// OnShare is the function that is called when the Share button is clicked.
func (c *Webview) OnShare() {
	app.NewShare("Hello world")
}

// OnShareURL is the function that is called when the Share URL button is
// clicked.
func (c *Webview) OnShareURL() {
	u, err := url.Parse("https://github.com/murlokswarm/app")
	if err != nil {
		app.Logf("%s", err)
	}
	app.NewShare(u)
}

// OnDirPanel is the function that is called when the Dir panel button is
// clicked.
func (c *Webview) OnDirPanel() {
	app.NewFilePanel(app.FilePanelConfig{
		IgnoreFiles: true,
		OnSelect: func(filenames []string) {
			app.Logf("filenames: %v", filenames)
		},
	})
}

// OnFilePanel is the function that is called when the File panel button is
// clicked.
func (c *Webview) OnFilePanel() {
	app.NewFilePanel(app.FilePanelConfig{
		IgnoreDirectories: true,
		ShowHiddenFiles:   true,
		FileTypes:         []string{"public.jpeg", "gif"},
		OnSelect: func(filenames []string) {
			app.Logf("filenames: %v", filenames)
		},
	})
}

// OnMultipleFilePanel is the function that is called when the Multiple file
// panel button is clicked.
func (c *Webview) OnMultipleFilePanel() {
	app.NewFilePanel(app.FilePanelConfig{
		IgnoreDirectories: true,
		MultipleSelection: true,
		OnSelect: func(filenames []string) {
			app.Logf("filenames: %v", filenames)
		},
	})
}

// OnSavePanel is the function that is called when the Save panel button is
// clicked.
func (c *Webview) OnSavePanel() {
	app.NewSaveFilePanel(app.SaveFilePanelConfig{
		// FileTypes: []string{"public.jpeg"},
		OnSelect: func(filename string) {
			app.Log(filename)
		},
	})
}

// OnNotification is the function that is called when the Notification button is
// clicked.
func (c *Webview) OnNotification() {
	app.NewNotification(app.NotificationConfig{
		Title:     "hello",
		Subtitle:  "world",
		Text:      uuid.New().String(),
		ImageName: filepath.Join(app.Resources(), "logo.png"),
		Sound:     true,
	})
}

// OnNotificationWithReply is the function that is called when the Notification
// with reply button is clicked.
func (c *Webview) OnNotificationWithReply() {
	id := uuid.New().String()

	if err := app.NewNotification(app.NotificationConfig{
		Title:     "hello",
		Subtitle:  "world",
		Text:      id,
		ImageName: filepath.Join(app.Resources(), "logo.png"),
		Sound:     true,
		OnReply: func(reply string) {
			app.NewNotification(app.NotificationConfig{
				Title:    "reply to",
				Subtitle: id,
				Text:     reply,
				Sound:    true,
			})
		},
	}); err != nil {
		app.Logf("%s", err)
	}
}

// OnPrevious is the function that is called when the previous button is
// clicked.
func (c *Webview) OnPrevious() {
	app.ElemByCompo(c).WhenNavigator(func(n app.Navigator) {
		n.Previous()
	})
}

// OnReload is the function that is called when the reload button is clicked.
func (c *Webview) OnReload() {
	app.ElemByCompo(c).WhenNavigator(func(n app.Navigator) {
		n.Reload()
	})
}

// OnNext is the function that is called when the next button is clicked.
func (c *Webview) OnNext() {
	app.ElemByCompo(c).WhenNavigator(func(n app.Navigator) {
		n.Next()
	})
}

// OnDragStart is the function that is called when the node is dragged.
func (c *Webview) OnDragStart(e app.DragAndDropEvent) {
	data, _ := json.MarshalIndent(e, "", "  ")
	app.Logf("drag: %s", data)
}

// OnDrop is the function that is called when something is dropped to the
// drop zone.
func (c *Webview) OnDrop(e app.DragAndDropEvent) {
	data, _ := json.MarshalIndent(e, "", "  ")
	app.Logf("drop: %s", data)
}
