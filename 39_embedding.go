//go:build ignore

package main

import "fmt"

type Engine struct {
	Power int
}

// A method on Engine; promoted to the embedding struct.
func (engine Engine) Start() string {
	return fmt.Sprintf("engine start (%dhp)", engine.Power)
}

// Car embeds Engine; its fields and methods are promoted.
// Prefer a plain named field unless you want that promotion.
type Car struct {
	Engine
	Name string
}

// Shadowing, not overriding: no virtual dispatch into Car.
func (car Car) Start() string {
	return car.Name + ": " + car.Engine.Start()
}

func main() {
	car := Car{Engine: Engine{Power: 120}, Name: "coupe"}

	// Promoted field, then the shadowing method.
	fmt.Println("power:", car.Power)
	fmt.Println("start:", car.Start())

	// The embedded method is still reachable via the field.
	fmt.Println("inner:", car.Engine.Start())
}
