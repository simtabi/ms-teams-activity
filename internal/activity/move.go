package activity

import "math/rand"

// naturalDelta returns a small, non-zero pixel offset (magnitude 1–3) in a
// random direction. Using a varied offset each tick produces less mechanical,
// more natural cursor movement than a fixed step. The activator moves by this
// delta and immediately back, so the cursor does not drift.
func naturalDelta() int {
	d := 1 + rand.Intn(3) // 1..3
	if rand.Intn(2) == 0 {
		return -d
	}
	return d
}

// naturalVertical reports whether this tick should nudge vertically instead of
// horizontally, so movement isn't always along one axis.
func naturalVertical() bool { return rand.Intn(4) == 0 }
