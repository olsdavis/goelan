package world

import (
	"testing"
)

var (
	firstLocation  *Location = &Location{Orientation{0, 0}, 32, 53, 0, nil}
	secondLocation *Location = &Location{Orientation{0, 0}, 32, 53, 0, nil}
)

func BenchmarkSquaredDistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		firstLocation.DistanceSquared(secondLocation)
	}
}

func BenchmarkDistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		firstLocation.Distance(secondLocation)
	}
}
