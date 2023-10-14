package app

import "context"

type engineX struct {
	localStorage   BrowserStorage
	sessionStorage BrowserStorage
	nodes          nodeManager
}

func newEngineX(ctx context.Context, resolveURL func(string) string) *engineX {
	// var localStorage BrowserStorage
	// var sessionStorage BrowserStorage
	// if IsServer {
	// 	localStorage = newMemoryStorage()
	// 	sessionStorage = newMemoryStorage()
	// } else {
	// 	localStorage = newJSStorage("localStorage")
	// 	sessionStorage = newJSStorage("sessionStorage")
	// }

	return &engineX{
		nodes: nodeManager{},
	}
}
