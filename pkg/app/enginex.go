package app

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type engineX struct {
	routes          *router
	internalURLs    []string
	resolveURL      func(string) string
	originPage      requestPage
	lastVisitedURL  *url.URL
	localStorage    BrowserStorage
	sessionStorage  BrowserStorage
	nodes           nodeManager
	newBody         func() HTMLBody
	body            UI
	initBrowserOnce sync.Once
	browser         browser
}

func newEngineX(ctx context.Context, routes *router, resolveURL func(string) string, origin *url.URL, newBody func() HTMLBody) *engineX {
	var localStorage BrowserStorage
	var sessionStorage BrowserStorage
	if IsServer {
		localStorage = newMemoryStorage()
		sessionStorage = newMemoryStorage()
	} else {
		localStorage = newJSStorage("localStorage")
		sessionStorage = newJSStorage("sessionStorage")
	}

	if resolveURL == nil {
		resolveURL = func(v string) string { return v }
	}

	return &engineX{
		routes:     routes,
		resolveURL: resolveURL,
		originPage: requestPage{
			url:                   origin,
			resolveStaticResource: resolveURL,
		},
		localStorage:   localStorage,
		lastVisitedURL: &url.URL{},
		sessionStorage: sessionStorage,
		newBody:        newBody,
		nodes:          nodeManager{},
	}
}

func (e *engineX) Navigate(destination *url.URL, updateHistory bool) {
	e.initBrowserOnce.Do(e.initBrowser)

	switch {
	case e.internalURL(destination),
		e.mailTo(destination):
		Window().Get("location").Set("href", destination.String())
		return

	case e.externalNavigation(destination):
		Window().Call("open", destination.String())
		return
	}

	if destination.String() == e.lastVisitedURL.String() {
		return
	}
	defer func() {
		if updateHistory {
			Window().addHistory(destination)
		}
		e.lastVisitedURL = destination

		e.nodes.NotifyComponentEvent(e.baseContext(), e.body, nav{})
		if destination.Fragment != "" {
			e.defere(func() {
				Window().ScrollToID(destination.Fragment)
			})
		}
	}()

	if destination.Path == e.lastVisitedURL.Path &&
		destination.Fragment != e.lastVisitedURL.Fragment {
		return
	}

	path := strings.TrimPrefix(destination.Path, Getenv("GOAPP_ROOT_PREFIX"))
	if path == "" {
		path = "/"
	}
	root, ok := e.routes.createComponent(path)
	if !ok {
		root = &notFound{}
	}
	e.load(root)
}

func (e *engineX) initBrowser() {
	if IsServer {
		return
	}
	e.browser.HandleEvents(e.baseContext())
}

func (e *engineX) externalNavigation(v *url.URL) bool {
	return v.Host != e.originPage.URL().Host
}

func (e *engineX) mailTo(v *url.URL) bool {
	return v.Scheme == "mailto"
}

func (e *engineX) internalURL(v *url.URL) bool {
	if e.internalURLs == nil {
		json.Unmarshal([]byte(Getenv("GOAPP_INTERNAL_URLS")), &e.internalURLs)
	}

	url := v.String()
	for _, u := range e.internalURLs {
		if strings.HasPrefix(url, u) {
			return true
		}
	}
	return false
}

func (e *engineX) baseContext() Context {
	return nodeContext{
		resolveURL:             e.resolveURL,
		page:                   e.page,
		navigate:               e.Navigate,
		localStorage:           e.localStorage,
		sessionStorage:         e.sessionStorage,
		dispatch:               e.dispatch,
		defere:                 e.defere,
		updateComponent:        func(c Composer) {},
		preventComponentUpdate: func(c Composer) {},
	}
}

func (e *engineX) page() Page {
	if IsClient {
		return browserPage{resolveStaticResource: e.resolveURL}
	}
	return &e.originPage
}

func (e *engineX) load(v Composer) {
	if e.body == nil {
		body, err := e.nodes.Mount(e.baseContext(), 0, e.newBody().privateBody(v))
		if err != nil {
			panic(errors.New("mounting root failed").Wrap(err))
		}
		e.body = body
		return
	}

	body, err := e.nodes.Update(e.baseContext(), e.body, e.newBody().privateBody(v))
	if err != nil {
		panic(errors.New("updating root failed").Wrap(err))
	}
	e.body = body
}

func (e *engineX) dispatch(v func()) {
	// TODO implementd
	v()
}

func (e *engineX) defere(v func()) {
	// TODO implement
	v()
}
