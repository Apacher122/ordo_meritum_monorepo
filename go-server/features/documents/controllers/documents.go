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

var logger = log.With().Str("service", "documents-controller").Logger()
var service = "documents-controller"

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

// generateDocumentHandler generates a handler function for queuing a document for generation based on the generationFunc provided.
// It takes in a request body and returns a jobID of the queued document and an error if any errors occur.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the request body is malformed, it will return an error with ErrorCode set to ERR_INVALID_REQUEST_BODY.
// If the document cannot be queued, it will return an error with ErrorCode set to ERR_FAILED_TO_QUEUE_DOCUMENT.
// It will log a message with the jobID and the status of the document.
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

// decodeDocumentRequest is a helper function that decodes a JSON request body into a DocumentRequest.
// If the decoding fails, it returns an error with the relevant error message.
func decodeDocumentRequest(r *http.Request) (requests.DocumentRequest, error) {
	var requestBody requests.DocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return requests.DocumentRequest{}, err
	}
	return requestBody, nil
}

// handleDecodeError is a helper function that handles the decoding of a JSON request body.
// It checks if the error is a SyntaxError or an UnmarshalTypeError and returns a JSON response with a relevant error message.
// If the error is of neither type, it returns a generic error message.
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
