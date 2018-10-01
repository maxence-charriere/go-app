
3.2.9 / 2018-09-30
==================

## General

* `EventSource` now maps source value.

```go
type EventSource struct {
  GoappID string
  CompoID string
  ID      string
  Class   string
  Data    map[string]string
  Value   string
}
```

3.2.7 - 3.2.8 / 2018-09-26
==================

## mac

* fix MacOS retro compatibility up to 10.11 (El Capitan).

3.2.6 / 2018-09-26
==================

## mac

* The way to make a window translucent changed:

```go
win := app.NewWindow(app.WindowConfig{
  Title:             "frosted window",
  Width:             1024,
  Height:            720,
  TitlebarHidden:    true,
  FrostedBackground: true, // Enable translucent effect.
  URL:               "/MyCompo",
})
```

* `WindowConfig.Mac` is deprecated.
* `MacWindowConfig` is deprecated.
* `Vibrancy` is deprecated.
* Vibrancy const are deprecated:
  * `VibeNone`
  * `VibeLight`
  * `VibeDark`
  * `VibeTitlebar`
  * `VibeSelection`
  * `VibeMenu`
  * `VibePopover`
  * `VibeSidebar`
  * `VibeMediumLight`
  * `VibeUltraDark`

3.2.5 / 2018-09-25
==================

## general

* In HTML based elements (windows and pages):
  * `font-family`: Default is the system font.
  * `font-size`: Default is 11px.

## mac

* Fix MacOS Mojave deprecated warnings.
* `goapp mac build` and `goapp mac run`:
  * `-deployment-target`: Specify for which version of MacOS is the build.
  * `-a`: Force rebuild go executable.
  * `-race`: Build with data race detection.

3.2.4 / 2018-09-24
==================

* mac default main window background take the system default color.

3.2.3 / 2018-09-23
==================

* EventArgs structs now have info about the source element that invoked the
  event.

3.2.2 / 2018-09-13
==================

* Fixes `goapp web` on windows.
* Chinese README.

3.2.1 / 2018-09-10
==================

## General

* `goapp update` command has been added.
* Readme has been rewritten.
* Small dom internal refactors.

## mac

* `goapp mac build` has been refactored to be more maintainable and reliable.
* `goapp mac` got various bug fixes.
* `mac.Driver.URL` field has been added. It allows to run and reopen a
  main window without having to do in by hand with `driver.OnRun` and
  `driver.OnReopen`.

## web

* `goapp web build` has been refactored. It now produces a .wapp that contains
  the executable and the resources.
* `goapp web run` command has been added. It can launch the web server and the
  the gopherjs client in the default web browser.

3.2.0 / 2018-08-29
==================

## General

* Logger interface become a function:

  ```go
  var Logger func(format string, a ...interface{}
  ```

* Log functions have been modified:

  ```go
  Log()    // Behave like fmt.Println.
  Logf()   // New - behave like fmt.Printf.
  Panic()  // New - behave like Log followed by panic.
  Panicf() // New - behave like Logf followed by panic.

  Debug() // Removed.
  ```

  The reason of this change is because packaged app like a .app (MacOS) or
  .appx (Windows) does not print their log into the terminal.
  These functions embed the logic to make terminal logs possible.

* `goapp` commands now have verbose mode.

## Mac

* `goapp mac build` now produces a `.app`.
* `goapp mac run` command has been added. It runs a `.app` and capture the
  logs in the terminal.

3.1.2 / 2018-08-27
==================

* goapp is building on Windows and Linux.
* `goapp web init` is now working on Windows.
* `goapp web build` is now working on Windows.

3.1.1 / 2018-08-26
==================

## Mac

* menuitem icon support.
* menuitem checked support.

3.1.0 / 2018-08-25
==================

## General

* New dom engine.

## Web

* Notfound component is now rendered on the client side.

3.0.1 / 2018-08-24
==================

* `goapp mac build -bundle` now produces a valid package when specifying an `AppName`.

3.0.0 / 2018-08-07
==================

## API changes

### General

|v2|v3|
|---|---|
|func NewWindow(c WindowConfig) (Window, error)|func NewWindow(c WindowConfig) Window|
|func NewPage(c PageConfig) error|func NewPage(c PageConfig) Page|
|func NewContextMenu(c MenuConfig) (Menu, error)|func NewContextMenu(c MenuConfig) Menu|
|func NewFilePanel(c FilePanelConfig) error|func NewFilePanel(c FilePanelConfig) Elem|
|func NewSaveFilePanel(c SaveFilePanelConfig) error|func NewSaveFilePanel(c SaveFilePanelConfig) Elem|
|func NewShare(v interface{}) error|func NewShare(v interface{}) Elem|
|func NewNotification(c NotificationConfig) error|func NewNotification(c NotificationConfig) Elem|
|func MenuBar() (Menu, error)|func MenuBar() Menu|
|func NewStatusMenu(c StatusMenuConfig) (StatusMenu, error)|func NewStatusMenu(c StatusMenuConfig) StatusMenu|
|func Dock() (DockTile, error)|func Dock() DockTile|
|func NewAction(name string, arg interface{})|func PostAction(name string, arg interface{})|
|func NewActions(a ...Action)|func PostActions(a ...Action)|
|func Handle(name string, h ActionHandler)|func HandleAction(name string, h ActionHandler)|
|func NewEventSubscriber() EventSubscriber|func NewEventSubscriber() *EventSubscriber|
|func CSSResources() []string|*removed*|
|func CompoNameFromURL(u *url.URL) string|*removed*|
|func CompoNameFromURLString(rawurl string) string|*removed*|
|func NewErrNotFound(object string) error|*removed*|
|func NotFound(err error) bool|*removed*|
|func NewErrNotSupported(feature string) error|*removed*|
|func NotSupported(err error) bool|*removed*|

### Window

|v2|v3|
|---|---|
|Load(url string, v ...interface{}) error|Load(url string, v ...interface{})|
|Render(Compo) error|Render(Compo)|
|LastFocus() time.Time|*removed*|
|Reload() error|Reload()|
|Previous() error|Previous()|
|Next() error|Next()|
|Move(x, y float64) error|Move(x, y float64)|
|Center() error|Center()|
|Resize(width, height float64) error|Resize(width, height float64)|
|Focus() error|Focus()|
|ToggleFullScreen() error|FullScreen()|
||ExitFullScreen()|
|ToggleMinimize() error|Minimize()|
||Deminimize()|
|Close() error|Close()|

### Page

|v2|v3|
|---|---|
|Load(url string, v ...interface{}) error|Load(url string, v ...interface{})|
|Render(Compo) error|Render(Compo)|
|LastFocus() time.Time|*removed*|
|Reload() error|Reload()|
|Previous() error|Previous()|
|Next() error|Next()|
|Close() error|Close()|

### Menu

|v2|v3|
|---|---|
|Load(url string, v ...interface{}) error|Load(url string, v ...interface{})|
|Render(Compo) error|Render(Compo)|
|LastFocus() time.Time|*removed*|

### StatusMenu

|v2|v3|
|---|---|
|Load(url string, v ...interface{}) error|Load(url string, v ...interface{})|
|Render(Compo) error|Render(Compo)|
|LastFocus() time.Time|*removed*|
|SetIcon(name string) error|SetIcon(path string)|
|SetText(text string) error|SetText(text string)|
|Close() error|Close()|

### DockTile

|v2|v3|
|---|---|
|Load(url string, v ...interface{}) error|Load(url string, v ...interface{})|
|Render(Compo) error|Render(Compo)|
|LastFocus() time.Time|*removed*|
|SetIcon(name string) error|SetIcon(path string)|
|SetBadge(v interface{}) error|SetBadge(v interface{})|

## Misc

- A lot of code has been refactored.
- `appjs` and `html` has been moved to `internal`.

2.6.3 / 2018-07-14
==================

* Replace UUID by string (#152)

2.6.2 / 2018-07-14
==================

* refactor history test (#151)
* Move packages to internal (#149)

2.6.1 / 2018-07-13
==================

* interfaces and funcs renamed with shorter names:
  * interfaces:
    * `Component` => `Compo`
    * `ElementWithComponent` => `ElemWithCompo`
    * `ComponentWithExtendedRender` => `CompoWithExtendedRender`
  * functions:
    * `ComponentNameFromURLString` => `CompoNameFromURLString`

2.6.0 / 2018-07-12
==================

* `app.Element` is renamed `app.Elem`

* func to get an element changed.

  Old way (deprecated/deleted):

  ```go
  func ElementByCompo(c Compo) (ElemWithCompo, error) { ... }
  func NavigatorByCompo(c Compo) (Navigator, error)   { ... }
  func WindowByCompo(c Compo) (Window, error)         { ... }
  func PageByCompo(c Compo) (Page, error)             { ... }
  ```

  New way:

  ```go
  func ElemByCompo(c Compo) Elem { ... }
  ```

  Example:

  ```go
  app.ElemByCompo(myCompo).WhenWindow(func(w app.Window) {
    w.Center()
  })
  ```

  See [Elem](https://godoc.org/github.com/murlokswarm/app#Elem) definition for more info.

* Add compatibility with 10.11 (#144)

* MacOS driver is now compatible with MacOS El Capitan (10.11). Thanks to [tectiv3](https://github.com/tectiv3)

2.5.1 / 2018-06-07
==================

* Fix mac dock (#140)

2.5.0 / 2018-06-05
==================

* Status menu (#139)

2.4.3 / 2018-05-27
==================

* fix + travis (#137)

2.4.2 / 2018-05-26
==================

* Refactor logs decorators (#135)

2.4.1 / 2018-05-15
==================

* goapp mac

2.4.0 / 2018-03-25
==================

* Event and subscriber
* Actions

2.3.3 / 2018-03-23
==================

* Refactor decorators func

2.3.2 / 2018-03-18
==================

* bridge go rpc refactor

2.3.1 / 2018-03-16
==================

* bridge platform rpc refactor

2.3.0 / 2018-03-09
==================

* web driver implementation
* goapp tool to build web app

2.2.0 / 2018-02-25
==================

* Drag and drop

2.1.0 / 2018-02-24
==================

* Save file panel
* InnerText in HTML event handler when contenteditable is set

2.0.0 / 2018-02-23
==================

* Use of standard HTML in templates
* Get rid of hidden imports
* Improve code quality
* Improve architecture to make multiplaform easier to implement
* Centralize app related code into a single package (drivers, markup)

1.0.1 / 2018-02-03
==================
  
* Add history.md

1.0.0/ 2018-02-03
==================
  
* V1 release
