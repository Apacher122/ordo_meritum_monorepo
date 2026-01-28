package web

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ordo_meritum/config"
	"github.com/ordo_meritum/shared/middleware"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type SecureRouter struct{ *mux.Router }

type AuthenticatedRouter struct{ *mux.Router }

// InitializeFirebase sets up an fx.Lifecycle hook to initialize the global Firebase application
// during application startup.
func InitializeFirebase(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().
				Str("service", "startup").
				Msg("Initializing Firebase App")
			config.InitializeFirebaseApp()
			return nil
		},
	})
}

// NewHTTPServer creates a new http.Server instance configured with the main application router
// and integrates its lifecycle (start/stop) with the provided fx.Lifecycle.
// It listens on the address specified within the function (currently ":8080").
//
// Parameters:
//   - lc: The fx.Lifecycle instance used to manage application startup and shutdown hooks.
//   - router: The main mux.Router containing all application routes and middleware.
//
// Returns:
//   - A configured *http.Server instance.
func NewHTTPServer(lc fx.Lifecycle, router *mux.Router) *http.Server {

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				listener, err := net.Listen("tcp", server.Addr)
				if err != nil {
					return err
				}
				go server.Serve(listener)
				log.Info().
					Str("service", "startup").
					Msg(fmt.Sprintf("HTTP server listening on %s", server.Addr))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().
					Str("service", "shutdown").
					Msg("Stopping HTTP server.")
				return server.Shutdown(ctx)
			},
		},
	)

	return server
}

// NewAuthenticatedRouter creates and configures a subrouter specifically for authenticated endpoints.
// It mounts under the "/api/auth" path prefix relative to the main router and applies
// the authentication middleware (middleware.Authenticate) to all routes defined under it.
// Intended for use as an fx provider to inject an AuthenticatedRouter instance.
//
// Parameters:
//   - mainRouter: The root mux.Router of the application.
//
// Returns:
//   - An *AuthenticatedRouter instance wrapping the configured subrouter.
func NewAuthenticatedRouter(mainRouter *mux.Router) *AuthenticatedRouter {
	log.Info().
		Str("service", "startup").
		Msg("Configuring authenticated-only router...")
	authenticatedRouter := mainRouter.PathPrefix("/api/auth").Subrouter()
	authenticatedRouter.Use(middleware.Authenticate)
	return &AuthenticatedRouter{Router: authenticatedRouter}
}

// NewSecureRouter creates and configures a subrouter for endpoints requiring both decryption and authentication.
// It mounts under the "/api/secure" path prefix relative to the main router.
// It applies decryption middleware (middleware.Decrypt) first, followed by authentication
// middleware (middleware.Authenticate) to all routes defined under it.
// Intended for use as an fx provider to inject a SecureRouter instance.
//
// Parameters:
//   - mainRouter: The root mux.Router of the application.
//   - privateKey: The RSA private key used for the decryption middleware.
//
// Returns:
//   - A *SecureRouter instance wrapping the configured subrouter.
func NewSecureRouter(mainRouter *mux.Router, privateKey *rsa.PrivateKey) *SecureRouter {
	log.Info().
		Str("service", "startup").
		Msg("Configuring secure (decrypt & auth) router...")
	secureRouter := mainRouter.PathPrefix("/api/secure").Subrouter()
	secureRouter.Use(middleware.Decrypt(privateKey))
	secureRouter.Use(middleware.Authenticate)
	return &SecureRouter{Router: secureRouter}
}
