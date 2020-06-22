package app

import "context"

// Context represents a context that is tied to a UI element. It is canceled
// when the element is dismounted.
//
// It implements the context.Context interface.
//  https://golang.org/pkg/context/#Context
type Context struct {
	context.Context

	// The UI element tied to the context.
	Src UI

	// The JavaScript value of the element tied to the context. This is a
	// shorthand for:
	//  ctx.Src.JSValue()
	JSSrc Value
}
