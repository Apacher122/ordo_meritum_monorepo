package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	db_models "github.com/ordo_meritum/database/models"
)

type UserUpdate struct {
	Name *string `db:"name"`
}

type Repository interface {
	CreateUser(ctx context.Context, firebaseUID string) (*db_models.User, error)
	GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*db_models.User, error)
	UpdateUser(ctx context.Context, firebaseUID string, updates *UserUpdate) (*db_models.User, error)
	DeleteUser(ctx context.Context, firebaseUID string) error
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

var ErrUserNotFound = errors.New("user not found")

func (r *postgresRepository) CreateUser(ctx context.Context, firebaseUID string) (*db_models.User, error) {
	var user db_models.User
	query := "INSERT INTO users (firebase_uid) VALUES ($1) RETURNING *"
	err := r.db.GetContext(ctx, &user, query, firebaseUID)
	if err != nil {
		return nil, err
	}
	log.Printf("DATABASE: Created new user with Firebase UID: %s", user.FirebaseUID)
	return &user, nil
}

func (r *postgresRepository) GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*db_models.User, error) {
	var user db_models.User
	query := "SELECT * FROM users WHERE firebase_uid = $1"
	err := r.db.GetContext(ctx, &user, query, firebaseUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	log.Printf("DATABASE: Fetched user with Firebase UID: %s", user.FirebaseUID)
	return &user, nil
}

func (r *postgresRepository) UpdateUser(ctx context.Context, firebaseUID string, updates *UserUpdate) (*db_models.User, error) {
	var setClauses []string
	var args []interface{}
	argId := 1

	if updates.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argId))
		args = append(args, *updates.Name)
		argId++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no update fields provided")
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf("UPDATE users SET %s WHERE firebase_uid = $%d RETURNING *",
		strings.Join(setClauses, ", "), argId)
	args = append(args, firebaseUID)

	var updatedUser db_models.User
	err := r.db.GetContext(ctx, &updatedUser, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	log.Printf("DATABASE: Updated user with Firebase UID: %s", updatedUser.FirebaseUID)
	return &updatedUser, nil
}

func (r *postgresRepository) DeleteUser(ctx context.Context, firebaseUID string) error {
	query := "DELETE FROM users WHERE firebase_uid = $1"
	result, err := r.db.ExecContext(ctx, query, firebaseUID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("DATABASE: No user found to delete with Firebase UID: %s", firebaseUID)
		return ErrUserNotFound
	}
	log.Printf("DATABASE: Deleted user with Firebase UID: %s", firebaseUID)
	return nil
}

var _ Repository = (*postgresRepository)(nil)
