package domain

import (
	"errors"
	"time"
)

type Report struct {
	id string

	userID string
	groupID string

	time time.Time
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

func NewReport(id, userID, groupID, content string) (Report, error) {
	if id == "" {
		return Report{}, errors.New("id can not be empty")
	}
	if userID == "" {
		return Report{}, errors.New("user id can not be empty")
	}
	if groupID == "" {
		return Report{}, errors.New("group id can not be empty")
	}
	if content == "" {
		return Report{}, errors.New("content can not be empty")
	}
	return Report{
		id: id,
		userID: userID,
		groupID: groupID,
		content: content,
		time: time.Now(),
	}, nil
}

func UnmarshalFromDatabase(id, userID, groupID, content string, time time.Time) (Report, error)  {
	report, err := NewReport(id, userID, groupID, content)
	if err != nil {
		return Report{}, err
	}
	report.time = time
	return report, nil
}