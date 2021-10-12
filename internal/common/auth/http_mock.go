package auth

import (
	"context"
	"net/http"
	"time"
	"wake_up_backend/internal/common/server/httperr"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// HttpMockMiddleware is used in the local environment (which doesn't depend on database)
func HttpMockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims jwt.MapClaims
		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (i interface{}, e error) {
				return []byte("mock_secret"), nil
			},
			request.WithClaims(&claims),
		)
		if err != nil {
			httperr.BadRequest("unable-to-get-jwt", err, w, r)
			return
		}

		if !token.Valid {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, TokenInfo{
			UserID:      claims["user_id"].(string),
			DisplayName: claims["display_name"].(string),
			NBF:         claims["nbf"].(time.Time),
			EXP:         claims["exp"].(time.Time),
			IAT:         claims["iat"].(time.Time),
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
