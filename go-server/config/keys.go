package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"
)

func NewPrivateKey() (*rsa.PrivateKey, error) {
	pemKey := os.Getenv("PRIVATE_KEY")
	if pemKey == "" {
		return nil, errors.New("PRIVATE_KEY environment variable not set")
	}

	pemKey = strings.ReplaceAll(pemKey, "\\n", "\n")

	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
