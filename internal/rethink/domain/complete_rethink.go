package domain

import (
	"errors"
	"time"
)

type CompleteRethink struct {
	id     string
	userID string

	time           time.Time
	reportContent  string
	rethinkContent string
}

func (r CompleteRethink) ID() string {
	return r.id
}

func (r CompleteRethink) UserID() string {
	return r.userID
}

func (r CompleteRethink) ReportContent() string {
	return r.reportContent
}

func (r CompleteRethink) RethinkContent() string {
	return r.rethinkContent
}

func (r CompleteRethink) Time() time.Time {
	return r.time
}

func (r *CompleteRethink) UpdateRethinkContent(content string) error {
	if len(content) > 500 {
		return errors.New("the max length of content is 500")
	}
	r.rethinkContent = content
	return nil
}

func NewCompleteRethink(rethinkID, userID string, reportContent, rethinkContent string) (CompleteRethink, error) {
	if userID == "" || rethinkID == "" {
		return CompleteRethink{}, errors.New("user id and rethink id can not be empty")
	}
	if reportContent == "" && rethinkContent == "" {
		return CompleteRethink{}, errors.New("record content and rethink content all empty")
	}
	return CompleteRethink{
		id:             rethinkID,
		userID:         userID,
		reportContent:  reportContent,
		rethinkContent: rethinkContent,
	}, nil
}

func (r CompleteRethink) needComplete() bool {
	return r.reportContent != "" || r.rethinkContent != ""
}

func (r CompleteRethink) Complete(oldRethink Rethink) (Rethink, error) {
	if r.userID != oldRethink.userID {
		return Rethink{}, errors.New("user does not own this rethink")
	}
	if !r.needComplete() {
		return Rethink{}, nil
	}

	if r.reportContent != "" {
		if err := oldRethink.UpdateContent(r.reportContent); err != nil {
			return Rethink{}, err
		}
		oldRethink.UpdateReportTime(time.Now())
	}
	if r.rethinkContent != "" {
		if err := oldRethink.UpdateRethinkContent(r.rethinkContent); err != nil {
			return Rethink{}, err
		}
		oldRethink.UpdateRethinkTime(time.Now())
	}
	return oldRethink, nil
}
