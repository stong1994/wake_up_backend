package domain

import (
	"context"
	"errors"
	"github.com/stong1994/kit_golang/sstr"
	"time"
)

type Report struct {
	id string

	userID  string
	groupID string

	time    time.Time
	content string
}

func (r Report) ID() string {
	return r.id
}

func (r Report) UserID() string {
	return r.userID
}

func (r Report) GroupID() string {
	return r.groupID
}

func (r Report) Content() string {
	return r.content
}

func (r Report) Time() time.Time {
	return r.time
}

func (r *Report) UpdateContent(content string) error {
	if len(content) > 500 {
		return errors.New("the max length of content is 500")
	}
	r.content = content
	return nil
}

type checkReportGroupFunc func(ctx context.Context, userID, groupID string) (bool, error)

func NewReport(ctx context.Context, userID, groupID string,
	checkGroup checkReportGroupFunc) (Report, error) {
	if userID == "" {
		return Report{}, errors.New("user id can not be empty")
	}
	if groupID == "" {
		return Report{}, errors.New("group id can not be empty")
	}

	if checkGroup != nil {
		pass, err := checkGroup(ctx, userID, groupID)
		if err != nil {
			return Report{}, err
		}
		if !pass {
			return Report{}, errors.New("user does not own this group")
		}
	}

	return Report{
		id:      sstr.UUIDHex(),
		userID:  userID,
		groupID: groupID,
		time:    time.Now(),
	}, nil
}
