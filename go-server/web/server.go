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

func NewAuthenticatedRouter(mainRouter *mux.Router) *AuthenticatedRouter {
	log.Info().
		Str("service", "startup").
		Msg("Configuring authenticated-only router...")
	authenticatedRouter := mainRouter.PathPrefix("/api/auth").Subrouter()
	authenticatedRouter.Use(middleware.Authenticate)
	return &AuthenticatedRouter{Router: authenticatedRouter}
}

func NewSecureRouter(mainRouter *mux.Router, privateKey *rsa.PrivateKey) *SecureRouter {
	log.Info().
		Str("service", "startup").
		Msg("Configuring secure (decrypt & auth) router...")
	secureRouter := mainRouter.PathPrefix("/api/secure").Subrouter()
	secureRouter.Use(middleware.Decrypt(privateKey))
	secureRouter.Use(middleware.Authenticate)
	return &SecureRouter{Router: secureRouter}
}
