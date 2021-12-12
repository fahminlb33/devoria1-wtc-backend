package mocks

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
)

func GenerateRSAKeyPair() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Extract public component.
	publicKey := privateKey.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	publicKeyPacked, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyPacked,
		},
	)

	return base64.StdEncoding.EncodeToString(publicKeyPem), base64.StdEncoding.EncodeToString(privateKeyPem)
}
