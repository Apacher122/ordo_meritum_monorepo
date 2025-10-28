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

// generateDocumentHandler is a factory function that creates an http.HandlerFunc for handling
// asynchronous document generation requests. It accepts a specific generation function
// (e.g., for resumes or cover letters) and returns a handler that performs common tasks:
//
//  1. Checks for the user context.
//  2. Decodes the incoming JSON request body into a requests.DocumentRequest.
//  3. Immediately sends an HTTP 202 Accepted response to the client, indicating the request
//     has been received and is being processed.
//  4. Launches a background goroutine to execute the provided generationFunc.
//  5. Creates a new background context (context.Background) for the goroutine, ensuring
//     that the process continues even if the client disconnects.
//  6. Populates a new UserContext within the background context for the generationFunc.
//  7. Calls the provided generationFunc with the background context and request body.
//  8. Logs the final jobID upon successful queuing or logs an error if the background
//     generation fails.
//
// This pattern allows for handling long-running document generation tasks without blocking
// the client or being interrupted by client disconnections.
//
// The generationFunc parameter is expected to perform the specific logic for queuing
// a document (e.g., interacting with LLMs, saving initial data, sending to Kafka)
// and should return the resulting jobID and any error encountered.
func (c *Controller) generateDocumentHandler(
	generationFunc func(
		ctx context.Context,
		requestBody requests.DocumentRequest,
	) (int, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userCtx, ok := contexts.FromContext(r.Context())
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

		middleware.JSON(w, http.StatusAccepted, map[string]any{
			"jobId":  requestBody.Options.JobID,
			"status": "processing_queued",
		})

		go func() {
			bgCtx := context.Background()
			bgUserCtx := &contexts.UserContext{
				UID:    userCtx.UID,
				ApiKey: userCtx.ApiKey,
			}
			ctxWithUser := context.WithValue(bgCtx, contexts.UserContextKey, bgUserCtx)

			jobID, err := generationFunc(ctxWithUser, requestBody)
			if err != nil {
				lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: err}.ErrorLog()
				return
			}
			logger.Info().Int("jobID", jobID).Msg("Document successfully queued for generation in background")
		}()
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
