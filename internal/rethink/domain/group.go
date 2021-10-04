package domain

import (
	"errors"
	"time"
)

type ReportGroup struct {
	id         string
	userID     string
	name       string
	createTime time.Time
}

func (r ReportGroup) ID() string {
	return r.id
}

func (r ReportGroup) UserID() string {
	return r.userID
}

func (r ReportGroup) Name() string {
	return r.name
}

func (r ReportGroup) CreateTime() time.Time {
	return r.createTime
}

func NewReportGroup(id, userID, name string) (ReportGroup, error) {
	if id == "" {
		return ReportGroup{}, errors.New("id can not be empty")
	}
	if userID == "" {
		return ReportGroup{}, errors.New("user id can not be empty")
	}
	if name == "" {
		return ReportGroup{}, errors.New("name can not be empty")
	}

	return ReportGroup{
		id:         id,
		userID:     userID,
		name:       name,
		createTime: time.Now(),
	}, nil
}
