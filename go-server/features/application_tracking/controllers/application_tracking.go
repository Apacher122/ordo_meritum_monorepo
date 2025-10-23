package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	"github.com/ordo_meritum/features/application_tracking/services"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/middleware"
	"github.com/ordo_meritum/shared/webrender"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	service *services.AppTrackerService
}

func NewController(service *services.AppTrackerService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(secureRouter *mux.Router, authRouter *mux.Router) {
	secureRouter.HandleFunc("/apps/track", c.HandleTrackApplication).Methods("POST")
	authRouter.HandleFunc("/apps/track/list", c.HandleListApplications).Methods("GET")
	authRouter.HandleFunc("/apps/update", c.HandleUpdateApplication).Methods("PATCH")
	authRouter.HandleFunc("/track/{id:[0-9]+}", c.HandleGetTrackedApplication).Methods("GET")
}

func parseIDFromVars(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}

// HandleTrackApplication queues a job for tracking based on the request body and returns the job ID of the queued job.
// It takes in a request body and returns a job ID of the queued job and an error if any errors occur.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the request body is malformed, it will return an error with ErrorCode set to ERR_INVALID_REQUEST_BODY.
// If the job cannot be queued, it will return an error with ErrorCode set to ERR_FAILED_TO_QUEUE_JOB.
// It will log a message with the jobID and the status of the job.
func (c *Controller) HandleTrackApplication(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	_, ok := contexts.FromContext(r.Context())
	if !ok {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}
	var requestBody request.JobPostingRequest
	if webrender.DecodeJSONBody(w, r, &requestBody) != nil {
		return
	}

	jobID, err := c.service.QueueApplicationTracking(r.Context(), requestBody)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	middleware.JSON(w, http.StatusCreated, jobID)
}

// HandleListapplications returns a list of all tracked applications for the user.
// It requires the user context to be present in the request context.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the applications cannot be retrieved, it will return an error with ErrorCode set to ERR_FAILED_TO_RETRIEVE_APPLICATIONS.
// It will log a message with the applications and the status of the request.
func (c *Controller) HandleListApplications(w http.ResponseWriter, r *http.Request) {
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

// HandleGetTrackedApplication returns a single tracked application by ID.
// It requires the user context to be present in the request context and the roleID to be present in the request variables.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the roleID is not present or cannot be parsed, it will return an error with ErrorCode set to ERR_INVALID_REQUEST_BODY.
// If the application cannot be retrieved, it will return an error with ErrorCode set to ERR_FAILED_TO_RETRIEVE_APPLICATION.
// It will log a message with the application and the status of the request.
func (c *Controller) HandleGetTrackedApplication(w http.ResponseWriter, r *http.Request) {
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

// HandleUpdateApplication updates the status of an application based on the request body.
// It requires the user context to be present in the request context.
// If the request body is malformed, it will return an error with ErrorCode set to ERR_INVALID_REQUEST_BODY.
// If the application status cannot be updated, it will return an error with ErrorCode set to ERR_FAILED_TO_UPDATE_STATUS.
// It will log a message with the application and the status of the request.
func (c *Controller) HandleUpdateApplication(w http.ResponseWriter, r *http.Request) {
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
