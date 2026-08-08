//go:build ignore

package main

import "fmt"

// Rust: Drop runs at end of scope, guaranteed, and
// composes through ownership without a call site. Go has
// no destructors — cleanup is a call someone must make.
type File struct{ name string }

func open(name string) *File {
	fmt.Println("open", name)
	return &File{name}
}

func (f *File) Close() { fmt.Println("close", f.name) }

// defer is function-scoped, not block-scoped: every File
// stays open until the function returns. In Rust each one
// would drop at the end of its iteration.
func perFunction(names []string) {
	for _, name := range names {
		f := open(name)
		defer f.Close()
	}
	fmt.Println("loop done, nothing closed yet")
}

// Narrowing the scope takes a function call, because a
// function return is the only thing that fires a defer.
func perIteration(names []string) {
	for _, name := range names {
		func() {
			f := open(name)
			defer f.Close()
		}()
	}
	fmt.Println("loop done, all closed")
}

func main() {
	perFunction([]string{"a", "b"})
	perIteration([]string{"c", "d"})

	// runtime.AddCleanup and SetFinalizer are not a
	// substitute: no ordering, and no guarantee they run
	// before the process exits. Nothing here uses them.
}
