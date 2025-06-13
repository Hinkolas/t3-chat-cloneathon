package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/auth/argon2id"
)

func (s *Service) Authenticate(ctx context.Context, username string, password string) (*User, error) {

	// Retrieve the user with the provided username from the database
	var user User
	err := s.db.QueryRowContext(ctx, `
		SELECT
			id, username, email, password_hash,
			created_at, updated_at,
			is_verified, mfa_active
		FROM users
		WHERE username = ?
		LIMIT 1;`,
		username,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt,
		&user.IsVerified, &user.MFAActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: username", ErrInvalidCredentials)
		} else {
			return nil, err
		}
	}

	// Authenticate user by verifiyng provided password against stored hash
	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("%w: password", ErrInvalidCredentials)
	}

	return &user, nil

}

func (s *Service) Authorize(ctx context.Context, sessionToken string) (*Session, error) {

	// Retrieve the session with that token from the database
	var session Session
	err := s.db.QueryRowContext(ctx, `
		SELECT
			id, user_id, token,
			issued_at, renewed_at, time_to_live,
			is_verified
		FROM sessions
		WHERE token = ?
		ORDER BY issued_at DESC
		LIMIT 1;`,
		sessionToken,
	).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.IssuedAt, &session.RenewedAt, &session.TimeToLive,
		&session.IsVerified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		} else {
			return nil, err
		}
	}

	// Check if the session is verified (only relevant for when MFA is enabled by the user)
	if !session.IsVerified {
		return nil, ErrSessionUnverified
	}

	// Check if the session is expired
	if time.Now().Unix() >= session.RenewedAt+session.TimeToLive {
		return nil, ErrSessionExpired
	}

	return &session, nil

}
