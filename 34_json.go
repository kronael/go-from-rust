//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
)

type payload struct {
	Name     string `json:"name"`
	Admin    bool   `json:"admin,omitempty"`
	password string
}

func main() {
	// encoding/json generates no type-specific code. On first use it reflects on
	// a type and caches its encoder plus field and tag metadata by reflect.Type.
	encoded, err := json.Marshal(payload{Name: "Ana", password: "secret"})
	if err != nil {
		panic(err)
	}
	fmt.Println("encoded; false and unexported fields omitted:", string(encoded))

	// Every Marshal still walks the value at runtime and returns a fresh []byte.
	// Every Unmarshal parses its input. Default struct decoding writes through
	// reflection; a custom json.Unmarshaler can take over. Cached type metadata
	// avoids rediscovering this struct on every call.
	var decoded payload
	input := []byte(`{"name":"Bob","password":"ignored"}`)
	err = json.Unmarshal(input, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("decoded exported name:", decoded.Name)
	fmt.Println("unexported password ignored:", decoded.password == "")

	// Rust serde derives type-specific code at compile time, usually avoiding
	// this reflection dispatch. Both approaches still parse and produce bytes.
}
