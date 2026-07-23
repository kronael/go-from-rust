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
	// Default encoding visits exported fields at runtime. Tags
	// rename or omit them; the ordinary unexported named field
	// password is ignored. No derive or generated code is
	// required.
	encoded, err := json.Marshal(
		payload{Name: "Ana", password: "secret"})
	if err != nil {
		panic(err)
	}
	fmt.Println(
		"encoded; false and unexported fields omitted:",
		string(encoded))

	// Each Unmarshal parses its input. The current
	// implementation uses reflection and caches type metadata,
	// but those are implementation details. Marshaler or
	// Unmarshaler methods can replace the default behavior.
	var decoded payload
	input := []byte(`{"name":"Bob","password":"ignored"}`)
	err = json.Unmarshal(input, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("decoded exported name:", decoded.Name)
	fmt.Println("unexported password ignored:",
		decoded.password == "")

	// Rust serde commonly derives type-specific code. Go
	// favors a runtime default; both approaches still parse or
	// produce bytes, so benchmark relevant payloads.
}
