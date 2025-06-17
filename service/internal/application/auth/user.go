package auth

import (
	"context"
	"time"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/auth/argon2id"
	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `json:"id"`            // Unique identifier for the user as UUIDv4.
	Email        string    `json:"email"`         // Primary email address associated with the user.
	Username     string    `json:"username"`      // Name chosen by the user and used for login.
	PasswordHash string    `json:"password_hash"` // Hashed password for the user as argon2id string.

	CreatedAt int64 `json:"created_at"` // The time a user's account was created as Unix Timestamp
	UpdatedAt int64 `json:"updated_at"` // The time a user's infos where last changed as Unix Timestamp

	IsVerified bool `json:"is_verified"` // Indicates whether an account is verified using the configured method
	MFAActive  bool `json:"mfa_active"`  // Indicates whether an account uses MFA (Multi-Factor-Auth) for login
}

func (s *Service) CreateUser(ctx context.Context, email string, username string, password string) (*User, error) {

	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	// Get current unix timestamp for issued_at and renewed_at fields
	now := time.Now().UnixMilli()

	var user = User{
		ID:           userID,
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsVerified:   false,
		MFAActive:    false,
	}

	// Insert new session into database
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO users (
			id, email, username, password_hash,
			created_at, updated_at,
			is_verified, mfa_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
		user.ID, user.Email, user.Username, user.PasswordHash,
		user.CreatedAt, user.UpdatedAt,
		user.IsVerified, user.MFAActive,
	)
	// TODO: fix bug where device info cant be converted when nil/empty
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {

	// Insert new session into database
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM users WHERE id = ?;`,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSessionNotFound
	}

	return nil

}
