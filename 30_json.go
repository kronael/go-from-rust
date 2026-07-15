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
	// Rust serde uses derives. Go's json package sees exported fields and tags.
	encoded, err := json.Marshal(payload{Name: "Ana", password: "secret"})
	if err != nil {
		panic(err)
	}
	fmt.Println("encoded:", string(encoded))

	var decoded payload
	input := []byte(`{"name":"Bob","password":"ignored"}`)
	err = json.Unmarshal(input, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("decoded:", decoded.Name, "password empty:", decoded.password == "")
}
