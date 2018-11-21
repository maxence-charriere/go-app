// +build darwin,amd64

package mac

// Settings contains settings that define how the app is set up and what are its
// capabilities.
// It is used by goapp to define how the app is built and packaged.
type Settings struct {
	// The app name.
	// It is displayed in the menubar and dock.
	Name string

	// The UTI representing the app.
	ID string

	// The URL scheme that launches the app."
	URLScheme string

	// The version of the app (minified form eg 1.42).
	Version string

	// The build number.
	BuildNumber int

	// The app icon name:
	// - Must be located at the root of the resources directory.
	// - Must be a .png.
	// Provide a big one! Other required icon sizes will be auto generated.
	Icon string

	// The development region.
	DevRegion string

	// A human readable copyright.
	Copyright string

	// The application role.
	Role Role

	// The application category.
	Category Category

	// Reports wheter the app runs in background mode.
	// Background apps does not appear in the dock and menubar.
	Background bool

	// Reports whether the app is a server (accepts incoming connections).
	Server bool

	// Reports whether the app uses the camera.
	Camera bool

	// Reports whether the app uses the microphone.
	Microphone bool

	// Reports whether the app uses the USB devices.
	USB bool

	// Reports whether the app uses printers.
	Printers bool

	// Reports whether the app uses bluetooth.
	Bluetooth bool

	// Reports whether the app has access to contacts.
	Contacts bool

	// Reports whether the app has access to device location.
	Location bool

	// Reports whether the app has access to calendars.
	Calendar bool

	// The file picker access mode.
	FilePickers FileAccess

	// The Download directory access mode.
	Downloads FileAccess

	// The Pictures directory access mode.
	Pictures FileAccess

	// The Music directory access mode.
	Music FileAccess

	// The Movies directory access mode.
	Movies FileAccess

	// The UTIs representing the file types that the app can open.
	SupportedFiles []string
}

// Role represents the role of an application.
type Role string

// Constants that enumerate application roles.
const (
	NoRole     Role = "None"
	EditorRole Role = "Editor"
	ViewerRole Role = "Viewer"
	ShellRole  Role = "Shell"
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
type FileAccess struct {
	Read  bool
	Write bool
}

// MarshalJSON satisfies the json.Marshaler interface.
func (f FileAccess) MarshalJSON() ([]byte, error) {
	switch {
	case f.Read && f.Write:
		return []byte(`"read-write"`), nil

	case f.Read && !f.Write:
		return []byte(`"read-only"`), nil

	default:
		return []byte(`""`), nil
	}
}
