package app

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Element is the interface that describes an app element.
type Element interface {
	// ID returns the element identifier.
	ID() uuid.UUID
}

// ElementWithComponent is the interface that describes an element that hosts
// components.
type ElementWithComponent interface {
	Element

	// Load loads the page specified by the URL.
	// Calls with an URL which contains a component name will load the named
	// component.
	// e.g. hello will load the component named hello.
	// It returns an error if the component is not imported.
	Load(url string) error

	// Contains reports whether the component is mounted in the element.
	Contains(c Component) bool

	// Render renders the component.
	Render(c Component) error

	// LastFocus returns the last time when the element was focused.
	LastFocus() time.Time
}

// ElementWithNavigation is the interface that describes an element that
// supports navigation.
type ElementWithNavigation interface {
	ElementWithComponent

	// CanPrevious reports whether load the previous page is possible.
	CanPrevious() bool

	// Previous loads the previous page.
	// It returns an error if there is no previous page to load.
	Previous() error

	// CanNext indicates if loading next page is possible.
	CanNext() bool

	// Next loads the next page.
	// It returns an error if there is no next page to load.
	Next() error
}

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

// Menu is the interface that describes a menu.
type Menu ElementWithComponent

// MenuConfig is a struct that describes a menu.
type MenuConfig struct {
	DefaultURL string
}

// DockTile is the interface that describes a dock tile.
type DockTile interface {
	ElementWithComponent

	// SetIcon set the dock tile icon with the named file.
	// It returns an error if the file doesn't exist or if it is not a supported
	// image.
	SetIcon(name string) error

	// SetBadge set the dock tile badge with the string representation of the
	// value.
	SetBadge(v interface{})
}

// FilePanelConfig is a struct that describes a file panel.
type FilePanelConfig struct {
	MultipleSelection bool
	IgnoreDirectories bool
	IgnoreFiles       bool
	OnSelect          func(filenames []string)
}

// PopupNotificationConfig is a struct that describes a popup notification.
type PopupNotificationConfig struct {
	Message      string
	ComponentURL string
}

// ElementDB is the interface that describes an element database.
type ElementDB interface {
	// Add adds the element in the database.
	Add(e Element) error

	// Remove removes the element from the database.
	Remove(e Element)

	// Element returns the element with the given identifier.
	Element(id uuid.UUID) (e Element, ok bool)

	// ElementByComponent returns the element where the component is mounted.
	// It returns an error if the component is not mounted in any element.
	ElementByComponent(c Component) (e ElementWithComponent, err error)

	// ElementsWithComponents returns the elements that contains components.
	ElementsWithComponents() []ElementWithComponent

	// Sort sorts the elements that hosts components. Latest focused elements
	// will be at the beginning.
	Sort()

	// Len returns the number of element.
	Len() int
}

// NewElementDB creates an element database with the given capacity.
// It is safe for concurrent access.
func NewElementDB(capacity int) ElementDB {
	db := newElementDB(capacity)
	return newConcurrentElemDB(db)
}

// elementDB is an element database that implements ElementDB.
type elementDB struct {
	capacity               int
	elements               map[uuid.UUID]Element
	elementsWithComponents elementWithComponentList
}

func newElementDB(capacity int) *elementDB {
	return &elementDB{
		capacity:               capacity,
		elements:               make(map[uuid.UUID]Element, capacity),
		elementsWithComponents: make(elementWithComponentList, 0, capacity),
	}
}

func (db *elementDB) Add(e Element) error {
	if len(db.elements) == db.capacity {
		return errors.Errorf("can't handle more than %d elements simultaneously", db.capacity)
	}

	if _, ok := db.elements[e.ID()]; ok {
		return errors.Errorf("element with id %s is already added", e.ID())
	}

	db.elements[e.ID()] = e

	if elemWithComp, ok := e.(ElementWithComponent); ok {
		db.elementsWithComponents = append(db.elementsWithComponents, elemWithComp)
		sort.Sort(db.elementsWithComponents)
	}
	return nil
}

func (db *elementDB) Remove(e Element) {
	delete(db.elements, e.ID())

	if _, ok := e.(ElementWithComponent); ok {
		elements := db.elementsWithComponents
		for i, elem := range elements {
			if elem == e {
				copy(elements[i:], elements[i+1:])
				elements[len(elements)-1] = nil
				elements = elements[:len(elements)-1]
				db.elementsWithComponents = elements
				return
			}
		}
	}
}

func (db *elementDB) Element(id uuid.UUID) (e Element, ok bool) {
	e, ok = db.elements[id]
	return
}

func (db *elementDB) ElementByComponent(c Component) (e ElementWithComponent, err error) {
	for _, elem := range db.elementsWithComponents {
		if elem.Contains(c) {
			e = elem
			return
		}
	}

	err = errors.Errorf("component %+v is not mounted in any elements", c)
	return
}

func (db *elementDB) ElementsWithComponents() []ElementWithComponent {
	elems := make([]ElementWithComponent, len(db.elementsWithComponents))
	copy(elems, db.elementsWithComponents)
	return elems
}

func (db *elementDB) Sort() {
	sort.Sort(db.elementsWithComponents)
}

func (db *elementDB) Len() int {
	return len(db.elements)
}

// concurrentElemDB is a concurrent element database that implements
// ElementDB.
// It is safe for concurrent access.
type concurrentElemDB struct {
	mutex sync.Mutex
	base  ElementDB
}

func newConcurrentElemDB(db ElementDB) *concurrentElemDB {
	return &concurrentElemDB{
		base: db,
	}
}

func (db *concurrentElemDB) Add(e Element) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.base.Add(e)
}

func (db *concurrentElemDB) Remove(e Element) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.base.Remove(e)
}

func (db *concurrentElemDB) Element(id uuid.UUID) (e Element, ok bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.base.Element(id)
}

func (db *concurrentElemDB) ElementByComponent(c Component) (e ElementWithComponent, err error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.base.ElementByComponent(c)
}

func (db *concurrentElemDB) ElementsWithComponents() []ElementWithComponent {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.base.ElementsWithComponents()
}

func (db *concurrentElemDB) Sort() {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.base.Sort()
}

func (db *concurrentElemDB) Len() int {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.base.Len()
}

// Slice of []ElementWithComponent that implements sort.Interface.
type elementWithComponentList []ElementWithComponent

func (l elementWithComponentList) Len() int {
	return len(l)
}

func (l elementWithComponentList) Less(i, j int) bool {
	return l[i].LastFocus().After(l[j].LastFocus())
}

func (l elementWithComponentList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
