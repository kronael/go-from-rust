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
	// Satisfied implicitly: no impl or registration needed.
	var namer Namer = User("Ana")

	// Interface holds a dynamic type and value.
	fmt.Printf("namer: type=%T name=%q\n", namer, namer.Name())

	// *Counter's method set has Reset; Counter's does not.
	counter := Counter(7)
	var resetter interface{ Reset() } = &counter
	resetter.Reset()
	fmt.Printf("resetter: type=%T counter=%d\n",
		resetter, counter)
}
