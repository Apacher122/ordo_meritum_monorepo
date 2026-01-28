package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	request "github.com/ordo_meritum/features/job_guide/models/requests"
	"github.com/ordo_meritum/features/job_guide/services"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/middleware"
	"github.com/ordo_meritum/shared/webrender"
)

type Controller struct {
	service *services.JobGuideService
}

func NewController(service *services.JobGuideService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(secureRouter *mux.Router, authRouter *mux.Router) {
	secureRouter.HandleFunc("/company-info", c.HandleGetCompanyInfo).Methods("POST")
	secureRouter.HandleFunc("/match-summary", c.getMatchSummaryHandler).Methods("POST")
	secureRouter.HandleFunc("/guiding-answers", c.HandleGetGuidingAnswers).Methods("POST")
}

func (c *Controller) HandleGetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var requestBody request.JobGuideRequest
	if webrender.DecodeJSONBody(w, r, &requestBody) != nil {
		return
	}
	// TODO
}

// getMatchSummaryHandler is an HTTP handler that takes in a JobGuideRequest and returns a match summary.
// It first checks if a user is authenticated and if not, it returns an HTTP 401 status with an error message.
// Then, it attempts to decode the request body into a JobGuideRequest and if that fails, it returns an HTTP 400 status with an error message.
// If the request body is successfully decoded, it attempts to retrieve the match summary from the service and if that fails, it returns an HTTP 500 status with an error message.
// If the match summary is successfully retrieved, it returns the match summary as JSON with an HTTP 200 status.
func (c *Controller) getMatchSummaryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	_, ok := contexts.FromContext(r.Context())
	if !ok {
		middleware.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	var requestBody request.JobGuideRequest
	if webrender.DecodeJSONBody(w, r, &requestBody) != nil {
		return
	}

	summary, err := c.service.GetMatchSummary(r.Context(), &requestBody)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	middleware.JSON(w, http.StatusOK, summary)
}

func (c *Controller) HandleGetGuidingAnswers(w http.ResponseWriter, r *http.Request) {
	// TODO
}
