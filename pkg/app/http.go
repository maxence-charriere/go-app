package app

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

const (
	defaultThemeColor = "#2d2c2c"
)

// Handler is an HTTP handler that serves an HTML page that loads a Go wasm app
// and its resources.
type Handler struct {
	// The name of the web application as it is usually displayed to the user.
	Name string

	// The name of the web application displayed to the user when there is not
	// enough space to display Name.
	ShortName string

	// The icon that is used for the PWA, favicon, loading and default not
	// found component.
	Icon Icon

	// A placeholder background color for the application page to display before
	// its stylesheets are loaded.
	//
	// Default: #2d2c2c.
	BackgroundColor string

	// The theme color for the application. This affects how the OS displays the
	// app (e.g., PWA title bar or Android's task switcher).
	//
	// DEFAULT: #2d2c2c.
	ThemeColor string

	// The text displayed while loading a page. Load progress can be inserted by
	// including "{progress}" in the loading label.
	//
	// DEFAULT: "{progress}%".
	LoadingLabel string

	// The page language.
	//
	// DEFAULT: en.
	Lang string

	// The custom libraries to load with the page.
	Libraries []Library

	// The page title.
	Title string

	// The page description.
	Description string

	// The page authors.
	Author string

	// The page keywords.
	Keywords []string

	// The path of the default image that is used by social networks when
	// linking the app.
	Image string

	// The paths or urls of the CSS files to use with the page.
	//
	// eg:
	//  app.Handler{
	//      Styles: []string{
	//          "/web/test.css",            // Static resource
	//          "https://foo.com/test.css", // External resource
	//      },
	//  },
	Styles []string

	// The paths or urls of the font files to preload with the page.
	//
	// eg:
	//  app.Handler{
	//      Fonts: []string{
	//          "/web/test.woff2",            // Static resource
	//          "https://foo.com/test.woff2", // External resource
	//      },
	//  },
	Fonts []string

	// The paths or urls of the JavaScript files to use with the page.
	//
	// eg:
	//  app.Handler{
	//      Scripts: []string{
	//          "/web/test.js",            // Static resource
	//          "https://foo.com/test.js", // External resource
	//      },
	//  },
	Scripts []string

	// The path of the static resources that the browser is caching in order to
	// provide offline mode.
	//
	// Note that Icon, Styles and Scripts are already cached by default.
	//
	// Paths are relative to the root directory.
	CacheableResources []string

	// Additional headers to be added in head element.
	RawHeaders []string

	// The page HTML element.
	//
	// Default: Html().
	HTML func() HTMLHtml

	// The page body element.
	//
	// Note that the lang attribute is always overridden by the Handler.Lang
	// value.
	//
	// Default: Body().
	Body func() HTMLBody

	// The interval between each app auto-update while running in a web browser.
	// Zero or negative values deactivates the auto-update mechanism.
	//
	// Default is 0.
	AutoUpdateInterval time.Duration

	// The environment variables that are passed to the progressive web app.
	//
	// Reserved keys:
	// - GOAPP_VERSION
	// - GOAPP_GOAPP_STATIC_RESOURCES_URL
	Env Environment

	// The URLs that are launched in the app tab or window.
	//
	// By default, URLs with a different domain are launched in another tab.
	// Specifying internal URLs is to override that behavior. A good use case
	// would be the URL for an OAuth authentication.
	InternalURLs []string

	// The URLs of the origins to preconnect in order to improve the user
	// experience by preemptively initiating a connection to those origins.
	// Preconnecting speeds up future loads from a given origin by preemptively
	// performing part or all of the handshake (DNS+TCP for HTTP, and
	// DNS+TCP+TLS for HTTPS origins).
	Preconnect []string

	// The static resources that are accessible from custom paths. Files that
	// are proxied by default are /robots.txt, /sitemap.xml and /ads.txt.
	ProxyResources []ProxyResource

	// The resource provider that provides static resources. Static resources
	// are always accessed from a path that starts with "/web/".
	//
	// eg:
	//  "/web/main.css"
	//
	// Default: LocalDir("")
	Resources ResourceProvider

	// The version number. This is used in order to update the PWA application
	// in the browser. It must be set when deployed on a live system in order to
	// prevent recurring updates.
	//
	// Default: Auto-generated in order to trigger pwa update on a local
	// development system.
	Version string

	// The HTTP header to retrieve the WebAssembly file content length.
	//
	// Content length finding falls back to the Content-Length HTTP header when
	// no content length is found with the defined header.
	WasmContentLengthHeader string

	// The template used to generate app-worker.js. The template follows the
	// text/template package model.
	//
	// By default set to DefaultAppWorkerJS, changing the template have very
	// high chances to mess up go-app usage. Any issue related to a custom app
	// worker template is not supported and will be closed.
	ServiceWorkerTemplate string

	once                 sync.Once
	etag                 string
	libraries            map[string][]byte
	proxyResources       map[string]ProxyResource
	cachedProxyResources *memoryCache
	cachedPWAResources   *memoryCache
}

func (h *Handler) init() {
	h.initVersion()
	h.initStaticResources()
	h.initImage()
	h.initLibraries()
	h.initLinks()
	h.initScripts()
	h.initServiceWorker()
	h.initCacheableResources()
	h.initIcon()
	h.initPWA()
	h.initPageContent()
	h.initPWAResources()
	h.initProxyResources()
}

func (h *Handler) initVersion() {
	if h.Version == "" {
		t := time.Now().UTC().String()
		h.Version = fmt.Sprintf(`%x`, sha1.Sum([]byte(t)))
	}
	h.etag = `"` + h.Version + `"`
}

func (h *Handler) initStaticResources() {
	if h.Resources == nil {
		h.Resources = LocalDir("")
	}
}

func (h *Handler) initImage() {
	if h.Image != "" {
		h.Image = h.resolveStaticPath(h.Image)
	}
}

func (h *Handler) initLibraries() {
	libs := make(map[string][]byte)
	for _, l := range h.Libraries {
		path, styles := l.Styles()
		if !strings.HasPrefix(path, "/") || len(styles) == 0 {
			continue
		}
		libs[path] = []byte(styles)
	}
	h.libraries = libs
}

func (h *Handler) initLinks() {
	for i, path := range h.Preconnect {
		h.Preconnect[i] = h.resolveStaticPath(path)
	}

	styles := []string{h.resolvePackagePath("/app.css")}
	for path := range h.libraries {
		styles = append(styles, h.resolvePackagePath(path))
	}
	for _, path := range h.Styles {
		styles = append(styles, h.resolveStaticPath(path))
	}
	h.Styles = styles

	for i, path := range h.Fonts {
		h.Fonts[i] = h.resolveStaticPath(path)
	}
}

func (h *Handler) initScripts() {
	for i, path := range h.Scripts {
		h.Scripts[i] = h.resolveStaticPath(path)
	}
}

func (h *Handler) initServiceWorker() {
	if h.ServiceWorkerTemplate == "" {
		h.ServiceWorkerTemplate = DefaultAppWorkerJS
	}
}

func (h *Handler) initCacheableResources() {
	for i, path := range h.CacheableResources {
		h.CacheableResources[i] = h.resolveStaticPath(path)
	}
}

func (h *Handler) initIcon() {
	if h.Icon.Default == "" {
		h.Icon.Default = "https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.png"
		h.Icon.Large = "https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.png"
	}

	if h.Icon.AppleTouch == "" {
		h.Icon.AppleTouch = h.Icon.Default
	}

	if h.Icon.SVG == "" {
		h.Icon.SVG = "https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.svg"
	}

	h.Icon.Default = h.resolveStaticPath(h.Icon.Default)
	h.Icon.Large = h.resolveStaticPath(h.Icon.Large)
	h.Icon.SVG = h.resolveStaticPath(h.Icon.SVG)
	h.Icon.AppleTouch = h.resolveStaticPath(h.Icon.AppleTouch)
}

func (h *Handler) initPWA() {
	if h.Name == "" && h.ShortName == "" && h.Title == "" {
		h.Name = "App PWA"
	}
	if h.ShortName == "" {
		h.ShortName = h.Name
	}
	if h.Name == "" {
		h.Name = h.ShortName
	}

	if h.BackgroundColor == "" {
		h.BackgroundColor = defaultThemeColor
	}
	if h.ThemeColor == "" {
		h.ThemeColor = defaultThemeColor
	}

	if h.Lang == "" {
		h.Lang = "en"
	}

	if h.LoadingLabel == "" {
		h.LoadingLabel = "{progress}%"
	}
}

func (h *Handler) initPageContent() {
	if h.HTML == nil {
		h.HTML = Html
	}

	if h.Body == nil {
		h.Body = Body
	}

}

func (h *Handler) initPWAResources() {
	h.cachedPWAResources = newMemoryCache(5)

	h.cachedPWAResources.Set(cacheItem{
		Path:        "/wasm_exec.js",
		ContentType: "application/javascript",
		Body:        []byte(wasmExecJS()),
	})

	h.cachedPWAResources.Set(cacheItem{
		Path:        "/app.js",
		ContentType: "application/javascript",
		Body:        h.makeAppJS(),
	})

	h.cachedPWAResources.Set(cacheItem{
		Path:        "/app-worker.js",
		ContentType: "application/javascript",
		Body:        h.makeAppWorkerJS(),
	})

	h.cachedPWAResources.Set(cacheItem{
		Path:        "/manifest.webmanifest",
		ContentType: "application/manifest+json",
		Body:        h.makeManifestJSON(),
	})

	h.cachedPWAResources.Set(cacheItem{
		Path:        "/app.css",
		ContentType: "text/css",
		Body:        []byte(appCSS),
	})
}

func (h *Handler) makeAppJS() []byte {
	if h.Env == nil {
		h.Env = make(map[string]string)
	}
	internalURLs, _ := json.Marshal(h.InternalURLs)
	h.Env["GOAPP_INTERNAL_URLS"] = string(internalURLs)
	h.Env["GOAPP_VERSION"] = h.Version
	h.Env["GOAPP_STATIC_RESOURCES_URL"] = h.Resources.Static()
	h.Env["GOAPP_ROOT_PREFIX"] = h.Resources.Package()

	for k, v := range h.Env {
		if err := os.Setenv(k, v); err != nil {
			Log(errors.New("setting app env variable failed").
				WithTag("name", k).
				WithTag("value", v).
				Wrap(err))
		}
	}

	var b bytes.Buffer
	if err := template.
		Must(template.New("app.js").Parse(appJS)).
		Execute(&b, struct {
			Env                     string
			LoadingLabel            string
			Wasm                    string
			WasmContentLengthHeader string
			WorkerJS                string
			AutoUpdateInterval      int64
		}{
			Env:                     jsonString(h.Env),
			LoadingLabel:            h.LoadingLabel,
			Wasm:                    h.Resources.AppWASM(),
			WasmContentLengthHeader: h.WasmContentLengthHeader,
			WorkerJS:                h.resolvePackagePath("/app-worker.js"),
			AutoUpdateInterval:      h.AutoUpdateInterval.Milliseconds(),
		}); err != nil {
		panic(errors.New("initializing app.js failed").Wrap(err))
	}
	return b.Bytes()
}

func (h *Handler) makeAppWorkerJS() []byte {
	resources := make(map[string]struct{})
	setResources := func(res ...string) {
		for _, r := range res {
			if r == "" {
				continue
			}

			r, _, _ = parseSrc(r)
			resources[r] = struct{}{}
		}
	}
	setResources(
		h.resolvePackagePath("/app.css"),
		h.resolvePackagePath("/app.js"),
		h.resolvePackagePath("/manifest.webmanifest"),
		h.resolvePackagePath("/wasm_exec.js"),
		h.resolvePackagePath("/"),
		h.Resources.AppWASM(),
	)
	setResources(h.Icon.Default, h.Icon.Large, h.Icon.AppleTouch)
	setResources(h.Styles...)
	setResources(h.Fonts...)
	setResources(h.Scripts...)
	setResources(h.CacheableResources...)

	resourcesTocache := make([]string, 0, len(resources))
	for k := range resources {
		resourcesTocache = append(resourcesTocache, k)
	}
	sort.Slice(resourcesTocache, func(a, b int) bool {
		return strings.Compare(resourcesTocache[a], resourcesTocache[b]) < 0
	})

	var b bytes.Buffer
	if err := template.
		Must(template.New("app-worker.js").Parse(h.ServiceWorkerTemplate)).
		Execute(&b, struct {
			Version          string
			ResourcesToCache string
		}{
			Version:          h.Version,
			ResourcesToCache: jsonString(resourcesTocache),
		}); err != nil {
		panic(errors.New("initializing app-worker.js failed").Wrap(err))
	}
	return b.Bytes()
}

func (h *Handler) makeManifestJSON() []byte {
	normalize := func(s string) string {
		if !strings.HasPrefix(s, "/") {
			s = "/" + s
		}
		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		return s
	}

	var b bytes.Buffer
	if err := template.
		Must(template.New("manifest.webmanifest").Parse(manifestJSON)).
		Execute(&b, struct {
			ShortName       string
			Name            string
			Description     string
			DefaultIcon     string
			LargeIcon       string
			SVGIcon         string
			BackgroundColor string
			ThemeColor      string
			Scope           string
			StartURL        string
		}{
			ShortName:       h.ShortName,
			Name:            h.Name,
			Description:     h.Description,
			DefaultIcon:     h.Icon.Default,
			LargeIcon:       h.Icon.Large,
			SVGIcon:         h.Icon.SVG,
			BackgroundColor: h.BackgroundColor,
			ThemeColor:      h.ThemeColor,
			Scope:           normalize(h.Resources.Package()),
			StartURL:        normalize(h.Resources.Package()),
		}); err != nil {
		panic(errors.New("initializing manifest.webmanifest failed").Wrap(err))
	}
	return b.Bytes()
}

func (h *Handler) initProxyResources() {
	h.cachedProxyResources = newMemoryCache(len(h.ProxyResources))
	resources := make(map[string]ProxyResource)

	for _, r := range h.ProxyResources {
		switch r.Path {
		case "/wasm_exec.js",
			"/goapp.js",
			"/app.js",
			"/app-worker.js",
			"/manifest.json",
			"/manifest.webmanifest",
			"/app.css",
			"/app.wasm",
			"/goapp.wasm",
			"/":
			continue

		default:
			if strings.HasPrefix(r.Path, "/") && strings.HasPrefix(r.ResourcePath, "/web/") {
				resources[r.Path] = r
			}
		}
	}

	if _, ok := resources["/robots.txt"]; !ok {
		resources["/robots.txt"] = ProxyResource{
			Path:         "/robots.txt",
			ResourcePath: "/web/robots.txt",
		}
	}
	if _, ok := resources["/sitemap.xml"]; !ok {
		resources["/sitemap.xml"] = ProxyResource{
			Path:         "/sitemap.xml",
			ResourcePath: "/web/sitemap.xml",
		}
	}
	if _, ok := resources["/ads.txt"]; !ok {
		resources["/ads.txt"] = ProxyResource{
			Path:         "/ads.txt",
			ResourcePath: "/web/ads.txt",
		}
	}

	h.proxyResources = resources
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", h.etag)

	etag := r.Header.Get("If-None-Match")
	if etag == h.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	path := r.URL.Path

	fileHandler, isServingStaticResources := h.Resources.(http.Handler)
	if isServingStaticResources && strings.HasPrefix(path, "/web/") {
		fileHandler.ServeHTTP(w, r)
		return
	}

	switch path {
	case "/goapp.js":
		path = "/app.js"

	case "/manifest.json":
		path = "/manifest.webmanifest"

	case "/app.wasm", "/goapp.wasm":
		if isServingStaticResources {
			r2 := *r
			r2.URL.Path = h.Resources.AppWASM()
			fileHandler.ServeHTTP(w, &r2)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return

	}

	if res, ok := h.cachedPWAResources.Get(path); ok {
		h.serveCachedItem(w, res)
		return
	}

	if proxyResource, ok := h.proxyResources[path]; ok {
		h.serveProxyResource(proxyResource, w, r)
		return
	}

	if library, ok := h.libraries[path]; ok {
		h.serveLibrary(w, r, library)
		return
	}

	h.servePage(w, r)
}

func (h *Handler) serveCachedItem(w http.ResponseWriter, i cacheItem) {
	w.Header().Set("Content-Length", strconv.Itoa(i.Len()))
	w.Header().Set("Content-Type", i.ContentType)

	if i.ContentEncoding != "" {
		w.Header().Set("Content-Encoding", i.ContentEncoding)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(i.Body)
}

func (h *Handler) serveProxyResource(resource ProxyResource, w http.ResponseWriter, r *http.Request) {
	var u string
	if _, ok := h.Resources.(http.Handler); ok {
		var protocol string
		if r.TLS != nil {
			protocol = "https://"
		} else {
			protocol = "http://"
		}
		u = protocol + r.Host + resource.ResourcePath
	} else {
		u = h.Resources.Static() + resource.ResourcePath
	}

	if i, ok := h.cachedProxyResources.Get(resource.Path); ok {
		h.serveCachedItem(w, i)
		return
	}

	res, err := http.Get(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Log(errors.New("getting proxy static resource failed").
			WithTag("url", u).
			WithTag("proxy-path", resource.Path).
			WithTag("static-resource-path", resource.ResourcePath).
			Wrap(err),
		)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Log(errors.New("reading proxy static resource failed").
			WithTag("url", u).
			WithTag("proxy-path", resource.Path).
			WithTag("static-resource-path", resource.ResourcePath).
			Wrap(err),
		)
		return
	}

	item := cacheItem{
		Path:            resource.Path,
		ContentType:     res.Header.Get("Content-Type"),
		ContentEncoding: res.Header.Get("Content-Encoding"),
		Body:            body,
	}
	h.cachedProxyResources.Set(item)
	h.serveCachedItem(w, item)
}

func (h *Handler) servePage(w http.ResponseWriter, r *http.Request) {
	content, ok := routes.createComponent(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	url := *r.URL
	url.Host = r.Host
	url.Scheme = "http"

	page := requestPage{
		url:                   &url,
		resolveStaticResource: h.resolveStaticPath,
	}
	page.SetTitle(h.Title)
	page.SetLang(h.Lang)
	page.SetDescription(h.Description)
	page.SetAuthor(h.Author)
	page.SetKeywords(h.Keywords...)
	page.SetLoadingLabel(strings.ReplaceAll(h.LoadingLabel, "{progress}", "0"))
	page.SetImage(h.Image)

	disp := engine{
		Page:                   &page,
		IsServerSide:           true,
		StaticResourceResolver: h.resolveStaticPath,
		ActionHandlers:         actionHandlers,
	}
	body := h.Body().privateBody(
		Div(), // Pre-rendeging placeholder
		Aside().
			ID("app-wasm-loader").
			Class("goapp-app-info").
			Body(
				Img().
					ID("app-wasm-loader-icon").
					Class("goapp-logo goapp-spin").
					Src(h.Icon.Default),
				P().
					ID("app-wasm-loader-label").
					Class("goapp-label").
					Text(page.loadingLabel),
			),
	)
	if err := mount(&disp, body); err != nil {
		panic(errors.New("mounting pre-rendering container failed").
			WithTag("server-side", disp.isServerSide()).
			WithTag("body-type", reflect.TypeOf(disp.Body)).
			Wrap(err))
	}
	disp.Body = body
	disp.init()
	defer disp.Close()

	disp.Mount(content)

	for len(disp.dispatches) != 0 {
		disp.Consume()
		disp.Wait()
	}

	icon := h.Icon.SVG
	if icon == "" {
		icon = h.Icon.Default
	}

	var b bytes.Buffer
	b.WriteString("<!DOCTYPE html>\n")
	PrintHTML(&b, h.HTML().
		Lang(page.Lang()).
		privateBody(
			Head().Body(
				Meta().Charset("UTF-8"),
				Meta().
					Name("author").
					Content(page.Author()),
				Meta().
					Name("description").
					Content(page.Description()),
				Meta().
					Name("keywords").
					Content(page.Keywords()),
				Meta().
					Name("theme-color").
					Content(h.ThemeColor),
				Meta().
					Name("viewport").
					Content("width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0, viewport-fit=cover"),
				Meta().
					Property("og:url").
					Content(page.URL().String()),
				Meta().
					Property("og:title").
					Content(page.Title()),
				Meta().
					Property("og:description").
					Content(page.Description()),
				Meta().
					Property("og:type").
					Content("website"),
				Meta().
					Property("og:image").
					Content(page.Image()),
				Range(page.twitterCardMap).Map(func(k string) UI {
					v := page.twitterCardMap[k]
					if v == "" {
						return nil
					}
					return Meta().
						Name(k).
						Content(v)
				}),
				Title().Text(page.Title()),
				Range(h.Preconnect).Slice(func(i int) UI {
					url, crossOrigin, _ := parseSrc(h.Preconnect[i])
					if url == "" {
						return nil
					}

					link := Link().
						Rel("preconnect").
						Href(url)

					if crossOrigin != "" {
						link = link.CrossOrigin(strings.Trim(crossOrigin, "true"))
					}

					return link
				}),
				Range(h.Fonts).Slice(func(i int) UI {
					url, crossOrigin, _ := parseSrc(h.Fonts[i])
					if url == "" {
						return nil
					}

					link := Link().
						Type("font/" + strings.TrimPrefix(filepath.Ext(url), ".")).
						Rel("preload").
						Href(url).
						As("font")

					if crossOrigin != "" {
						link = link.CrossOrigin(strings.Trim(crossOrigin, "true"))
					}

					return link
				}),
				Range(page.Preloads()).Slice(func(i int) UI {
					p := page.Preloads()[i]
					if p.Href == "" || p.As == "" {
						return nil
					}

					url, crossOrigin, _ := parseSrc(p.Href)
					if url == "" {
						return nil
					}

					link := Link().
						Type(p.Type).
						Rel("preload").
						Href(url).
						As(p.As).
						FetchPriority(p.FetchPriority)

					if crossOrigin != "" {
						link = link.CrossOrigin(strings.Trim(crossOrigin, "true"))
					}

					return link
				}),
				Range(h.Styles).Slice(func(i int) UI {
					url, crossOrigin, _ := parseSrc(h.Styles[i])
					if url == "" {
						return nil
					}

					link := Link().
						Type("text/css").
						Rel("preload").
						Href(url).
						As("style")

					if crossOrigin != "" {
						link = link.CrossOrigin(strings.Trim(crossOrigin, "true"))
					}

					return link
				}),
				Link().
					Rel("icon").
					Href(icon),
				Link().
					Rel("apple-touch-icon").
					Href(h.Icon.AppleTouch),
				Link().
					Rel("manifest").
					Href(h.resolvePackagePath("/manifest.webmanifest")),
				Range(h.Styles).Slice(func(i int) UI {
					url, crossOrigin, _ := parseSrc(h.Styles[i])
					if url == "" {
						return nil
					}

					link := Link().
						Rel("stylesheet").
						Type("text/css").
						Href(url)

					if crossOrigin != "" {
						link = link.CrossOrigin(strings.Trim(crossOrigin, "true"))
					}

					return link
				}),
				Script().
					Defer(true).
					Src(h.resolvePackagePath("/wasm_exec.js")),
				Script().
					Defer(true).
					Src(h.resolvePackagePath("/app.js")),
				Range(h.Scripts).Slice(func(i int) UI {
					url, crossOrigin, loading := parseSrc(h.Scripts[i])
					if url == "" {
						return nil
					}

					script := Script().Src(url)

					if crossOrigin != "" {
						script = script.CrossOrigin(strings.Trim(crossOrigin, "true"))
					}

					switch loading {
					case "defer":
						script = script.Defer(true)

					case "async":
						script.Async(true)
					}

					return script
				}),
				Range(h.RawHeaders).Slice(func(i int) UI {
					return Raw(h.RawHeaders[i])
				}),
			),
			body,
		))

	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	w.Header().Set("Content-Type", "text/html")
	w.Write(b.Bytes())
}

func (h *Handler) serveLibrary(w http.ResponseWriter, r *http.Request, library []byte) {
	w.Header().Set("Content-Length", strconv.Itoa(len(library)))
	w.Header().Set("Content-Type", "text/css")
	w.Write(library)
}

func (h *Handler) resolvePackagePath(path string) string {
	var b strings.Builder

	b.WriteByte('/')
	appResources := strings.Trim(h.Resources.Package(), "/")
	b.WriteString(appResources)

	path = strings.Trim(path, "/")
	if b.Len() != 1 && path != "" {
		b.WriteByte('/')
	}
	b.WriteString(path)

	return b.String()
}

func (h *Handler) resolveStaticPath(path string) string {
	if isRemoteLocation(path) || !isStaticResourcePath(path) {
		return path
	}

	var b strings.Builder
	staticResources := strings.TrimSuffix(h.Resources.Static(), "/")
	b.WriteString(staticResources)
	path = strings.Trim(path, "/")
	b.WriteByte('/')
	b.WriteString(path)
	return b.String()
}

// Icon describes a square image that is used in various places such as
// application icon, favicon or loading icon.
type Icon struct {
	// The path or url to a square image/png file. It must have a side of 192px.
	//
	// Path is relative to the root directory.
	Default string

	// The path or url to larger square image/png file. It must have a side of
	// 512px.
	//
	// Path is relative to the root directory.
	Large string

	// The path or url to a svg file.
	SVG string

	// The path or url to a square image/png file that is used for IOS/IPadOS
	// home screen icon. It must have a side of 192px.
	//
	// Path is relative to the root directory.
	//
	// DEFAULT: Icon.Default
	AppleTouch string
}

// Environment describes the environment variables to pass to the progressive
// web app.
type Environment map[string]string

func isRemoteLocation(path string) bool {
	return strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "http://")
}

func isStaticResourcePath(path string) bool {
	return strings.HasPrefix(path, "/web/") ||
		strings.HasPrefix(path, "web/")
}

func parseSrc(link string) (url, crossOrigin, loading string) {
	for _, p := range strings.Split(link, " ") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		switch {
		case p == "crossorigin":
			crossOrigin = "true"

		case strings.HasPrefix(p, "crossorigin="):
			crossOrigin = strings.TrimPrefix(p, "crossorigin=")

		case p == "defer":
			loading = "defer"

		case p == "async":
			loading = "async"

		default:
			url = p
		}
	}

	return url, crossOrigin, loading
}
