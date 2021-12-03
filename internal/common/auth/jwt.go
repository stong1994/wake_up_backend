package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenInfo struct {
	UserID      string
	DisplayName string
	NBF         time.Time // valid begin time
	EXP         time.Time // valid after time
	IAT         time.Time // now
}

func NewTokenInfo(userId, displayName string) TokenInfo {
	return TokenInfo{
		UserID:      userId,
		DisplayName: displayName,
		IAT:         time.Now(),
		NBF:         time.Now(),
		EXP:         time.Now().Add(time.Hour * 24 * 7),
	}
}

var hmacSampleSecret = []byte("wake_up")
var InvalidTokenErr = errors.New("invalid token")

func GenToken(data TokenInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      data.UserID,
		"display_name": data.DisplayName,
		"exp":          data.EXP.Unix(),
		"nbf":          data.NBF.Unix(),
		"iat":          data.IAT.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenStr, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", "Bearer", tokenStr), nil
}

func DecodeToken(token string) (TokenInfo, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if err != nil {
		return TokenInfo{}, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return TokenInfo{
			UserID:      claims["user_id"].(string),
			DisplayName: claims["display_name"].(string),
		}, nil
	}

	return TokenInfo{}, InvalidTokenErr
}
