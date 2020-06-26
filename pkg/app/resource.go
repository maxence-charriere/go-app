package app

import (
	"net/http"
	"strings"
)

var (
	staticResourcesURL string
)

// StaticResource makes a static resource path point to the right
// location whether the root directory is remote or not.
//
// Static resources are resources located in the web directory.
//
// This call is used internally to resolve paths within Cite, Data, Href, Src,
// and SrcSet. Paths already resolved are skipped.
func StaticResource(path string) string {
	if !strings.HasPrefix(path, "/web/") &&
		!strings.HasPrefix(path, "web/") {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return staticResourcesURL + path
}

// Window returns the JavaScript "window" object.
func Window() BrowserWindow {
	return window
}

// ResourceProvider is the interface that describes a provider for static
// resources.
//
// In order to avoid conflict with files required to run a wasm application,
// whether they are located on a local machine or a remote bucket, static
// resources URL paths are always prefixed by "/web".
//
// If the resource provider is an http.handler, the handler is used to serve
// requests with a path that starts by "/web/".
type ResourceProvider interface {
	// The URL of the remote location where static resources are stored.
	//
	// "/web" prefix must not be included in this URL.
	//
	// It must be empty if the resource provider server files from a local
	// location.
	URL() string

	// The URL of the app.wasm file. This must match the pattern:
	//  URL/web/WASM_FILE.
	AppWASM() string

	// The URL of the robots.txt file. This must match the pattern:
	//  URL/web/robots.txt.
	RobotsTxt() string
}

// LocalDir returns a resource provider that serves static resources from a
// local directory located at the given path.
func LocalDir(path string) ResourceProvider {
	return localDir{
		Handler: http.StripPrefix("/web/", http.FileServer(http.Dir(path))),
		path:    path,
	}
}

type localDir struct {
	http.Handler
	path string
}

func (d localDir) URL() string {
	return ""
}

func (d localDir) AppWASM() string {
	return "/web/app.wasm"
}

func (d localDir) RobotsTxt() string {
	return "/web/robots.txt"
}

// RemoteBucket returns a resource provider that provides resources from a
// remote bucket such as Amazon S3 or Google Cloud Storage.
func RemoteBucket(url string) ResourceProvider {
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, "/web")

	return remoteBucket{
		url: url,
	}
}

type remoteBucket struct {
	url string
}

func (b remoteBucket) URL() string {
	return b.url
}

func (b remoteBucket) AppWASM() string {
	return b.URL() + "/web/app.wasm"
}

func (b remoteBucket) RobotsTxt() string {
	return b.URL() + "/web/robots.txt"
}
