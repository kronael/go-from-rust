//go:build ignore

package main

import "fmt"

// Rust: Drop runs automatically at the end of a scope.
// Go has no destructor hook: Close is an ordinary call,
// and defer schedules that call for function return.
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
	fmt.Println("perFunction loop ended")
}

// Narrowing the scope takes a function call, because a
// defer fires when its own function exits, and a block
// is not a function.
func perIteration(names []string) {
	for _, name := range names {
		func() {
			f := open(name)
			defer f.Close()
		}()
	}
	fmt.Println("perIteration loop ended")
}

func main() {
	fmt.Println("defer belongs to perFunction:")
	perFunction([]string{"a", "b"})
	fmt.Println("perFunction returned")

	fmt.Println("one function return per iteration:")
	perIteration([]string{"c", "d"})
	fmt.Println("perIteration returned")
}
