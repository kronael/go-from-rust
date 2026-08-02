//go:build ignore

package main

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
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

	// Go 1.27: json (v1) runs on the v2 engine but
	// keeps lenient defaults; a duplicate key just
	// takes the last value, no error.
	var dup payload
	_ = json.Unmarshal(
		[]byte(`{"name":"a","name":"b"}`), &dup)
	fmt.Println("v1 duplicate key, last wins:", dup.Name)

	// encoding/json/v2 opts into stricter defaults:
	// a duplicate key is a hard error instead.
	var strict payload
	dupErr := jsonv2.Unmarshal(
		[]byte(`{"name":"a","name":"b"}`), &strict)
	fmt.Println("v2 rejects duplicate key:", dupErr != nil)
}
