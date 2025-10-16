package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ordo_meritum/config"
	"github.com/ordo_meritum/database"
	"github.com/ordo_meritum/database/candidate_forms"
	"github.com/ordo_meritum/database/guides"
	"github.com/ordo_meritum/database/jobs"
	"github.com/ordo_meritum/database/questionnaires"
	"github.com/ordo_meritum/database/resumes"
	"github.com/ordo_meritum/database/users"
	"github.com/ordo_meritum/database/writingsamples"
	apptracking_controllers "github.com/ordo_meritum/features/application_tracking/controllers"
	apptracking_services "github.com/ordo_meritum/features/application_tracking/services"
	auth_controllers "github.com/ordo_meritum/features/auth/controllers"
	auth_services "github.com/ordo_meritum/features/auth/services"
	candidate_form_controllers "github.com/ordo_meritum/features/candidate_forms/controllers"
	candidate_form_services "github.com/ordo_meritum/features/candidate_forms/services"
	doc_controllers "github.com/ordo_meritum/features/documents/controllers"
	doc_services "github.com/ordo_meritum/features/documents/services"
	jobguide_controllers "github.com/ordo_meritum/features/job_guide/controllers"
	jobguide_services "github.com/ordo_meritum/features/job_guide/services"
	"github.com/ordo_meritum/kafka"
	"github.com/ordo_meritum/web"
	"github.com/ordo_meritum/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func main() {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "/app/logs/server"
		log.Warn().Str("service", "startup").Msgf("LOG_DIR not set, defaulting to %s", logDir)
	}

	logFileDir := filepath.Dir(logDir)
	if err := os.MkdirAll(logFileDir, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create log directory")
	}

	logFile, err := os.OpenFile(
		logFileDir+"/server.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	defer logFile.Close()

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	err = godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Warning: .env file not found. Using system environment variables.")
	}

	fx.New(
		fx.Provide(
			web.NewHTTPServer,
			mux.NewRouter,
			web.NewAuthenticatedRouter,
			web.NewSecureRouter,

			websocket.NewHub,

			config.NewPrivateKey,
			database.NewDB,

			candidate_forms.NewPostgresRepository,
			writingsamples.NewPostgresRepository,
			jobs.NewPostgresRepository,
			guides.NewPostgresRepository,
			users.NewPostgresRepository,
			questionnaires.NewPostgresRepository,
			resumes.NewPostgresRepository,

			kafka.NewLatexWriter,

			auth_services.NewAuthService,
			auth_controllers.NewController,
			candidate_form_services.NewCandidateFormService,
			candidate_form_controllers.NewController,
			apptracking_services.NewAppTrackerService,
			apptracking_controllers.NewController,
			doc_services.NewDocumentService,
			doc_controllers.NewDocumentController,
			jobguide_services.NewJobGuideService,
			jobguide_controllers.NewController,

			web.NewRouteDependencies,
		),

		fx.Invoke(web.InitializeFirebase),
		fx.Invoke(kafka.RegisterCompletionConsumer),
		fx.Invoke(web.RegisterRoutes),
		fx.Invoke(func(lc fx.Lifecycle, hub *websocket.Hub) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go hub.Run()
					log.Info().Str("service", "startup").Msg("WebSocket Hub started")
					return nil
				},
			})
		}),

		fx.Invoke(func(*http.Server) { /*intentionally left empty*/ }),
	).Run()
}
