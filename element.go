package app

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app/markup"
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

	// Contains reports whether component c is mounted in the element.
	Contains(c markup.Component) bool

	// Render renders component c.
	Render(c markup.Component) error

	// LastFocus returns the last time when the element has got focus.
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

	// Size returns the window size.
	Size() (width, height float64)

	// Resize resizes the window to width x height.
	Resize(width, height float64)

	// Focus gives the focus to the window.
	// The window will be put in front, above the other elements.
	Focus()

	// Close closes the element.
	Close()
}

// WindowConfig is a struct that describes a window.
type WindowConfig struct {
	Title           string
	X               float64
	Y               float64
	Width           float64
	MinWidth        float64
	MaxWidth        float64
	Height          float64
	MinHeight       float64
	MaxHeight       float64
	BackgroundColor string
	Borderless      bool
	DisableResize   bool
	DefaultURL      string
	Mac             MacWindowConfig

	OnMinimize       func()
	OnDeminimize     func()
	OnFullScreen     func()
	OnExitFullScreen func()
	OnMove           func(x, y float64)
	OnResize         func(width float64, height float64)
	OnFocus          func()
	OnBlur           func()
	OnClose          func() bool
}

// MacWindowConfig is a struct that describes window fields specific to MacOS.
type MacWindowConfig struct {
	BackgroundVibrancy Vibrancy
	HideCloseButton    bool
	HideMinimizeButton bool
	HideTitleBar       bool
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

	// SetBadge set the dock tile badge with the string representation of v.
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

// ElementStore is the interface that decribes a store of elements.
// It should thread safe.
type ElementStore interface {
	// Add adds an element in the store.
	Add(e Element) error

	// Remove removes an element from the store.
	Remove(e Element)

	// Element returns the element with identifier id.
	Element(id uuid.UUID) (e Element, ok bool)

	// ElementByComponent returns the element where component c is mounted.
	ElementByComponent(c markup.Component) (e ElementWithComponent, ok bool)

	// Sort sorts the elements that hosts components.
	Sort()

	// Len returns the number of elements.
	Len() int
}

// NewElementStore creates an element store.
func NewElementStore() ElementStore {
	return newElementStore(256)
}

type elementStore struct {
	mutex                  sync.Mutex
	capacity               int
	elements               map[uuid.UUID]Element
	elementsWithComponents elementWithComponentList
}

func newElementStore(capacity int) *elementStore {
	return &elementStore{
		capacity:               capacity,
		elements:               make(map[uuid.UUID]Element, capacity),
		elementsWithComponents: make(elementWithComponentList, 0, capacity),
	}
}

func (s *elementStore) Add(e Element) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.elements) == s.capacity {
		return errors.Errorf("can't handle more than %d elements simultaneously", s.capacity)
	}
	s.elements[e.ID()] = e

	if elemWithComp, ok := e.(ElementWithComponent); ok {
		s.elementsWithComponents = append(s.elementsWithComponents, elemWithComp)
		sort.Sort(s.elementsWithComponents)
	}
	return nil
}

func (s *elementStore) Remove(e Element) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.elements, e.ID())

	if _, ok := e.(ElementWithComponent); !ok {
		return
	}

	elements := s.elementsWithComponents
	for i, elem := range elements {
		if elem == e {
			copy(elements[i:], elements[i+1:])
			elements[len(elements)-1] = nil
			elements = elements[:len(elements)-1]
			s.elementsWithComponents = elements
			return
		}
	}
}

func (s *elementStore) Element(id uuid.UUID) (e Element, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e, ok = s.elements[id]
	return
}

func (s *elementStore) ElementByComponent(c markup.Component) (e ElementWithComponent, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, elem := range s.elementsWithComponents {
		if elem.Contains(c) {
			e = elem
			ok = true
			return
		}
	}
	return
}

func (s *elementStore) Sort() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sort.Sort(s.elementsWithComponents)
}

func (s *elementStore) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.elements)
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
