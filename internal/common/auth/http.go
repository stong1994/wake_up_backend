package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"wake_up_backend/internal/common/server/httperr"
)

type HttpMiddleware struct{}

func (a HttpMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bearerToken := a.tokenFromHeader(r)
		if bearerToken == "" {
			httperr.Unauthorised("empty-bearer-token", nil, w, r)
			return
		}

		token, err := DecodeToken(bearerToken)
		if err != nil {
			httperr.Unauthorised("unable-to-verify-jwt", err, w, r)
			return
		}

		// it's always a good idea to use custom type as context value (in this case ctxKey)
		// because nobody from the outside of the package will be able to override/read this value
		ctx = context.WithValue(ctx, userContextKey, User{
			ID:          token.UserID,
			DisplayName: token.DisplayName,
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a HttpMiddleware) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}

type User struct {
	ID          string
	DisplayName string
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

var (
	// if we expect that the user of the function may be interested with concrete error,
	// it's a good idea to provide variable with this error
	NoUserInContextError = errors.New("no user in context")
)

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}

	return User{}, NoUserInContextError
}
