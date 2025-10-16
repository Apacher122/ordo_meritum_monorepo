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

type RouteDependencies struct {
	AuthController       *auth_controllers.Controller
	UserController       *user_controllers.Controller
	AppTrackerController *apptracking_controllers.Controller
	DocController        *doc_controllers.Controller
	JobGuideController   *jobguide_controllers.Controller
	WebSocketHub         *websocket.Hub
}

func NewRouteDependencies(
	authController *auth_controllers.Controller,
	userController *user_controllers.Controller,
	appTrackerController *apptracking_controllers.Controller,
	docController *doc_controllers.Controller,
	jobGuideController *jobguide_controllers.Controller,
	hub *websocket.Hub,
) *RouteDependencies {
	return &RouteDependencies{
		AuthController:       authController,
		UserController:       userController,
		AppTrackerController: appTrackerController,
		DocController:        docController,
		JobGuideController:   jobGuideController,
		WebSocketHub:         hub,
	}
}

func RegisterRoutes(
	mainRouter *mux.Router,
	authenticatedRouter *AuthenticatedRouter,
	secureRouter *SecureRouter,
	deps *RouteDependencies,
	hub *websocket.Hub,
) {
	log.Info().Str("service", "startup").Msg("Registering feature routes")

	wsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	mainRouter.Handle("/ws", wsHandler)

	authenticatedRouter.HandleFunc("/downloads", HandleDownload).Methods("POST")
	mainRouter.HandleFunc("/public-key", security.GetPublicKeyHandler).Methods("GET")
	mainRouter.HandleFunc("/public-key-stream", security.PublicKeyStreamHandler)

	deps.AuthController.RegisterRoutes(authenticatedRouter.PathPrefix("/").Subrouter())
	deps.UserController.RegisterRoutes(secureRouter.PathPrefix("/user").Subrouter())
	deps.AppTrackerController.RegisterRoutes(secureRouter.Router, authenticatedRouter.Router)
	deps.DocController.RegisterRoutes(secureRouter.Router)
	deps.JobGuideController.RegisterRoutes(secureRouter.Router, authenticatedRouter.Router)
}
