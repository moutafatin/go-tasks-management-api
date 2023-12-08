package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

type usersModel struct {
	DB *pgxpool.Pool
}

func (u usersModel) Insert(user *User) error {
	stmt := `
INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	err := u.DB.QueryRow(context.Background(), stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation && strings.Contains(pgError.Message, "users_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (u usersModel) GetByEmail(email string) (*User, error) {
	stmt := `
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1`

	var user User
	err := u.DB.QueryRow(context.Background(), stmt, email).Scan(&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u usersModel) Update(user *User) error {
	stmt := `
UPDATE users
SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	err := u.DB.QueryRow(context.Background(), stmt, args...).Scan(&user.Version)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation && strings.Contains(pgError.Message, "users_email") {
				return ErrDuplicateEmail
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (u usersModel) GetForToken(token, scope string) (*User, error) {
	hash := sha256.Sum256([]byte(token))
	stmt := `
      SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
      FROM users
      INNER JOIN tokens
      ON users.id = tokens.user_id
      WHERE tokens.hash = $1
      AND tokens.scope = $2
      AND tokens.expiry > $3`

	var user User
	err := u.DB.QueryRow(context.Background(), stmt, hash[:], scope, time.Now()).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.Password.hash, &user.Activated, &user.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &user, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plainText != nil {
		ValidatePasswordPlaintext(v, *user.Password.plainText)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
