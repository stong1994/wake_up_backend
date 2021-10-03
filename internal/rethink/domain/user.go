package domain

import (
	"errors"
	"fmt"
)

type User struct {
	id string
}

func (u User) ID() string {
	return u.id
}

func NewUser(id string) (User, error) {
	if id == "" {
		return User{}, errors.New("id can not be empty")
	}
	return User{id: id}, nil
}

type ForbiddenToSeeReportError struct {
	RequestingUserID string
	ReportOwnerID    string
}

func (f ForbiddenToSeeReportError) Error() string {
	return fmt.Sprintf(
		"user '%s' can't see user '%s' report",
		f.RequestingUserID, f.ReportOwnerID,
	)
}

func CanUserSeeReport(user User, report Report) error {
	if user.ID() == report.UserID() {
		return nil
	}

	return ForbiddenToSeeReportError{user.ID(), report.ID()}
}
