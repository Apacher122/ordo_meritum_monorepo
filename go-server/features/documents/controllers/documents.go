package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/features/documents/services"
	"github.com/ordo_meritum/shared/middleware"
	"github.com/rs/zerolog/log"
)

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
	generationFunc func(ctx context.Context, apiKey string, requestBody requests.DocumentRequest, uid string) (int, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		apiKey, ok := r.Context().Value(middleware.APIKeyContextKey).(string)
		if !ok {
			middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "API key not found"})
			return
		}
		verifiedToken, _ := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)
		if verifiedToken == nil {
			middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}

		requestBody, err := decodeDocumentRequest(r)
		if err != nil {
			handleDecodeError(w, err)
			return
		}

		jobID, err := generationFunc(r.Context(), apiKey, requestBody, verifiedToken.UID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to queue document for generation")
			middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to queue document for generation"})
			return
		}

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

func (c *Controller) GenerateResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	apiKey := r.Context().Value(middleware.APIKeyContextKey)
	verifiedToken, _ := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)

	requestBody := requests.DocumentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": msg})
			return

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Invalid type for field '%s'. Expected '%s' but received a '%s'.", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
			middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": msg})
			return
		default:
			middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
			return
		}
	}

	jobID, err := c.docService.QueueResumeGeneration(
		r.Context(),
		apiKey.(string),
		requestBody,
		verifiedToken.UID,
	)

	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to queue resume for generation"})
		return
	}

	middleware.JSON(w, http.StatusAccepted, map[string]any{
		"jobId":  jobID,
		"status": "processing_queued",
	})
}

func (c *Controller) GenerateCoverLetter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey := r.Context().Value(middleware.APIKeyContextKey)
	verifiedToken, _ := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)

	requestBody := requests.DocumentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Error().
			Err(err).
			Str("service", "documents-controller").
			Msg("Failed to decode request body")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	jobID, err := c.docService.QueueCoverLetterGeneration(
		r.Context(),
		apiKey.(string),
		requestBody,
		verifiedToken.UID,
	)

	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to queue resume for generation"})
		return
	}

	middleware.JSON(w, http.StatusAccepted, map[string]any{
		"jobId":  jobID,
		"status": "processing_queued",
	})
}
