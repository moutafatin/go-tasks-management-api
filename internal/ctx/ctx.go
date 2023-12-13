package ctx

import (
	"context"
	"net/http"

	"github.com/moutafatin/go-tasks-management-api/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user key in request context")
	}

	return user
}
