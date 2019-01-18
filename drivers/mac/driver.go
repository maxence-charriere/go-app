// Package mac is the driver to be used for apps that run on MacOS.
// It is build on the top of Cocoa and Webkit.
package mac

import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
)

// Driver is the app.Driver implementation for MacOS.
type Driver struct {
	core.Driver `json:"-"`

	// The URL of the component to load in the default window. A non empty value
	// triggers the creation of the default window when the app in openened. It
	// overrides the DefaultWindow.URL value.
	URL string `json:"-"`

	// The default window configuration.
	DefaultWindow app.WindowConfig `json:"-"`

	// Menubar configuration
	MenubarConfig app.MenuBarConfig `json:"-"`

	// The URL of the component to load in the dock.
	DockURL string `json:"-"`

	// The app name. It is displayed in the menubar and dock. The default value
	// is the package directory name.
	//
	// It is used only for goapp packaging.
	Name string `json:",omitempty"`

	// The UTI representing the app.
	//
	// It is used only for goapp packaging.
	ID string `json:",omitempty"`

	// The URL scheme that launches the app.
	//
	// It is used only for goapp packaging.
	URLScheme string `json:",omitempty"`

	// The version of the app (minified form eg. 1.42).
	//
	// It is used only for goapp packaging.
	Version string `json:",omitempty"`

	// The build number.
	//
	// It is used only for goapp packaging.
	BuildNumber int `json:",omitempty"`

	// The app icon path relative to the resources directory. It must be a
	// ".png". Provide a big one! Other required icon sizes will be auto
	// generated.
	//
	// It is used only for goapp packaging.
	Icon string `json:",omitempty"`

	// The development region.
	//
	// It is used only for goapp packaging.
	DevRegion string `json:",omitempty"`

	// A human readable copyright.
	//
	// It is used only for goapp packaging.
	Copyright string `json:",omitempty"`

	// The application category.
	//
	// It is used only for goapp packaging.
	Category Category `json:",omitempty"`

	// Reports wheter the app runs in background mode. Background apps does not
	// appear in the dock and menubar.
	//
	// It is used only for goapp packaging.
	Background bool `json:",omitempty"`

	// Reports whether the app is a server (accepts incoming connections).
	//
	// It is used only for goapp packaging.
	Server bool `json:",omitempty"`

	// Reports whether the app uses the camera.
	//
	// It is used only for goapp packaging.
	Camera bool `json:",omitempty"`

	// Reports whether the app uses the microphone.
	//
	// It is used only for goapp packaging.
	Microphone bool `json:",omitempty"`

	// Reports whether the app uses the USB devices.
	//
	// It is used only for goapp packaging.
	USB bool `json:",omitempty"`

	// Reports whether the app uses printers.
	//
	// It is used only for goapp packaging.
	Printers bool `json:",omitempty"`

	// Reports whether the app uses bluetooth.
	//
	// It is used only for goapp packaging.
	Bluetooth bool `json:",omitempty"`

	// Reports whether the app has access to contacts.
	//
	// It is used only for goapp packaging.
	Contacts bool `json:",omitempty"`

	// Reports whether the app has access to device location.
	//
	// It is used only for goapp packaging.
	Location bool `json:",omitempty"`

	// Reports whether the app has access to calendars.
	//
	// It is used only for goapp packaging.
	Calendar bool `json:",omitempty"`

	// The file picker access mode.
	//
	// It is used only for goapp packaging.
	FilePickers FileAccess `json:",omitempty"`

	// The Download directory access mode.
	//
	// It is used only for goapp packaging.
	Downloads FileAccess `json:",omitempty"`

	// The Pictures directory access mode.
	//
	// It is used only for goapp packaging.
	Pictures FileAccess `json:",omitempty"`

	// The Music directory access mode.
	//
	// It is used only for goapp packaging.
	Music FileAccess `json:",omitempty"`

	// The Movies directory access mode.
	//
	// It is used only for goapp packaging.
	Movies FileAccess `json:",omitempty"`

	// The file types that can be opened by the app.
	//
	// It is used only for goapp packaging.
	SupportedFiles []FileType `json:",omitempty"`

	devID        string
	menubar      *core.Menu
	droppedFiles []string
	stop         func()
}

// Target satisfies the app.Driver interface.
func (d *Driver) Target() string {
	return "macos"
}

// MenuBarConfig contains the menu bar configuration.
type MenuBarConfig struct {
	// The URL of the component to load in the menu bar.
	// Set this to customize the whole menu bar.
	//
	// Default is mac.menubar.
	URL string

	// The URL of the app menu.
	// Set this to customize only the app menu.
	//
	// Default is mac.appmenu.
	AppURL string

	// The URL of the edit menu.
	// Set this to customize only the edit menu.
	//
	// Default is mac.editmenu.
	EditURL string

	// The URL of the window menu.
	// Set this to customize only the window menu.
	//
	// Default is mac.windowmenu.
	WindowURL string

	// An array that contains additional menu URLs.
	CustomURLs []string

	// The URL of the help menu.
	// Set this to customize only the help menu.
	//
	// Default is mac.helpmenu.
	HelpURL string
}

// Role represents the role of an application.
type Role string

// Constants that enumerate application roles.
const (
	Editor Role = "Editor"
	Viewer Role = "Viewer"
	Shell  Role = "Shell"
)

// Category represents the app style.
// The App Store uses this string to determine the appropriate categorization
// for the app.
type Category string

// Constants that enumerate application categories.
const (
	BusinessApp             Category = "public.app-category.business"
	DeveloperToolsApp                = "public.app-category.developer-tools"
	EducationApp                     = "public.app-category.education"
	EntertainmentApp                 = "public.app-category.entertainment"
	FinanceApp                       = "public.app-category.finance"
	GamesApp                         = "public.app-category.games"
	GraphicsAndDesignApp             = "public.app-category.graphics-design"
	HealthcareAndFitnessApp          = "public.app-category.healthcare-fitness"
	LifestyleApp                     = "public.app-category.lifestyle"
	MedicalApp                       = "public.app-category.medical"
	MusicApp                         = "public.app-category.music"
	NewsApp                          = "public.app-category.news"
	PhotographyApp                   = "public.app-category.photography"
	ProductivityApp                  = "public.app-category.productivity"
	ReferenceApp                     = "public.app-category.reference"
	SocialNetworkingApp              = "public.app-category.social-networking"
	SportsApp                        = "public.app-category.sports"
	TravelApp                        = "public.app-category.travel"
	UtilitiesApp                     = "public.app-category.utilities"
	VideoApp                         = "public.app-category.video"
	WeatherApp                       = "public.app-category.weather"
	ActionGamesApp                   = "public.app-category.action-games"
	AdventureGamesApp                = "public.app-category.adventure-games"
	ArcadeGamesApp                   = "public.app-category.arcade-games"
	BoardGamesApp                    = "public.app-category.board-games"
	CardGamesApp                     = "public.app-category.card-games"
	CasinoGamesApp                   = "public.app-category.casino-games"
	DiceGamesApp                     = "public.app-category.dice-games"
	EducationalGamesApp              = "public.app-category.educational-games"
	FamilyGamesApp                   = "public.app-category.family-games"
	KidsGamesApp                     = "public.app-category.kids-games"
	MusicGamesApp                    = "public.app-category.music-games"
	PuzzleGamesApp                   = "public.app-category.puzzle-games"
	RacingGamesApp                   = "public.app-category.racing-games"
	RolePlayingGamesApp              = "public.app-category.role-playing-games"
	SimulationGamesApp               = "public.app-category.simulation-games"
	SportsGamesApp                   = "public.app-category.sports-games"
	StrategyGamesApp                 = "public.app-category.strategy-games"
	TriviaGamesApp                   = "public.app-category.trivia-games"
	WordGamesApp                     = "public.app-category.word-games"
)

// FileAccess represents a file access mode.
type FileAccess string

// Constants that enumerate file access modes.
const (
	NoAccess  FileAccess = ""
	ReadOnly  FileAccess = "read-only"
	ReadWrite FileAccess = "read-write"
)

// FileType describes a file type that can be opened by the app.
type FileType struct {
	// The  name.
	// Must be non empty, a single word and lowercased.
	Name string

	// The app’s role with respect to the type.
	Role Role

	// The icon path:
	// - Must be relative to the resources directory.
	// - Must be a ".png".
	Icon string

	// A list of UTI defining a supported file type.
	// Eg. "public.png" for ".png" files.
	UTIs []string
}
