package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

const defaultTimeout = 3 * time.Second

// insert user into database
func (db *DB) InsertUser(username, hashedPassword string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var id int

	query := `
		INSERT INTO users (created_at, username, hashed_password)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := db.GetContext(ctx, &id, query, time.Now(), username, hashedPassword)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				// Handle the unique violation error
				logger.Error(pqErr.Detail)
				return 0, errors.New("username already exists")
			} else {
				// Handle other errors

			}
		} else {
			// Handle other non-PostgreSQL-related errors
		}
		return 0, err
	}

	return id, err
}

// get user by id
func (db *DB) GetUser(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE id = $1`

	err := db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}

	return &user, err
}

// get user by email
func (db *DB) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	user := &User{}

	query := `SELECT * FROM users WHERE username = $1`

	err := db.GetContext(ctx, user, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return user, err
}

// update user password
func (db *DB) UpdateUserHashedPassword(id int, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE users SET hashed_password = $1 WHERE id = $2`

	_, err := db.ExecContext(ctx, query, hashedPassword, id)
	return err
}

// user model
type User struct {
	ID             int          `db:"id"`
	Username       string       `db:"username"`
	HashedPassword string       `db:"hashed_password"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	CreatedAt      time.Time    `db:"created_at"`
}
