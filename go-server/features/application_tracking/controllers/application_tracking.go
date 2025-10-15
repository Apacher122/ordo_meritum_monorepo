package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"github.com/ordo_meritum/database/models"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	"github.com/ordo_meritum/features/application_tracking/services"
	"github.com/ordo_meritum/shared/middleware"
	errors "github.com/ordo_meritum/shared/types/errors"
	"github.com/ordo_meritum/shared/webrender"
)

var logger = log.With().
	Str("controller", "application_tracking").
	Logger()

type Controller struct {
	service *services.AppTrackerService
}

func NewController(service *services.AppTrackerService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(secureRouter *mux.Router, authRouter *mux.Router) {
	secureRouter.HandleFunc("/apps/track", c.TrackApplication).Methods("POST")
	authRouter.HandleFunc("/apps/track/list", c.ListApplications).Methods("GET")
	authRouter.HandleFunc("/track/{id:[0-9]+}", c.GetTrackedApplication).Methods("GET")
	authRouter.HandleFunc("/track/{id:[0-9]+}/status", c.UpdateStatus).Methods("PUT")
}

func getUserID(r *http.Request) (string, error) {
	verifiedToken, ok := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)
	if !ok || verifiedToken == nil {
		return "", errors.ErrNoUserID
	}
	return verifiedToken.UID, nil
}

func parseIDFromVars(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}

func (c *Controller) TrackApplication(w http.ResponseWriter, r *http.Request) {
	uid, _ := getUserID(r)
	apiKey := r.Context().Value(middleware.APIKeyContextKey)

	var requestBody request.JobPostingRequest
	if webrender.DecodeJSONBody(w, r, &requestBody) != nil {
		return
	}

	jobID, err := c.service.QueueApplicationTracking(
		context.Background(),
		apiKey.(string),
		uid,
		requestBody,
	)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	middleware.JSON(w, http.StatusCreated, jobID)
}

func (c *Controller) ListApplications(w http.ResponseWriter, r *http.Request) {
	uid, _ := getUserID(r)

	applications, err := c.service.ListTrackedApplications(context.Background(), uid)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}
	middleware.JSON(w, http.StatusOK, applications)
}

func (c *Controller) GetTrackedApplication(w http.ResponseWriter, r *http.Request) {
	roleID, err := parseIDFromVars(r)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid role ID format"})
		return
	}

	application, err := c.service.GetTrackedApplicationByID(context.Background(), roleID)
	if err != nil {
		middleware.JSON(w, http.StatusNotFound, nil)
		return
	}
	middleware.JSON(w, http.StatusOK, application)
}

func (c *Controller) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	uid, _ := getUserID(r)

	roleID, err := parseIDFromVars(r)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid role ID format"})
		return
	}

	var payload struct {
		Status models.AppStatus `json:"status"`
	}
	if webrender.DecodeJSONBody(w, r, &payload) != nil {
		return
	}

	err = c.service.UpdateApplicationStatus(context.Background(), uid, roleID, payload.Status)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Status updated successfully"})
}
