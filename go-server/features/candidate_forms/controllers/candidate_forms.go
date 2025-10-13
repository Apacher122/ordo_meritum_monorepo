package controllers

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"github.com/ordo_meritum/features/candidate_forms/services"
	"github.com/ordo_meritum/shared/middleware"
	"github.com/ordo_meritum/shared/models/requests"
	"github.com/ordo_meritum/shared/utils/validators"
)

type Controller struct {
	service *services.CandidateFormsService
}

func NewController(service *services.CandidateFormsService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/upload-questions", c.PostQuestionnare).Methods("POST")
	// router.HandleFunc("/create-profile", c.CreateProfile).Methods("POST")

	// router.HandleFunc("/writings", c.UpsertWritingSamples).Methods("POST")
	// router.HandleFunc("/writings", c.GetWritingSamples).Methods("GET")
}

// func getUserID(r *http.Request) (string, error) {
// 	verifiedToken, ok := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)
// 	if !ok || verifiedToken == nil {
// 		return "", errors.New("no authenticated user found in context")
// 	}
// 	return verifiedToken.UID, nil
// }

func (c *Controller) PostQuestionnare(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	apiKey := r.Context().Value(middleware.APIKeyContextKey)
	verifiedToken, _ := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)

	requestBody := requests.RequestBody{}
	err := validators.DecodeJSON(w, r, &requestBody)
	if err != nil {
		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = c.service.SaveCandidateQuestionnaire(r.Context(), verifiedToken.UID, apiKey.(string), &requestBody)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

}

// func (c *Controller) UpsertQuestionnaire(w http.ResponseWriter, r *http.Request) {
// 	firebaseUID, err := getUserID(r)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
// 		return
// 	}

// 	var payload dto.CandidateQuestionsPayload
// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
// 		return
// 	}
// 	defer r.Body.Close()

// 	questionnaire, err := c.service.CreateOrUpdateQuestionnaire(context.Background(), firebaseUID, &payload)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save questionnaire"})
// 		return
// 	}

// 	middleware.JSON(w, http.StatusCreated, questionnaire)
// }

// func (c *Controller) GetQuestionnaire(w http.ResponseWriter, r *http.Request) {
// 	firebaseUID, err := getUserID(r)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
// 		return
// 	}

// 	questionnaire, err := c.service.GetQuestionnaireByUID(context.Background(), firebaseUID)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve questionnaire"})
// 		return
// 	}
// 	if questionnaire == nil {
// 		middleware.JSON(w, http.StatusNotFound, map[string]string{"message": "No questionnaire found for this user."})
// 		return
// 	}

// 	middleware.JSON(w, http.StatusOK, questionnaire)
// }

// func (c *Controller) UpsertWritingSamples(w http.ResponseWriter, r *http.Request) {
// 	firebaseUID, err := getUserID(r)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
// 		return
// 	}

// 	var payload dto.CandidateWritingPayload
// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		middleware.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
// 		return
// 	}
// 	defer r.Body.Close()

// 	err = c.service.CreateOrUpdateWritingSamples(context.Background(), firebaseUID, &payload)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save writing samples"})
// 		return
// 	}

// 	middleware.JSON(w, http.StatusCreated, map[string]string{"message": "Writing samples saved successfully"})
// }

// func (c *Controller) GetWritingSamples(w http.ResponseWriter, r *http.Request) {
// 	firebaseUID, err := getUserID(r)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
// 		return
// 	}

// 	samples, err := c.service.GetWritingSamplesByUID(context.Background(), firebaseUID)
// 	if err != nil {
// 		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve writing samples"})
// 		return
// 	}

// 	middleware.JSON(w, http.StatusOK, samples)
// }
