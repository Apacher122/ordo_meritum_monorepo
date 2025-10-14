package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ordo_meritum/features/job_guide/services"
)

type Controller struct {
	service *services.JobGuideService
}

func NewController(service *services.JobGuideService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/company-info", c.HandleGetCompanyInfo).Methods("POST")
	router.HandleFunc("/match-summary", c.HandleGetMatchSummary).Methods("POST")
	router.HandleFunc("/guiding-answers", c.HandleGetGuidingAnswers).Methods("POST")
}

func (c *Controller) HandleGetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// requestBody := requests.RequestBody{}

	// if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
	// 	middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	// 	return
	// }

	// info, err := c.service.GetCompanyInfo(r.Context(), requestBody.Payload.CompanyName, requestBody.Payload.JobID)
	// if err != nil {
	// 	middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// 	return
	// }
	// middleware.JSON(w, http.StatusOK, info)
}

func (c *Controller) HandleGetMatchSummary(w http.ResponseWriter, r *http.Request) {
	// verifiedToken, ok := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)
	// if !ok {
	// 	middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "No authenticated user found"})
	// 	return
	// }

	// requestBody := requests.RequestBody{}
	// if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
	// 	middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	// 	return
	// }

	// err := c.service.GetMatchSummary(r.Context(), verifiedToken.UID, requestBody.Payload.JobID)
	// if err != nil {
	// 	middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// 	return
	// }

	// middleware.JSON(w, http.StatusOK, map[string]string{"message": "Match summary generated and saved successfully."})
}

func (c *Controller) HandleGetGuidingAnswers(w http.ResponseWriter, r *http.Request) {
	// verifiedToken, ok := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)
	// if !ok {
	// 	middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "No authenticated user found"})
	// 	return
	// }

	// requestBody := requests.RequestBody{}
	// if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
	// 	middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	// 	return
	// }

	// answers, err := c.service.GetGuidingAnswers(r.Context(), verifiedToken.UID, requestBody.Payload.JobID)
	// if err != nil {
	// 	middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	// 	return
	// }

	// middleware.JSON(w, http.StatusOK, answers)
}
