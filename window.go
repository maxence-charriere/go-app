package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Window is the interface that describes a window.
type Window interface {
	ElementWithNavigation

	// Position returns the window position.
	Position() (x, y float64)

	// Move moves the window to the position (x, y).
	Move(x, y float64)

	// Center moves the window to the center of the screen.
	Center()

	// Size returns the window size.
	Size() (width, height float64)

	// Resize resizes the window to width * height.
	Resize(width, height float64)

	// Focus gives the focus to the window.
	// The window will be put in front, above the other elements.
	Focus()

	// ToggleFullScreen takes the window into or out of fullscreen mode.
	ToggleFullScreen()

	// Minimize takes the window into or out of minimized mode
	ToggleMinimize()

	// Close closes the element.
	Close()
}

// WindowConfig is a struct that describes a window.
type WindowConfig struct {
	Title           string          `json:"title"`
	X               float64         `json:"x"`
	Y               float64         `json:"y"`
	Width           float64         `json:"width"`
	MinWidth        float64         `json:"min-width"`
	MaxWidth        float64         `json:"max-width"`
	Height          float64         `json:"height"`
	MinHeight       float64         `json:"min-height"`
	MaxHeight       float64         `json:"max-height"`
	BackgroundColor string          `json:"background-color"`
	NoResizable     bool            `json:"no-resizable"`
	NoClosable      bool            `json:"no-closable"`
	NoMinimizable   bool            `json:"no-minimizable"`
	TitlebarHidden  bool            `json:"titlebar-hidden"`
	DefaultURL      string          `json:"default-url"`
	Mac             MacWindowConfig `json:"mac"`

	OnMove           func(x, y float64)                  `json:"-"`
	OnResize         func(width float64, height float64) `json:"-"`
	OnFocus          func()                              `json:"-"`
	OnBlur           func()                              `json:"-"`
	OnFullScreen     func()                              `json:"-"`
	OnExitFullScreen func()                              `json:"-"`
	OnMinimize       func()                              `json:"-"`
	OnDeminimize     func()                              `json:"-"`
	OnClose          func() bool                         `json:"-"`
}

// MacWindowConfig is a struct that describes window fields specific to MacOS.
type MacWindowConfig struct {
	BackgroundVibrancy Vibrancy `json:"background-vibrancy"`
}

// Vibrancy represents a constant that define Apple's frost glass effects.
type Vibrancy uint8

// Constants to specify vibrancy effects to use in Apple application elements.
const (
	VibeNone Vibrancy = iota
	VibeLight
	VibeDark
	VibeTitlebar
	VibeSelection
	VibeMenu
	VibePopover
	VibeSidebar
	VibeMediumLight
	VibeUltraDark
)

// NewWindowWithLogs returns a decorated version of the given window that logs
// all the operations.
// Use the default logger.
func NewWindowWithLogs(w Window) Window {
	return &windowWithLogs{
		base: w,
	}
}

type windowWithLogs struct {
	base Window
}

func (w *windowWithLogs) ID() uuid.UUID {
	id := w.base.ID()
	Log("window id is", id)
	return id
}

func (w *windowWithLogs) Load(url string, v ...interface{}) error {
	fmtURL := fmt.Sprintf(url, v...)
	Logf("window %s: loading %s", w.base.ID(), fmtURL)

	err := w.base.Load(url, v...)
	if err != nil {
		Errorf("window %s: loading %s failed: %s", w.base.ID(), fmtURL, err)
	}
	return err
}

func (w *windowWithLogs) Component() Component {
	c := w.base.Component()
	Logf("window %s: mounted component is %T", w.base.ID(), c)
	return c
}

func (w *windowWithLogs) Contains(c Component) bool {
	ok := w.base.Contains(c)
	Logf("window %s: contains %T is %v", w.base.ID(), c, ok)
	return ok
}

func (w *windowWithLogs) Render(c Component) error {
	Logf("window %s: rendering component %T", w.base.ID(), c)

	err := w.base.Render(c)
	if err != nil {
		Errorf("window %s: rendering %T failed: %s", w.base.ID(), c, err)
	}
	return err
}

func (w *windowWithLogs) LastFocus() time.Time {
	focused := w.base.LastFocus()
	Logf("window %s: last focus at %v", w.base.ID(), focused)
	return focused
}

func (w *windowWithLogs) Reload() error {
	Logf("window %s: reloading component %T", w.base.ID(), w.base.Component())

	err := w.base.Reload()
	if err != nil {
		Errorf("window %s: reloading component failed: %s", w.base.ID(), err)
	}
	return err
}

func (w *windowWithLogs) CanPrevious() bool {
	ok := w.base.CanPrevious()
	Logf("window %s: can navigate to previous component is %v", w.base.ID(), ok)
	return ok
}

func (w *windowWithLogs) Previous() error {
	Logf("window %s: navigating to previous component", w.base.ID())

	err := w.base.Previous()
	if err != nil {
		Errorf("window %s: navigating to previous component failed: %s",
			w.base.ID(),
			err,
		)
	}
	return err
}

func (w *windowWithLogs) CanNext() bool {
	ok := w.base.CanNext()
	Logf("window %s: can navigate to next component is %v", w.base.ID(), ok)
	return ok
}

func (w *windowWithLogs) Next() error {
	Logf("window %s: navigating to next component", w.base.ID())

	err := w.base.Next()
	if err != nil {
		Errorf("window %s: navigating to next component failed: %s",
			w.base.ID(),
			err,
		)
	}
	return err
}

func (w *windowWithLogs) Position() (x, y float64) {
	x, y = w.base.Position()
	Logf("window %s: position is (%.2f, %.2f)", w.base.ID(), x, y)
	return x, y
}

func (w *windowWithLogs) Move(x, y float64) {
	Logf("window %s: moving to (%.2f, %.2f)", w.base.ID(), x, y)
	w.base.Move(x, y)
}

func (w *windowWithLogs) Center() {
	Logf("window %s: centering", w.base.ID())
	w.base.Center()
}

func (w *windowWithLogs) Size() (width, height float64) {
	width, height = w.base.Size()
	Logf("window %s: size is %.2fx%.2f", w.base.ID(), width, height)
	return width, height
}

func (w *windowWithLogs) Resize(width, height float64) {
	Logf("window %s: resize to %.2fx%.2f", w.base.ID(), width, height)
	w.base.Resize(width, height)
}

func (w *windowWithLogs) Focus() {
	Logf("window %s: focusing", w.base.ID())
	w.base.Focus()
}

func (w *windowWithLogs) ToggleFullScreen() {
	Logf("window %s: toggle full screen", w.base.ID())
	w.base.ToggleFullScreen()
}

func (w *windowWithLogs) ToggleMinimize() {
	Logf("window %s: toggle minimize", w.base.ID())
	w.base.ToggleMinimize()
}

func (w *windowWithLogs) Close() {
	Logf("window %s: closing", w.base.ID())
	w.base.Close()
}

// NewConcurrentWindow returns a decorated version of the given window that is
// safe for concurrent operations.
func NewConcurrentWindow(w Window) Window {
	return &concurrentWindow{
		base: w,
	}
}

type concurrentWindow struct {
	mutex sync.Mutex
	base  Window
}

func (w *concurrentWindow) ID() uuid.UUID {
	w.mutex.Lock()
	id := w.base.ID()
	w.mutex.Unlock()
	return id
}

func (w *concurrentWindow) Load(url string, v ...interface{}) error {
	w.mutex.Lock()
	err := w.base.Load(url, v...)
	w.mutex.Unlock()
	return err
}

func (w *concurrentWindow) Component() Component {
	w.mutex.Lock()
	c := w.base.Component()
	w.mutex.Unlock()
	return c
}

func (w *concurrentWindow) Contains(c Component) bool {
	w.mutex.Lock()
	ok := w.base.Contains(c)
	w.mutex.Unlock()
	return ok
}

func (w *concurrentWindow) Render(c Component) error {
	w.mutex.Lock()
	err := w.base.Render(c)
	w.mutex.Unlock()
	return err
}

func (w *concurrentWindow) LastFocus() time.Time {
	w.mutex.Lock()
	focused := w.base.LastFocus()
	w.mutex.Unlock()
	return focused
}

func (w *concurrentWindow) Reload() error {
	w.mutex.Lock()
	err := w.base.Reload()
	w.mutex.Unlock()
	return err
}

func (w *concurrentWindow) CanPrevious() bool {
	w.mutex.Lock()
	ok := w.base.CanPrevious()
	w.mutex.Unlock()
	return ok
}

func (w *concurrentWindow) Previous() error {
	w.mutex.Lock()
	err := w.base.Previous()
	w.mutex.Unlock()
	return err
}

func (w *concurrentWindow) CanNext() bool {
	w.mutex.Lock()
	ok := w.base.CanNext()
	w.mutex.Unlock()
	return ok
}

func (w *concurrentWindow) Next() error {
	w.mutex.Lock()
	err := w.base.Next()
	w.mutex.Unlock()
	return err
}

func (w *concurrentWindow) Position() (x, y float64) {
	w.mutex.Lock()
	x, y = w.base.Position()
	w.mutex.Unlock()
	return x, y
}

func (w *concurrentWindow) Move(x, y float64) {
	w.mutex.Lock()
	w.base.Move(x, y)
	w.mutex.Unlock()
}

func (w *concurrentWindow) Center() {
	w.mutex.Lock()
	w.base.Center()
	w.mutex.Unlock()
}

func (w *concurrentWindow) Size() (width, height float64) {
	w.mutex.Lock()
	width, height = w.Size()
	w.mutex.Unlock()
	return width, height
}

func (w *concurrentWindow) Resize(width, height float64) {
	w.mutex.Lock()
	w.base.Resize(width, height)
	w.mutex.Unlock()
}

func (w *concurrentWindow) Focus() {
	w.mutex.Lock()
	w.base.Focus()
	w.mutex.Unlock()
}

func (w *concurrentWindow) ToggleFullScreen() {
	w.mutex.Lock()
	w.base.ToggleFullScreen()
	w.mutex.Unlock()
}

func (w *concurrentWindow) ToggleMinimize() {
	w.mutex.Lock()
	w.base.ToggleMinimize()
	w.mutex.Unlock()
}

func (w *concurrentWindow) Close() {
	w.mutex.Lock()
	w.base.Close()
	w.mutex.Unlock()
}
