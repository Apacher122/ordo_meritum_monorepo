package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ordo_meritum/config"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/rs/zerolog/log"
)

type tokenContextKey string

const VerifiedTokenKey tokenContextKey = "verifiedToken"

// Authenticate is a middleware that verifies a Firebase ID token from the Authorization header.
// It expects a "Bearer <token>" format. If the token is valid, it retrieves the user's
// UID and attaches a UserContext to the request's context before passing it to the next handler.
//
// If the Authorization header is missing, the token is invalid or expired, or the Firebase
// client fails to initialize, it writes an appropriate HTTP error (401 or 500) and stops
// the request chain.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error().
				Str("middleware", "authentication").
				Msg("Authenticate Middleware: FAILED - Authorization header missing")
			http.Error(w, "Authorization header must be provided", http.StatusUnauthorized)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		client, err := config.AuthClient()
		if err != nil {
			log.Error().
				Err(err).
				Str("middleware", "authentication").
				Msg("Authenticate Middleware: FAILED - Could not get auth client")
			http.Error(w, "Error initializing Firebase Auth client", http.StatusInternalServerError)
			return
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Error().
				Err(err).
				Str("middleware", "authentication").
				Msg("Authenticate Middleware: FAILED - Token verification error")
			http.Error(w, "Invalid or expired ID token", http.StatusUnauthorized)
			return
		}

		userCtx, ok := r.Context().Value(contexts.UserContextKey).(*contexts.UserContext)
		if !ok || userCtx == nil {
			userCtx = &contexts.UserContext{}
		}

		userCtx.Token = token
		userCtx.UID = token.UID

		ctx := context.WithValue(r.Context(), contexts.UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
