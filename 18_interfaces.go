//go:build ignore

package main

import "fmt"

type Namer interface {
	Name() string
}

type User string

func (user User) Name() string { return string(user) }

type Counter int

func (counter *Counter) Reset() { *counter = 0 }

func main() {
	// At assignment, the compiler checks that User's method
	// set contains Name. No impl declaration or runtime
	// registration is needed.
	var namer Namer = User("Ana")

	// The interface stores a dynamic concrete type and value.
	// Method calls dispatch through that type, like a Rust dyn
	// Trait call.
	fmt.Printf("namer: type=%T name=%q\n", namer, namer.Name())

	// Counter's method set lacks Reset; *Counter's method set
	// includes it. Therefore Counter does not satisfy this
	// interface, but *Counter does.
	counter := Counter(7)
	var resetter interface{ Reset() } = &counter
	resetter.Reset()
	fmt.Printf("resetter: type=%T counter=%d\n",
		resetter, counter)

	// Interface conversion and dispatch do not guarantee an
	// allocation or a fixed cost: the compiler may keep values
	// inline, devirtualize, or allocate. Prefer concrete types
	// in hot paths and benchmark instead of guessing.
}
