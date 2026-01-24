package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

// main generates a pair of RSA keys and prints them to the console.
//
// It prints the private key as a base64-encoded string, prefixed with "PRIVATE_KEY=".
// It prints the public key as a base64-encoded string, prefixed with "PUBLIC_KEY=".
//
// It is intended to be copied and pasted into a .env file.
func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	fmt.Println("### COPY AND PASTE THE FOLLOWING INTO YOUR .env FILE ###")
	fmt.Println("\n#--- PRIVATE KEY ---")
	fmt.Printf("PRIVATE_KEY=\"%s\"\n", string(privateKeyBytes))
	fmt.Println("\n#--- PUBLIC KEY ---")
	fmt.Printf("PUBLIC_KEY=\"%s\"\n", string(publicKeyPEM))
	fmt.Println("\n#########################################################")
}
