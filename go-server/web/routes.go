package web

import (
	"net/http"

	"github.com/gorilla/mux"
	apptracking_controllers "github.com/ordo_meritum/features/application_tracking/controllers"
	auth_controllers "github.com/ordo_meritum/features/auth/controllers"
	user_controllers "github.com/ordo_meritum/features/candidate_forms/controllers"
	doc_controllers "github.com/ordo_meritum/features/documents/controllers"
	jobguide_controllers "github.com/ordo_meritum/features/job_guide/controllers"
	"github.com/ordo_meritum/security"
	"github.com/ordo_meritum/websocket"
	"github.com/rs/zerolog/log"
)

func RegisterRoutes(
	mainRouter *mux.Router,
	authenticatedRouter *AuthenticatedRouter,
	secureRouter *SecureRouter,
	authController *auth_controllers.Controller,
	userController *user_controllers.Controller,
	appTrackerController *apptracking_controllers.Controller,
	docController *doc_controllers.Controller,
	jobGuideController *jobguide_controllers.Controller,
	hub *websocket.Hub,
) {
	log.Info().Str("service", "startup").Msg("Registering feature routes")

	// mainRouter.Use(loggingMiddleware)
	wsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	mainRouter.Handle("/ws", wsHandler)

	authenticatedRouter.HandleFunc("/downloads", HandleDownload).Methods("POST")
	mainRouter.HandleFunc("/public-key", security.GetPublicKeyHandler).Methods("GET")
	mainRouter.HandleFunc("/public-key-stream", security.PublicKeyStreamHandler)

	authController.RegisterRoutes(authenticatedRouter.PathPrefix("/").Subrouter())
	userController.RegisterRoutes(secureRouter.PathPrefix("/user").Subrouter())
	appTrackerController.RegisterRoutes(authenticatedRouter.PathPrefix("/apps").Subrouter())
	docController.RegisterRoutes(secureRouter.Router)
	jobGuideController.RegisterRoutes(secureRouter.PathPrefix("/guide").Subrouter())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Interface("headers", r.Header).
			Msg("Received request")
		next.ServeHTTP(w, r)
	})
}
