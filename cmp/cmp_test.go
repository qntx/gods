package cmp_test

import (
	"cmp"
	"math"
	"testing"
	"time"

	godscmp "github.com/qntx/gods/cmp"
)

// TestTimeComparator verifies TimeComparator's behavior with time.Time values.
//
// Ensures correct ordering using time.Time's After and Before methods.
func TestTimeComparator(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name string
		t1   time.Time
		t2   time.Time
		want int
	}{
		{name: "equal", t1: now, t2: now, want: 0},
		{name: "t1 > t2", t1: now.Add(2 * 7 * 24 * time.Hour), t2: now, want: 1},
		{name: "t1 < t2", t1: now, t2: now.Add(2 * 7 * 24 * time.Hour), want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := godscmp.TimeComparator(tt.t1, tt.t2)
			if got != tt.want {
				t.Errorf("TimeComparator(%v, %v) = %d, want %d", tt.t1, tt.t2, got, tt.want)
			}
		})
	}
}

// TestFloat64Comparator verifies Float64Comparator's behavior with float64 values.
//
// Tests include equality within epsilon, strict ordering, and special cases like NaN and ±0.
func TestFloat64Comparator(t *testing.T) {
	t.Parallel()

	// Compute at runtime to preserve IEEE 754 precision behavior.
	a := 0.1
	b := 0.2
	sum := a + b // ≈ 0.30000000000000004

	const epsilon = 1e-10

	tests := []struct {
		name    string
		x       float64
		y       float64
		epsilon float64
		want    int
	}{
		{name: "equal", x: 1.0, y: 1.0, epsilon: epsilon, want: 0},
		{name: "approx equal", x: sum, y: 0.3, epsilon: epsilon, want: 0},
		{name: "x > y", x: 2.0, y: 1.0, epsilon: epsilon, want: 1},
		{name: "x < y", x: 1.0, y: 2.0, epsilon: epsilon, want: -1},
		{name: "zero vs neg zero", x: 0.0, y: math.Copysign(0, -1), epsilon: epsilon, want: 0},
		{name: "NaN vs NaN", x: math.NaN(), y: math.NaN(), epsilon: epsilon, want: 0},
		{name: "NaN < non-NaN", x: math.NaN(), y: 1.0, epsilon: epsilon, want: -1},
		{name: "non-NaN > NaN", x: 1.0, y: math.NaN(), epsilon: epsilon, want: 1},
		{name: "invalid epsilon", x: 1.0, y: 1.1, epsilon: -1, want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := godscmp.Float64Comparator(tt.x, tt.y, tt.epsilon)
			if got != tt.want {
				t.Errorf("Float64Comparator(%v, %v, %v) = %d, want %d", tt.x, tt.y, tt.epsilon, got, tt.want)
			}
		})
	}
}

// TestFloat64ReverseComparator verifies Float64ReverseComparator's behavior.
//
// Ensures reverse ordering (descending) with epsilon tolerance, including NaN and ±0 cases.
func TestFloat64ReverseComparator(t *testing.T) {
	t.Parallel()

	// Compute at runtime to preserve IEEE 754 precision behavior.
	a := 0.1
	b := 0.2
	sum := a + b // ≈ 0.30000000000000004

	const epsilon = 1e-10

	tests := []struct {
		name    string
		x       float64
		y       float64
		epsilon float64
		want    int
	}{
		{name: "equal", x: 1.0, y: 1.0, epsilon: epsilon, want: 0},
		{name: "approx equal", x: sum, y: 0.3, epsilon: epsilon, want: 0},
		{name: "x > y (reverse)", x: 2.0, y: 1.0, epsilon: epsilon, want: -1},
		{name: "x < y (reverse)", x: 1.0, y: 2.0, epsilon: epsilon, want: 1},
		{name: "zero vs neg zero", x: 0.0, y: math.Copysign(0, -1), epsilon: epsilon, want: 0},
		{name: "NaN vs NaN", x: math.NaN(), y: math.NaN(), epsilon: epsilon, want: 0},
		{name: "NaN > non-NaN", x: math.NaN(), y: 1.0, epsilon: epsilon, want: 1},
		{name: "non-NaN < NaN", x: 1.0, y: math.NaN(), epsilon: epsilon, want: -1},
		{name: "invalid epsilon", x: 1.0, y: 1.1, epsilon: -1, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := godscmp.Float64ReverseComparator(tt.x, tt.y, tt.epsilon)
			if got != tt.want {
				t.Errorf("Float64ReverseComparator(%v, %v, %v) = %d, want %d", tt.x, tt.y, tt.epsilon, got, tt.want)
			}
		})
	}
}

// TestCmpFloat64Compare verifies cmp.Compare's behavior with float64 values.
//
// Highlights strict comparison without epsilon, including NaN and ±0 cases.
func TestCmpFloat64Compare(t *testing.T) {
	t.Parallel()

	// Compute at runtime to preserve IEEE 754 precision behavior.
	a := 0.1
	b := 0.2
	sum := a + b // ≈ 0.30000000000000004

	tests := []struct {
		name string
		x    float64
		y    float64
		want int
	}{
		{name: "equal", x: 1.0, y: 1.0, want: 0},
		{name: "sum > 0.3", x: sum, y: 0.3, want: 1},
		{name: "0.3 < sum", x: 0.3, y: sum, want: -1},
		{name: "x > y", x: 2.0, y: 1.0, want: 1},
		{name: "x < y", x: 1.0, y: 2.0, want: -1},
		{name: "zero vs neg zero", x: 0.0, y: math.Copysign(0, -1), want: 0},
		{name: "NaN vs NaN", x: math.NaN(), y: math.NaN(), want: 0},
		{name: "NaN < non-NaN", x: math.NaN(), y: 1.0, want: -1},
		{name: "non-NaN > NaN", x: 1.0, y: math.NaN(), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := cmp.Compare(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("cmp.Compare(%v, %v) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}
