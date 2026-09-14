//go:build ignore

package main

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
)

type Item struct {
	Name string `json:"name"`
}

func main() {
	// Same object key appears twice.
	dup := []byte(`{"name":"a","name":"b"}`)

	// encoding/json now runs on the v2
	// engine but keeps v1's lenient
	// defaults: last value wins, no error.
	var v1 Item
	err1 := json.Unmarshal(dup, &v1)
	fmt.Println("v1 name:", v1.Name)
	fmt.Println("v1 err:", err1)

	// encoding/json/v2 is stricter by
	// default: a duplicate key is an error.
	var v2 Item
	err2 := jsonv2.Unmarshal(dup, &v2)
	fmt.Println("v2 name:", v2.Name)
	fmt.Println("v2 err:", err2)
}
