package main

import "context"

// Package describes an app package.
type Package interface {
	// Init creates the required directories and installs or gives instructions
	// about required frameworks and tools to build an app.
	Init(ctx context.Context) error

	// Build builds the package.
	Build(ctx context.Context) error

	// Run builds and run the package.
	Run(ctx context.Context) error

	// Clean delete the package and its temporary build files.
	Clean(ctx context.Context) error
}
