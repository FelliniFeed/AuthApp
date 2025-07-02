package repository

import (
	"fmt"
	"strings"

	"github.com/FelliniFeed/AuthApp.git/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (username, password) values ($1, $2) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {

	var (
		user models.User
		query strings.Builder
		args []interface{}

	)

	query.WriteString(fmt.Sprintf("SELECT * FROM %s", usersTable))

	if username != "" && password != "" {
		query.WriteString(` WHERE username = $1 AND password = $2`)
		args = append(args, username, password)
	}

	err := r.db.Get(&user, query.String(), args...)

	return user, err
}

func(r *AuthPostgres) CreateRefreshToken(hash string, userID uuid.UUID) error {
	query := fmt.Sprintf("INSERT INTO %s (refresh_token, user_id) values ($1, $2)", refreshTokensTable)
	_, err := r.db.Exec(query, hash, userID)
	return err
}