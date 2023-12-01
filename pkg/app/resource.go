package app

import (
	"net/http"
	"strings"
)

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

// PrefixedLocation returns a ResourceResolver that resolves resources with
// a specified prefix. This resolver prepends the given prefix to resource paths,
// which is particularly useful when serving resources from a specific directory
// or URL path.
//
// The prefix is added to the beginning of the resource path, effectively
// modifying the path from which the resources are served.
//
// For example, if the prefix is "/assets", a resource path like "/web/main.css"
// will be resolved as "/assets/web/main.css".
func PrefixedLocation(prefix string) ResourceResolver {
	return prefixedResourceResolver{
		prefix: prefix,
	}
}

type prefixedResourceResolver struct {
	localResourceResolver
	prefix string
}

func (r prefixedResourceResolver) Resolve(location string) string {
	if remoteLocation(location) {
		return location
	}
	location = r.localResourceResolver.Resolve(location)
	if location == "/" {
		return strings.TrimRight(r.prefix, "/")
	}
	return strings.TrimRight(r.prefix, "/") + location
}

// GitHubPages returns a ResourceResolver tailored for serving resources
// from a GitHub Pages site. It creates a resolver with a prefix matching
// the given repository name. This is particularly useful when hosting
// a go-app project on GitHub Pages, where resources need to be served
// from a repository-specific subpath.
//
// For example, if the repository name is "myapp", the resources will be
// served from paths starting with "/myapp/web/".
func GitHubPages(repositoryName string) ResourceResolver {
	return PrefixedLocation("/" + strings.Trim(repositoryName, "/"))
}

func clientResourceResolver(resourcesLocation string) func(string) string {
	return func(location string) string {
		if remoteLocation(location) || !webLocation(location) {
			return location
		}
		location = strings.Trim(location, "/")
		return resourcesLocation + "/" + strings.TrimPrefix(location, "web/")
	}
}

func resolveOGResource(domain string, location string) string {
	if remoteLocation(location) {
		return location
	}
	return "https://" + domain + strings.TrimRight("/"+strings.Trim(location, "/"), "/")
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

// ProxyResource is a proxy descriptor that maps a given resource to an URL
// path.
type ProxyResource struct {
	// The URL path from where a static resource is accessible.
	Path string

	// The path of the static resource that is proxied. It must start with
	// "/web/".
	ResourcePath string
}
