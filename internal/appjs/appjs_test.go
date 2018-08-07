package appjs

import (
	"io/ioutil"
	"os/exec"
	"runtime"
	"testing"

	"github.com/murlokswarm/app/internal/html"
)

func TestAppJS(t *testing.T) {
	pagename := "test/appjs.html"
	page := html.NewPage(html.PageConfig{
		Title:       "app.js test",
		Javascripts: []string{"test.js"},
		AppJS:       AppJS("alert"),
		DefaultCompo: `
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
<button onclick="callGoEventHandler('compo-03', 'OnTest', this, event)">Launch</button>

<h2>input</h2>
<input onchange="callGoEventHandler('compo-04', 'OnTest', this, event)" value="Edit me">

<h2>contenteditable</h2>
<div contenteditable onkeyup="callGoEventHandler('compo-05', 'OnTest', this, event)">Edit me</div>

<h2>drag/drop</h2>
<div style="display:inline-block;width:200px;height:200px;background-color:grey;cursor:move;"
	 draggable="true"
	 data-drag="hello world"
	 ondragstart="callGoEventHandler('compo-06', 'OnTest', this, event)">
	Drag me!
</div>
<div style="display:inline-block;width:200px;height:200px;background-color:silver;"
	 ondragover="event.preventDefault()"
	 ondrop="callGoEventHandler('compo-06', 'OnTest', this, event)">
	Drop something here
</div>
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
