//go:build ignore

package main

import (
	"errors"
	"fmt"
)

// Rust: private fields plus Ratio::new() -> Result<Ratio>
// can force other modules through the constructor. In Go,
// every type also has a zero value that needs no constructor.
type Ratio struct{ num, den int }

func NewRatio(num, den int) (Ratio, error) {
	if den == 0 {
		return Ratio{}, errors.New("zero denominator")
	}
	return Ratio{num, den}, nil
}

func (r Ratio) Float() float64 {
	return float64(r.num) / float64(r.den)
}

// Unexported fields stop other packages from setting
// fields, never from creating a zeroed value. There is
// no way to make a type require initialization.
type Counter struct {
	seen map[string]int
	log  []string
}

// So the idiom inverts: make the zero value correct
// rather than defend a constructor nothing has to call.
// Appending to a nil slice works, so Counter needs no New.
func (c *Counter) Add(word string) {
	if c.seen == nil {
		c.seen = map[string]int{}
	}
	c.seen[word]++
	c.log = append(c.log, word)
}

func main() {
	ok, err := NewRatio(3, 4)
	fmt.Println("constructed:", ok.Float(), err)

	_, err = NewRatio(1, 0)
	fmt.Println("rejected:", err)

	// This declaration creates Ratio{0, 0} without NewRatio.
	var zero Ratio
	fmt.Println("var skips constructor:", zero, zero.Float())

	var c Counter
	c.Add("go")
	c.Add("go")
	fmt.Println("zero value made useful:", c.seen["go"], c.log)
}
