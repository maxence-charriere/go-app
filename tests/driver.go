package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/murlokswarm/app"
)

// TestDriver is a test suite that ensure that all driver implementations behave
// the same.
func TestDriver(t *testing.T, setup func(onRun func()) app.Driver, shutdown func() error) {
	var driver app.Driver

	app.Import(&Hello{})
	app.Import(&World{})
	app.Import(&Menu{})
	app.Import(&Menubar{})

	onRun := func() {
		defer shutdown()

		t.Log("testing driver", driver.Name())
		t.Run("window", func(t *testing.T) { testWindow(t, driver) })
		t.Run("page", func(t *testing.T) { testPage(t, driver) })
		t.Run("context menu", func(t *testing.T) { testContextMenu(t, driver) })
		t.Run("menubar", func(t *testing.T) { testMenubar(t, driver) })
		t.Run("status bar", func(t *testing.T) { testStatusBar(t, driver) })
		t.Run("dock", func(t *testing.T) { testDockTile(t, driver) })

		if err := driver.NewFilePanel(app.FilePanelConfig{}); !app.NotSupported(err) {
			assert.NoError(t, err)
		}

		if err := driver.NewSaveFilePanel(app.SaveFilePanelConfig{}); !app.NotSupported(err) {
			assert.NoError(t, err)
		}

		if err := driver.NewShare(42); !app.NotSupported(err) {
			assert.NoError(t, err)
		}

		if err := driver.NewNotification(app.NotificationConfig{
			Title: "test",
			Text:  "test",
		}); !app.NotSupported(err) {
			assert.NoError(t, err)
		}
	}

	driver = setup(onRun)

	err := app.Run(driver)
	if app.NotSupported(err) {
		return
	}
	require.NoError(t, err)
}
