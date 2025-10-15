package controllers

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"github.com/ordo_meritum/features/auth/models"
	"github.com/ordo_meritum/features/auth/services"
	"github.com/ordo_meritum/shared/middleware"
)

type Controller struct {
	service *services.AuthService
}

func NewController(service *services.AuthService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login-or-register", c.HandleLoginOrRegister).Methods("POST")
}

func (c *Controller) HandleLoginOrRegister(
	w http.ResponseWriter,
	r *http.Request,
) {
	verifiedToken, ok := r.Context().Value(middleware.VerifiedTokenKey).(*auth.Token)
	if !ok || verifiedToken == nil {
		middleware.JSON(w, http.StatusUnauthorized, map[string]string{"error": "No authenticated user found in context"})
		return
	}

	userToSync := &models.User{
		FirebaseUID: verifiedToken.UID,
	}

	user, err := c.service.LoginOrRegister(context.Background(), userToSync)
	if err != nil {
		middleware.JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error during login process"})
		return
	}

	middleware.JSON(w, http.StatusOK, user)

}
