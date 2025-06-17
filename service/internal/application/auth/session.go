package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// Session represents a authenticated session.
type Session struct {
	ID     uuid.UUID `json:"id"`      // Unique identifier for the session.
	UserID uuid.UUID `json:"user_id"` // Unique identifier of the user associated with the session.
	Token  string    `json:"token"`   // Token used to authenticate the session.

	// Timestamps for session issuance, renewal, and time to live.
	IssuedAt   int64 `json:"issued_at"`    // The time a session was created as Unix Timestamp
	RenewedAt  int64 `json:"renewed_at"`   // The time a session was last used as Unix Timestamp
	TimeToLive int64 `json:"time_to_live"` // The period of time since last renewed a session stays active in seconds

	IsVerified bool `json:"is_verified"` // Indicates whether a session is verified using the configured MFA method (always true when account does not use MFA)
}

func (s *Service) CreateSession(ctx context.Context, user *User) (*Session, error) {

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// Create a new 256bit session key from crypto random bytes
	sessionToken := make([]byte, 32)
	if _, err := rand.Read(sessionToken); err != nil {
		return nil, err
	}

	// Get current unix timestamp for issued_at and renewed_at fields
	now := time.Now().UnixMilli()

	var session = Session{
		ID:         sessionID,
		UserID:     user.ID,
		Token:      hex.EncodeToString(sessionToken),
		IssuedAt:   now,
		RenewedAt:  now,
		TimeToLive: 3600 * 1000,
		IsVerified: !user.MFAActive,
	}

	// Insert new session into database
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO sessions (
			id, user_id, token,
			issued_at, renewed_at, time_to_live,
			is_verified
		) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		session.ID, session.UserID, session.Token,
		session.IssuedAt, session.RenewedAt, session.TimeToLive,
		session.IsVerified,
	)
	// TODO: fix bug where device info cant be converted when nil/empty
	if err != nil {
		return nil, err
	}

	return &session, nil

}

func (s *Service) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {

	// Insert new session into database
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM sessions WHERE id = ?;`,
		sessionID,
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

func (s *Service) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	var session Session

	err := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, token, issued_at, renewed_at, time_to_live, is_verified
		FROM sessions WHERE id = ?;`,
		sessionID,
	).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.IssuedAt, &session.RenewedAt, &session.TimeToLive,
		&session.IsVerified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}
