package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func InitializeFirebaseApp() {
	opt := option.WithCredentialsFile("config/firebase-service-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}
	firebaseApp = app
}

func AuthClient() (*auth.Client, error) {
	if firebaseApp == nil {
		log.Fatalf("Firebase app has not been initialized. Call InitializeFirebaseApp first.")
	}
	return firebaseApp.Auth(context.Background())
}
