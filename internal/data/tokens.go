package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moutafatin/go-tasks-management-api/internal/validator"
)

const ScopeActivation = "activation"

type Token struct {
	PlainText string
	Hash      []byte
	UserID    int
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

type tokensModel struct {
	DB *pgxpool.Pool
}

func (t tokensModel) New(userID int, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.insert(token)
	return token, err
}

func (t tokensModel) insert(token *Token) error {
	stmt := `INSERT INTO tokens (user_id, hash, scope, expiry) VALUES ($1, $2, $3, $4)`

	args := []any{token.UserID, token.Hash, token.Scope, token.Expiry}
	_, err := t.DB.Exec(context.Background(), stmt, args...)
	return err
}

func (t tokensModel) DeleteAllForUser(userID int, scope string) error {
	stmt := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`

	_, err := t.DB.Exec(context.Background(), stmt, userID, scope)

	return err
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
