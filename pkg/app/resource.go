package app

import (
	"net/http"
	"strings"
)

// ResourceProvider is the interface that describes a resource provider that
// tells the Handler how to locate and get the package and static resources.
//
// Package resources are the resource required to operate go-app.
//
// Static resources are resources such as app.wasm, CSS files, images.
//
// The resource provider is used to serve static resources when it satisfies the
// http.Handler interface.
type ResourceProvider interface {
	// Package returns the path where the package resources are located.
	Package() string

	// Static returns the path where the static resources directory (/web) is
	// located.
	Static() string

	// AppWASM returns the app.wasm file path.
	AppWASM() string
}

// LocalDir returns a resource provider that serves static resources from a
// local directory located at the given path.
func LocalDir(root string) ResourceProvider {
	root = strings.Trim(root, "/")
	return localDir{
		Handler: http.FileServer(http.Dir(root)),
		root:    root,
		appWASM: root + "/web/app.wasm",
	}
}

type localDir struct {
	http.Handler
	root    string
	appWASM string
}

func (d localDir) Package() string {
	return d.root
}

func (d localDir) Static() string {
	return d.root
}

func (d localDir) AppWASM() string {
	return d.appWASM
}

// RemoteBucket returns a resource provider that provides resources from a
// remote bucket such as Amazon S3 or Google Cloud Storage.
func RemoteBucket(url string) ResourceProvider {
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, "/web")

	return remoteBucket{
		root:    url,
		appWASM: url + "/web/app.wasm",
	}
}

type remoteBucket struct {
	root    string
	appWASM string
}

func (b remoteBucket) Package() string {
	return ""
}

func (b remoteBucket) Static() string {
	return b.root
}

func (b remoteBucket) AppWASM() string {
	return b.appWASM
}

// GitHubPages returns a resource provider that provides resources from GitHub
// pages. This provider must only be used to generate static websites with the
// GenerateStaticWebsite function.
func GitHubPages(repoName string) ResourceProvider {
	return CustomProvider("", repoName)
}

// CustomProvider returns a resource provider that serves static resources from
// a local directory located at the given path and prefixes URL paths with the
// given prefix.
func CustomProvider(path, prefix string) ResourceProvider {
	root := strings.Trim(path, "/")
	prefix = "/" + strings.Trim(prefix, "/")

	return localDir{
		Handler: http.FileServer(http.Dir(root)),
		root:    prefix,
		appWASM: prefix + "/web/app.wasm",
	}
}

// ProxyResource is a proxy descriptor that maps a given resource to an URL
// path.
type ProxyResource struct {
	// The URL path from where a static resource is accessible.
	Path string

	// The path of the static resource that is proxied. It must start with
	// "/web/".
	ResourcePath string
}
