Package mold
============
![Project status](https://img.shields.io/badge/version-2.2.0-green.svg)
[![Build Status](https://travis-ci.org/go-playground/mold.svg?branch=v2)](https://travis-ci.org/go-playground/mold)
[![Coverage Status](https://coveralls.io/repos/github/go-playground/mold/badge.svg?branch=v2)](https://coveralls.io/github/go-playground/mold?branch=v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/mold)](https://goreportcard.com/report/github.com/go-playground/mold)
[![GoDoc](https://godoc.org/gopkg.in/go-playground/mold.v2?status.svg)](https://godoc.org/gopkg.in/go-playground/mold.v2)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Package mold is a general library to help modify or set data within data structures and other objects.

How can this help me you ask, please see the examples [here](_examples/full/main.go)

Installation
------------

Use go get.
```shell
go get -u gopkg.in/go-playground/mold.v2
```

Then import the form package into your own code.

	import "gopkg.in/go-playground/mold.v2"

Simple example
-----
```go
package main

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"gopkg.in/go-playground/mold.v2"
)

var tform *mold.Transformer

func main() {
	tform = mold.New()
	tform.Register("set", transformMyData)

	type Test struct {
		String string `mold:"set"`
	}

	var tt Test

	err := tform.Struct(context.Background(), &tt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", tt)

	var myString string
	err = tform.Field(context.Background(), &myString, "set")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(myString)
}

func transformMyData(ctx context.Context, t *mold.Transformer, value reflect.Value, param string) error {
	value.SetString("test")
	return nil
}
```

Full example
-----
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/go-playground/form"

	"gopkg.in/go-playground/mold.v2/modifiers"
	"gopkg.in/go-playground/mold.v2/scrubbers"

	"gopkg.in/go-playground/validator.v9"
)

// This example is centered around a form post, but doesn't have to be
// just trying to give a well rounded real life example.

// <form method="POST">
//   <input type="text" name="Name" value="joeybloggs"/>
//   <input type="text" name="Age" value="3"/>
//   <input type="text" name="Gender" value="Male"/>
//   <input type="text" name="Address[0].Name" value="26 Here Blvd."/>
//   <input type="text" name="Address[0].Phone" value="9(999)999-9999"/>
//   <input type="text" name="Address[1].Name" value="26 There Blvd."/>
//   <input type="text" name="Address[1].Phone" value="1(111)111-1111"/>
//   <input type="text" name="active" value="true"/>
//   <input type="submit"/>
// </form>

var (
	conform  = modifiers.New()
	scrub    = scrubbers.New()
	validate = validator.New()
	decoder  = form.NewDecoder()
)

// Address contains address information
type Address struct {
	Name  string `mod:"trim" validate:"required"`
	Phone string `mod:"trim" validate:"required"`
}

// User contains user information
type User struct {
	Name    string    `mod:"trim"      validate:"required"              scrub:"name"`
	Age     uint8     `                validate:"required,gt=0,lt=130"`
	Gender  string    `                validate:"required"`
	Email   string    `mod:"trim"      validate:"required,email"        scrub:"emails"`
	Address []Address `                validate:"required,dive"`
	Active  bool      `form:"active"`
}

func main() {
	// this simulates the results of http.Request's ParseForm() function
	values := parseForm()

	var user User

	// must pass a pointer
	err := decoder.Decode(&user, values)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Decoded:%+v\n\n", user)

	// great not lets conform our values, after all a human input the data
	// nobody's perfect
	err = conform.Struct(context.Background(), &user)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Conformed:%+v\n\n", user)

	// that's better all those extra spaces are gone
	// let's validate the data
	err = validate.Struct(user)
	if err != nil {
		log.Panic(err)
	}

	// ok now we know our data is good, let's do something with it like:
	// save to database
	// process request
	// etc....

	// ok now I'm done working with my data
	// let's log or store it somewhere
	// oh wait a minute, we have some sensitive PII data
	// let's make sure that's de-identified first
	err = scrub.Struct(context.Background(), &user)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Scrubbed:%+v\n\n", user)
}

// this simulates the results of http.Request's ParseForm() function
func parseForm() url.Values {
	return url.Values{
		"Name":             []string{"  joeybloggs  "},
		"Age":              []string{"3"},
		"Gender":           []string{"Male"},
		"Email":            []string{"Dean.Karn@gmail.com  "},
		"Address[0].Name":  []string{"26 Here Blvd."},
		"Address[0].Phone": []string{"9(999)999-9999"},
		"Address[1].Name":  []string{"26 There Blvd."},
		"Address[1].Phone": []string{"1(111)111-1111"},
		"active":           []string{"true"},
	}
}
```

Special Information
-------------------
- To use a comma(,) within your params replace use it's hex representation instead '0x2C' which will be replaced while caching.

Contributing
------------
I am definitly interested in the communities help in adding more scrubbers and modifiers.
Please send a PR with tests, and prefereably no extra dependencies, at lease until a solid base
has been built.

Complimentary Software
----------------------

Here is a list of software that compliments using this library post decoding.

* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving.
* [form](https://github.com/go-playground/form) - Decodes url.Values into Go value(s) and Encodes Go value(s) into url.Values. Dual Array and Full map support.

License
------
Distributed under MIT License, please see license file in code for more details.
