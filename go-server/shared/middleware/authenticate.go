package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ordo_meritum/config"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/rs/zerolog/log"
)

type tokenContextKey string

const VerifiedTokenKey tokenContextKey = "verifiedToken"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("middleware", "authentication").
			Msg("--- Authenticate Middleware: Running ---")

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

		log.Info().
			Str("middleware", "authentication").
			Msg(fmt.Sprintf("Authenticate Middleware: SUCCESS - Verified token for UID: %s", token.UID))

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
