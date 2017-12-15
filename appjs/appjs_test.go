package appjs

import (
	"io/ioutil"
	"os/exec"
	"runtime"
	"testing"

	"github.com/murlokswarm/app/html"
)

func TestAppJS(t *testing.T) {
	pagename := "test/appjs.html"
	page := html.Page(html.PageConfig{
		Title:       "app.js test",
		Javascripts: []string{"test.js"},
		AppJS:       AppJS("alert"),
	})

	err := ioutil.WriteFile(pagename, []byte(page), 0644)
	if err != nil {
		t.Fatal(err)
	}

	switch runtime.GOOS {
	case "darwin":
		openOnMacOS(t, pagename)

	case "windows":
		openOnWindows(t, pagename)

	default:
		defaultOpen(t, pagename)
	}
}

func openOnMacOS(t *testing.T, url string) {
	if err := exec.Command("open", "-a", "Google Chrome", url).Start(); err == nil {
		return
	}
	exec.Command("open", url).Start()
}

func openOnWindows(t *testing.T, url string) {
	exec.Command("cmd", "/c", "start", url).Start()
}

func defaultOpen(t *testing.T, url string) {
	exec.Command("xdg-open", url).Start()
}
