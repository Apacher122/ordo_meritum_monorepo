package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	"github.com/ordo_meritum/features/application_tracking/services"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/middleware"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	lg "github.com/ordo_meritum/shared/utils/logger"
	"github.com/ordo_meritum/shared/webrender"
	"github.com/rs/zerolog/log"
)

var service = "app-tracker-controller"

type Controller struct {
	service *services.AppTrackerService
}

func NewController(service *services.AppTrackerService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(secureRouter *mux.Router, authRouter *mux.Router) {
	secureRouter.HandleFunc("/apps/track", c.trackNewApplicationHandler).Methods("POST")
	authRouter.HandleFunc("/apps/track/list", c.getApplicationListHandler).Methods("GET")
	authRouter.HandleFunc("/apps/update", c.updateApplicationHandler).Methods("PATCH")
	authRouter.HandleFunc("/apps/delete", c.deleteApplicationHandler).Methods("DELETE")
	authRouter.HandleFunc("/track/{id:[0-9]+}", c.getApplicationByIDHandler).Methods("GET")
}

func parseIDFromVars(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}

// trackNewApplicationHandler is an HTTP handler that creates a new job posting for a given request body.
// It requires the user context to be present in the request context.
// If the request body is malformed, it will return an HTTP 400 status with an error message.
// If the user context is not found, it will return an HTTP 500 status with an error message.
// If the job posting is successfully created, it will return an HTTP 202 status with a message indicating that the job posting was received and is being processed.
// The actual creation of the job posting is done in the background using a goroutine.
func (c *Controller) trackNewApplicationHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userCtx, ok := contexts.FromContext(r.Context())
	if !ok {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	var requestBody request.JobPostingRequest
	if webrender.DecodeJSONBody(w, r, &requestBody) != nil {
		return
	}
	middleware.JSON(w, http.StatusAccepted, map[string]string{"status": "Job posting received and is being processed"})
	go func() {
		bgCtx := context.Background()
		bgUserCtx := &contexts.UserContext{
			UID:    userCtx.UID,
			ApiKey: userCtx.ApiKey,
		}
		ctxWithUser := context.WithValue(bgCtx, contexts.UserContextKey, bgUserCtx)

		jobID, err := c.service.QueueApplicationTracking(ctxWithUser, &requestBody)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create job posting in background")
			return
		}
		log.Info().Int("jobID", jobID).Msg("Successfully created job posting in background")
	}()
}

// getApplicationListHandler retrieves a list of all tracked job postings for a given user.
// It requires the user context to be present in the request context.
// If the user context is not found, it will return an HTTP 500 status with an error message.
// If the job postings cannot be retrieved, it will return an HTTP 500 status with an error message.
// It will log a message with the job postings and the status of the request.
func (c *Controller) getApplicationListHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	_, ok := contexts.FromContext(r.Context())
	if !ok {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	applications, err := c.service.ListTrackedApplications(r.Context())
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}
	middleware.JSON(w, http.StatusOK, applications)
}

// getApplicationByIDHandler retrieves a tracked job posting by ID.
// It requires the user context to be present in the request context.
// If the user context is not found, it will return an HTTP 500 status with an error message.
// If the job posting cannot be found, it will return an HTTP 404 status with an error message.
// It will log a message with the job posting and the status of the request.
func (c *Controller) getApplicationByIDHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	roleID, err := parseIDFromVars(r)
	_, ok := contexts.FromContext(r.Context())
	if !ok || err != nil {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	application, err := c.service.GetTrackedApplicationByID(r.Context(), roleID)
	if err != nil {
		middleware.JSON(w, http.StatusNotFound, nil)
		return
	}
	middleware.JSON(w, http.StatusOK, application)
}

// updateApplicationHandler updates the application status for a given job posting.
// It takes in a request which contains the role ID and the application status.
// It will return an error if any errors occur.
// If the role ID cannot be found, it will return an error with ErrorCode set to ERR_DB_FAILED_TO_GET.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the update process fails, it will return an HTTP 500 status with an error message.
// If the update process is successful, it will return an HTTP 200 status with a message indicating that the status was updated successfully.
func (c *Controller) updateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	_, ok := contexts.FromContext(r.Context())
	if !ok {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	requestBody := request.ApplicationUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Error().
			Err(err).
			Str("service", "documents-controller").
			Msg("Failed to decode request body")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := c.service.UpdateApplicationStatus(r.Context(), &requestBody)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Status updated successfully"})
}

// deleteApplicationHandler is an HTTP handler that deletes a tracked job posting by ID.
// It requires the user context to be present in the request context.
// If the user context is not found, it will return an HTTP 500 status with an error message.
// If the job posting cannot be found, it will return an HTTP 404 status with an error message.
// It will log a message with the job posting and the status of the request.
// If the deletion process fails, it will return an HTTP 500 status with an error message.
// If the deletion process is successful, it will return an HTTP 200 status with a message indicating that the application was deleted successfully.
func (c *Controller) deleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	_, ok := contexts.FromContext(r.Context())
	if !ok {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	requestBody := request.ApplicationUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Error().
			Err(err).
			Str("service", "documents-controller").
			Msg("Failed to decode request body")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := c.service.DeleteApplicationByID(r.Context(), &requestBody)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete application"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Application deleted successfully"})
}
