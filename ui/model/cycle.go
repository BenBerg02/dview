package model

// cycle moves i by d within [0, n), wrapping around at both ends.
func cycle(i, d, n int) int { return (i + d + n) % n }
