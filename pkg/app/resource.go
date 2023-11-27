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

// ResourceResolver is an interface that defines the method to resolve
// resources from /web/ path to its full URL or file location.
type ResourceResolver interface {
	// Resolve takes a resource path and returns its resolved URL or file path.
	Resolve(string) string
}

// LocalDir returns a ResourceResolver for local resources. It resolves paths
// starting with /web/ to their full file path based on the specified local directory.
// This resolver is suitable for handling resources stored in the local filesystem.
func LocalDir(directory string) ResourceResolver {
	directory = strings.TrimRight(directory, "/")
	return localResourceResolver{
		Handler:   http.FileServer(http.Dir(directory)),
		directory: directory,
	}
}

type localResourceResolver struct {
	http.Handler
	directory string
}

func (r localResourceResolver) Resolve(location string) string {
	if location == "/" || location == "" {
		return "/"
	}
	if remoteLocation(location) || !webLocation(location) {
		return location
	}
	return r.directory + "/" + strings.Trim(location, "/")
}

// RemoteBucket returns a ResourceResolver for remote resources. It resolves
// paths starting with /web/ to their full URL based on the specified remote URL,
// such as a cloud storage bucket. This resolver is ideal for resources hosted
// remotely.
func RemoteBucket(url string) ResourceResolver {
	return remoteResourceResolver{
		url: strings.Trim(url, "/"),
	}
}

type remoteResourceResolver struct {
	url string
}

func (r remoteResourceResolver) Resolve(location string) string {
	if location == "/" || location == "" {
		return "/"
	}
	if remoteLocation(location) || !webLocation(location) {
		return location
	}
	return r.url + "/" + strings.Trim(location, "/")
}

func remoteLocation(location string) bool {
	return strings.HasPrefix(location, "https://") ||
		strings.HasPrefix(location, "http://")
}

func webLocation(location string) bool {
	return strings.HasPrefix(location, "/web/") ||
		location == "/web" ||
		strings.HasPrefix(location, "web/") ||
		location == "web"
}
