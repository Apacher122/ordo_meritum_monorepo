package middleware

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/ordo_meritum/shared/contexts"
	"github.com/rs/zerolog/log"
)

type apiKeyContextKey string

const APIKeyContextKey apiKeyContextKey = "apiKey"

func Decrypt(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	log.Info().
		Str("middleware", "decryption").
		Msg("--- Decryption Middleware: Running ---")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().
					Err(err).
					Str("middleware", "decryption").
					Msg("Decryption Middleware: FAILED - Error reading request body")
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			encryptedAPIKeyB64 := r.Header.Get("X-Encrypted-API-Key")
			if encryptedAPIKeyB64 == "" {
				log.Error().
					Str("middleware", "decryption").
					Msg("Decryption Middleware: FAILED - Missing X-Encrypted-API-Key header")
				http.Error(w, "Missing X-Encrypted-API-Key header", http.StatusBadRequest)
				return
			}

			encryptedAPIKey, err := base64.StdEncoding.DecodeString(encryptedAPIKeyB64)
			if err != nil {
				log.Error().
					Err(err).
					Str("middleware", "decryption").
					Msg("Decryption Middleware: FAILED - Error decoding API key header Base64")
				http.Error(w, "Error decoding API key header Base64", http.StatusBadRequest)
				return
			}

			apiKeyBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedAPIKey, nil)
			if err != nil {
				log.Error().
					Err(err).
					Str("middleware", "decryption").
					Msg("Decryption Middleware: FAILED - Error decrypting API key from header")
				http.Error(w, "Error decrypting API key from header", http.StatusForbidden)
				return
			}
			apiKeyStr := string(apiKeyBytes)

			log.Info().
				Str("middleware", "decryption").
				Str("apiKey", apiKeyStr).
				Msg("Decryption Middleware: SUCCESS - Decrypted API key from header")

			userCtx := &contexts.UserContext{}
			userCtx.ApiKey = apiKeyStr
			ctx := context.WithValue(r.Context(), contexts.UserContextKey, userCtx)
			r.Body = io.NopCloser(bytes.NewReader(body))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
