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
	page := html.NewPage(html.PageConfig{
		Title:       "app.js test",
		Javascripts: []string{"test.js"},
		AppJS:       AppJS("alert"),
		DefaultComponent: `
<h1>Starting test</h1>
		
<h2>render</h2>
<button onclick="testRender()">Launch</button>
<h3>Output:</h3>
<p data-goapp-id="test-01"></p>
		
<h2>renderAttibutes</h2>
<button onclick="testRenderAttributes()">Launch</button>
<h3>Output:</h3>
<p data-goapp-id="test-02" data-remove="true" data-update="">
	<ul data-goapp-id="test-02-bis">
		<li>data-remove: true</li>
		<li>data-update:</li>
	</ul>
</p>
		
<h2>callback onclick</h2>
<button onclick="callGoEventHandler('compo-03', 'test', this, event)">Launch</button>
		`,
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
