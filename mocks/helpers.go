package mocks

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/pem"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func SetupGormMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	dbgorm, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
	}

	return db, dbgorm, mock
}
