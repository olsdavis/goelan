package world

import (
	"testing"
	"github.com/ojrac/opensimplex-go"
)
var seed int64 = 124124
var testWorld *World = &World{
	Seed: seed,
	ChunkManager: NewManager(seed),
	Type: Overworld,
}
var location SimpleLocation = SimpleLocation{316, 188, 527, nil}

func TestGenerateChunk(t *testing.T) {
	testWorld.ChunkManager.LoadChunk(location)
}

func BenchmarkChunkCoords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ChunkCoords(location)
	}
}

func BenchmarkOpenSimplex(b *testing.B) {
	noise := opensimplex.New()
	for i := 0; i < b.N; i++ {
		noise.Eval2(10, 20)
	}
}
