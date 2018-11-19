package main

import "context"

// Package describes an app package.
type Package interface {
	// Build builds the package.
	Build(ctx context.Context) error

	// Run builds and run the package.
	Run(ctx context.Context) error
}
