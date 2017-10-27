package encrypt

import (
	"crypto/rsa"
	"testing"
)

var key *rsa.PrivateKey = GeneratePrivateKey()

func BenchmarkGenerateKeyPair(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GeneratePrivateKey()
	}
}

func BenchmarkGeneratePublicKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GeneratePublicKey(key)
	}
}

func BenchmarkGenerateVerifyToken(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GenerateVerifyToken()
	}
}
