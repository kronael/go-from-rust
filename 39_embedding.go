//go:build ignore

package main

import "fmt"

type Engine struct {
	Power int
}

// A method on the embedded type; it will be promoted to the outer struct.
func (engine Engine) Start() string {
	return fmt.Sprintf("engine start (%dhp)", engine.Power)
}

// Car embeds Engine by naming the type with no field name. Engine's fields
// and methods are promoted, so car.Power and car.Start() work directly.
// Rust has no inheritance; you would hold `engine: Engine` and forward
// calls, or get default methods from a trait.
type Car struct {
	Engine
	Name string
}

// Embedding is NOT subclassing. This Start shadows Engine.Start; it does
// not override it: there is no virtual dispatch back into Car from Engine.
func (car Car) Start() string {
	return car.Name + ": " + car.Engine.Start()
}

func main() {
	car := Car{Engine: Engine{Power: 120}, Name: "coupe"}

	// Promoted field, then the shadowing method.
	fmt.Println("power:", car.Power)
	fmt.Println("start:", car.Start())

	// The embedded method is still reachable through the field name.
	fmt.Println("inner:", car.Engine.Start())
}
