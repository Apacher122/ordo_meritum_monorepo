package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ordo_meritum/database/models"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	"github.com/ordo_meritum/features/application_tracking/services"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/middleware"
	"github.com/ordo_meritum/shared/webrender"
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

func parseIDFromVars(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}

func (c *Controller) TrackApplication(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) ListApplications(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) GetTrackedApplication(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	roleID, err := parseIDFromVars(r)
	_, ok := contexts.FromContext(r.Context())
	if !ok || err != nil {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	var payload struct {
		Status models.AppStatus `json:"status"`
	}
	if webrender.DecodeJSONBody(w, r, &payload) != nil {
		return
	}

	err = c.service.UpdateApplicationStatus(r.Context(), roleID, payload.Status)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update status"})
		return
	}
	middleware.JSON(w, http.StatusOK, map[string]string{"message": "Status updated successfully"})
}
