//go:build ignore

package main

import "fmt"

// Rust's portable SIMD is nightly-only: on stable 1.97
// `use std::simd` is E0658. Go 1.27 ships package simd
// behind GOEXPERIMENT=simd. Neither language hands you
// portable vectors on its stable default.

// This is the code that runs here, and it is plain Go.
// The tour executes on the public Playground, which
// cannot set the experiment.
func sum(xs []float32) float32 {
	total := float32(0)
	for _, x := range xs {
		total += x
	}
	return total
}

// The vector version, not run here. It compiles under
// GOEXPERIMENT=simd and agrees with sum above:
//
//	var acc simd.Float32s  // zero vector is valid
//	width := acc.Len()
//	for len(xs) >= width {
//		acc = acc.Add(simd.LoadFloat32s(xs))
//		xs = xs[width:]
//	}
//	tail, _ := simd.LoadFloat32sPart(xs)
//	acc = acc.Add(tail)    // short load zero-fills
//	out := make([]float32, width)
//	acc.Store(out)
//	return sum(out)        // reduce the lanes
//
// Two things differ from Rust. Len is a run-time value,
// not a const generic like f32x4, so one binary adapts
// to the machine it lands on. And no unsafe block
// appears, because the package falls back to pure Go
// where the hardware is missing.

func main() {
	xs := make([]float32, 20)
	for i := range xs {
		xs[i] = float32(i)
	}

	// 20 values is not a whole number of vectors, which
	// is why the version above needs the tail load.
	fmt.Println("sum:", sum(xs))
	fmt.Println("computed one value at a time, no simd")
}
