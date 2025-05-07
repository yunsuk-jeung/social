package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	ID         int64    `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Password   password `json:"-"`
	CreatedAt  string   `json:"created_at"`
	IsActivate bool     `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at 
		FROM users	
		WHERE id = $1
	`
	user := &User{}

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		}
	}

	return user, nil
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, password, email)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	err := tx.QueryRowContext(
		ctx, query,
		user.Username,
		user.Password.hash,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {

		//  create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		//  create the user invite
		err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user that this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. update the user to be activated
		user.IsActivate = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. clean the invitation
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
    SELECT u.id, u.username, u.email, u.created_at, u.is_active
    FROM users u
    JOIN user_invitations ui ON u.id = ui.user_id
    WHERE ui.token = $1 AND ui.expiry > $2
  `
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActivate,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `
    INSERT INTO user_invitations (token, user_id, expiry)
    VALUES ($1, $2, $3)
  `
	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
    UPDATE users SET username = $1, email = $2, is_active = $3
    WHERE id=$4
  `

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActivate, user.ID)
	if err != nil {
		return err
	}
	return nil

}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
    DELETE FROM user_invitations WHERE user_id = $1
  `

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil

}
