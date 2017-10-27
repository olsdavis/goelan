package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

// Generates a 1024 bits RSA private key.
func GeneratePrivateKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic("Could not create private key")
	}
	err = key.Validate()
	if err != nil {
		panic("Could not validate key")
	}
	return key
}

// GeneratePublicKey generates the public key from the
// RSA private key.
func GeneratePublicKey(key *rsa.PrivateKey) []byte {
	pub := &key.PublicKey
	pk, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		panic("Could not generate the public key")
	}
	return pk
}

// GenerateVerifyToken generates a verify token used for
// players' authentication.
func GenerateVerifyToken() []byte {
	buf := make([]byte, 4)
	rand.Read(buf)
	return buf
}
