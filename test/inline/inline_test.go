package inline

import "testing"

//go:noinline
func maxNoinline(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func maxInline(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func BenchmarkMaxNoinline(b *testing.B) {
	x, y := 1, 2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maxNoinline(x, y)
	}
}

func BenchmarkInline(b *testing.B) {
	x, y := 1, 2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maxInline(x, y)
	}
}
