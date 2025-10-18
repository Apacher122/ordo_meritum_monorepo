package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/features/documents/services"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/middleware"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	lg "github.com/ordo_meritum/shared/utils/logger"
	"github.com/rs/zerolog/log"
)

var service = "documents-service"

var logger = log.With().Str("service", "documents-controller").Logger()

type Controller struct {
	docService *services.DocumentService
}

func NewDocumentController(docService *services.DocumentService) *Controller {
	return &Controller{docService: docService}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/documents/resume", c.generateDocumentHandler(c.docService.QueueResumeGeneration)).Methods("POST")
	router.HandleFunc("/documents/cover-letter", c.generateDocumentHandler(c.docService.QueueCoverLetterGeneration)).Methods("POST")
}

func (c *Controller) generateDocumentHandler(
	generationFunc func(
		ctx context.Context,
		requestBody requests.DocumentRequest,
	) (int, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		_, ok := contexts.FromContext(r.Context())
		if !ok {
			lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
			middleware.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		requestBody, err := decodeDocumentRequest(r)
		if err != nil {
			handleDecodeError(w, err)
			return
		}

		jobID, err := generationFunc(r.Context(), requestBody)
		if err != nil {
			lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_DB_FAILED_TO_INSERT, Error: err}.ErrorLog()
			middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to queue document for generation"})
			return
		}

		logger.Info().Int("jobID", jobID).Msg("Document queued for generation")
		middleware.JSON(w, http.StatusAccepted, map[string]any{
			"jobId":  jobID,
			"status": "processing_queued",
		})
	}
}

func decodeDocumentRequest(r *http.Request) (requests.DocumentRequest, error) {
	var requestBody requests.DocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return requests.DocumentRequest{}, err
	}
	return requestBody, nil
}

func handleDecodeError(w http.ResponseWriter, err error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": msg})

	case errors.As(err, &unmarshalTypeError):
		msg := fmt.Sprintf("Invalid type for field '%s'. Expected '%s' but received a '%s'.", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": msg})

	default:
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
	}
}
