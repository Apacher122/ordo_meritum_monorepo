package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"github.com/ordo_meritum/database/models"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	"github.com/ordo_meritum/features/application_tracking/services"
	"github.com/ordo_meritum/shared/middleware"
)

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
		return "", errors.New("no authenticated user found in context")
	}
	return verifiedToken.UID, nil
}

func (c *Controller) TrackApplication(w http.ResponseWriter, r *http.Request) {
	uid, err := getUserID(r)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to get user ID")
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to read request body")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot read request body"})
		return
	}
	defer r.Body.Close()

	apiKey := r.Context().Value(middleware.APIKeyContextKey)

	requestBody := request.JobPostingRequest{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to decode request body")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	jobID, err := c.service.QueueApplicationTracking(
		context.Background(),
		apiKey.(string),
		uid,
		requestBody,
	)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	middleware.JSON(w, http.StatusCreated, jobID)

}

func (c *Controller) ListApplications(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := getUserID(r)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to get user ID")
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	applications, err := c.service.ListTrackedApplications(context.Background(), firebaseUID)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to retrieve applications")
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve applications"})
		return
	}
	middleware.JSON(w, http.StatusOK, applications)
}

func (c *Controller) GetTrackedApplication(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := getUserID(r)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to get user ID")
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	roleIDStr := mux.Vars(r)["id"]
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Invalid role ID format, must be an integer")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid role ID format"})
		return
	}

	application, err := c.service.GetTrackedApplicationByID(context.Background(), firebaseUID, roleID)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to retrieve application")
		middleware.JSON(w, http.StatusNotFound, map[string]string{"error": "Application not found"})
		return
	}
	middleware.JSON(w, http.StatusOK, application)
}

func (c *Controller) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := getUserID(r)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to get user ID")
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	roleIDStr := mux.Vars(r)["id"]
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Invalid role ID format, must be an integer")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid role ID format"})
		return
	}

	var payload struct {
		Status models.AppStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Str("service", "application-tracking-controller").Msg("Failed to decode request body")
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err = c.service.UpdateApplicationStatus(context.Background(), firebaseUID, roleID, payload.Status)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Status updated successfully"})
}
