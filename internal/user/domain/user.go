package domain

import (
	"errors"
	"time"
	"wake_up_backend/internal/common/auth"
)

var (
	NotFoundUser = errors.New("not found user")
	NoNameUser   = errors.New("user not set name")
	DeletedUser  = errors.New("user has been deleted")
)

type User struct {
	ID          string
	DisplayName string
}

func (u User) GenToken() (string, error) {
	return auth.GenToken(auth.NewTokenInfo(u.ID, u.DisplayName))
}

func UnmarshalFromDB(userID, name, displayName string, deleteTime time.Time) (User, error) {
	if userID == "" {
		return User{}, NotFoundUser
	}
	if deleteTime.After(time.Now()) {
		return User{}, DeletedUser
	}
	if name == "" {
		return User{}, NoNameUser
	}
	if displayName == "" {
		displayName = name
	}
	return User{
		ID:          userID,
		DisplayName: displayName,
	}, nil
}
