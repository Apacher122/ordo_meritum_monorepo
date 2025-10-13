package services

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/database/users"
	auth_models "github.com/ordo_meritum/features/auth/models"
)

type AuthService struct {
	userRepo users.Repository
}

func NewAuthService(userRepo users.Repository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) LoginOrRegister(ctx context.Context, authUser *auth_models.User) (*db_models.User, error) {
	l := log.With().
		Str("service", "auth").
		Logger()

	existingUser, err := s.userRepo.GetUserByFirebaseUID(ctx, authUser.FirebaseUID)
	if err != nil && !errors.Is(err, users.ErrUserNotFound) {
		l.Error().Err(err).Str("uid", authUser.FirebaseUID).Msg("Error fetching user by UID")
		return nil, err
	}

	if existingUser != nil {
		l.Warn().Str("uid", authUser.FirebaseUID).Msg("User already exists")
		return existingUser, nil
	}

	l.Info().Str("uid", authUser.FirebaseUID).Msg("Creating new user")
	newUser, err := s.userRepo.CreateUser(ctx, authUser.FirebaseUID)
	if err != nil {
		l.Error().Err(err).Str("uid", authUser.FirebaseUID).Msg("Error creating user")
		return nil, err
	}
	return newUser, nil
}
