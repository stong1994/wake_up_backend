package auth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestJwtToken(t *testing.T) {
	data := NewTokenInfo("1234567890abc", "小明")
	token, err := GenToken(data)
	assert.NoError(t, err)
	fmt.Println("token", token)
	tokenSps := strings.SplitN(token, ".", 2)
	assert.Equal(t, tokenSps[0], "Bearer")

	token = tokenSps[1]
	decodeToken, err := DecodeToken(token)
	assert.NoError(t, err)
	assert.Equal(t, data.UserID, decodeToken.UserID)
	assert.Equal(t, data.DisplayName, decodeToken.DisplayName)
}
