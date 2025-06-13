package testutil

import (
	"math/rand"
	"time"
)

// GenerateRandomInts generates a slice of 'count' random integers,
// with each integer being in the range [0, maxVal).
// It uses a new random source for each call to ensure different sequences
// unless the test needs deterministic sequences (then seed could be a parameter).
func GenerateRandomInts(count int, maxVal int) []int {
	// It's generally good practice to create a local RNG for such utility functions
	// if you want non-deterministic sequences for different test runs or setups.
	// If deterministic sequences are needed across benchmark runs for comparability,
	// you might pass a seed or a pre-initialized *rand.Rand.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := make([]int, count)

	for i := range nums {
		nums[i] = rng.Intn(maxVal)
	}

	return nums
}

// GeneratePermutedInts generates a slice of integers from 0 to count-1
// in a random order.
func GeneratePermutedInts(count int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := rng.Perm(count)
	// rng.Perm(n) returns a permutation of integers in [0, n).
	// If you need [1, n] or other ranges, adjust accordingly.
	return p
}
