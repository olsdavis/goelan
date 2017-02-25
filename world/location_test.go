package world

import (
	"testing"
)

var (
	firstLocation  *Location = &Location{Orientation{0, 0}, SimpleLocation{128, 53, 0, nil}}
	secondLocation *Location = &Location{Orientation{0, 0}, SimpleLocation{2302, 57, -64, nil}}
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
