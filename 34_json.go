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
	// Marshal visits exported fields; tags rename or omit.
	encoded, err := json.Marshal(
		payload{Name: "Ana", password: "secret"})
	if err != nil {
		panic(err)
	}
	fmt.Println(
		"encoded; false and unexported fields omitted:",
		string(encoded))

	// Unmarshal parses at runtime; no generated code.
	var decoded payload
	input := []byte(`{"name":"Bob","password":"ignored"}`)
	err = json.Unmarshal(input, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("decoded exported name:", decoded.Name)
	fmt.Println("unexported password ignored:",
		decoded.password == "")
}
