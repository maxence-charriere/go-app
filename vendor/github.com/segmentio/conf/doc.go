// Package conf package provides tools for easily loading program configurations
// from multiple sources such as the command line arguments, environment, or a
// configuration file.
//
// Most applications only need to use the Load function to get their settings
// loaded into an object. By default, Load will read from a configurable file
// defined by the -config-file command line argument, load values present in the
// environment, and finally load the program arguments.
//
// The object in which the configuration is loaded must be a struct, the names
// and types of its fields are introspected by the Load function to understand
// how to load the configuration.
//
// The name deduction from the struct field obeys the same rules than those
// implemented by the standard encoding/json package, which means the program
// can set the "conf" tag to override the default field names in the command
// line arguments and configuration file.
//
// A "help" tag may also be set on the fields of the configuration object to
// add documentation to the setting, which will be shown when the program is
// asked to print its help.
//
// When values are loaded from the environment the Load function looks for
// variables matching the struct fields names in snake-upper-case form.
package conf
