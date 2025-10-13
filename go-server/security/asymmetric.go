package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ordo_meritum/shared/middleware"
)

func GenerateKeys(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("cannot generate RSA key: %w", err)
	}
	return privateKey, nil
}

func PrivateKeyToPEM(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
}

func PublicKeyToPEM(pub *rsa.PublicKey) ([]byte, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubBytes,
		},
	), nil
}

func GetPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting public key")
	publicKeyStr := os.Getenv("PUBLIC_KEY")
	if publicKeyStr == "" {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Public key not found on server"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"publicKey": publicKeyStr})
}
