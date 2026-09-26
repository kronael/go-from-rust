//go:build ignore

package main

import (
	"fmt"
	"simd"
)

// Rust's portable SIMD is nightly-only: on stable 1.97
// `use std::simd` is E0658. Go 1.27 ships package simd
// behind an experiment, so this lesson runs only as
//
//	GOEXPERIMENT=simd go run 47_simd.go
//
// Without it the build fails before anything runs. The
// public Playground cannot set the experiment, so the
// tour's Run fails here with a build error.

// The scalar baseline: one value per step.
func sum(xs []float32) float32 {
	total := float32(0)
	for _, x := range xs {
		total += x
	}
	return total
}

func main() {
	xs := make([]float32, 18)
	for i := range xs {
		xs[i] = float32(i)
	}

	// Len is a run-time value, not a const generic like
	// f32x4. The compiler emits one copy of main per
	// width, and the binary picks one at startup from what
	// the CPU supports, so the width and lanes lines differ
	// between machines. GODEBUG=simd=128 forces width 4 on
	// the same binary. The zero vector is valid.
	var acc simd.Float32s
	width := acc.Len()

	// Each Add covers every lane, in one instruction where
	// the CPU has one. No unsafe block appears, because
	// the package emulates in pure Go where it does not.
	rest := xs
	for len(rest) >= width {
		acc = acc.Add(simd.LoadFloat32s(rest))
		rest = rest[width:]
	}

	// Every width is a multiple of 4 and 18 is not, so a
	// tail is always left. A short load zero-fills the
	// missing lanes.
	tail, _ := simd.LoadFloat32sPart(rest)
	acc = acc.Add(tail)

	// The package has no horizontal sum like Rust's
	// reduce_sum, so the lanes go back to a slice.
	lanes := make([]float32, width)
	acc.Store(lanes)

	fmt.Println("width:", width)
	fmt.Println("lanes:", acc)
	fmt.Println("vector sum:", sum(lanes))
	fmt.Println("scalar sum:", sum(xs))
}
