package domain

import (
	"errors"
	"time"
)

type Rethink struct {
	id string

	userID     string
	reportTime time.Time
	groupID    string

	content        string
	recordTime     time.Time
	rethinkContent string
	rethinkTime    time.Time
}

func (r Rethink) ID() string {
	return r.id
}

func (r Rethink) UserID() string {
	return r.userID
}

func (r Rethink) GroupID() string {
	return r.groupID
}

func (r Rethink) Content() string {
	return r.content
}

func (r Rethink) RethinkContent() string {
	return r.rethinkContent
}

func (r Rethink) ReportTime() time.Time {
	return r.reportTime
}

func (r Rethink) RethinkTime() time.Time {
	return r.rethinkTime
}

func (r Rethink) RecordTime() time.Time {
	return r.recordTime
}

func (r *Rethink) UpdateContent(content string) error {
	if len(content) > 500 {
		return errors.New("the max length of content is 500")
	}
	r.content = content
	return nil
}

func (r *Rethink) UpdateRethinkContent(rethinkContent string) error {
	if len(rethinkContent) > 500 {
		return errors.New("the max length of rethink content is 500")
	}
	r.rethinkContent = rethinkContent
	return nil
}

func (r *Rethink) UpdateReportTime(reportTime time.Time) {
	r.reportTime = reportTime
}

func (r *Rethink) UpdateRethinkTime(rethinkTime time.Time) {
	r.rethinkTime = rethinkTime
}

func UnmarshalRethinkFromDB(id, userID, groupID, recordContent, rethinkContent string,
	reportTime, recordTime, rethinkTime time.Time) Rethink {

	if id == "" || userID == "" || groupID == "" {
		panic("rethink has been destroyed")
	}

	return Rethink{
		id:             id,
		userID:         userID,
		reportTime:     reportTime,
		groupID:        groupID,
		content:        recordContent,
		recordTime:     recordTime,
		rethinkContent: rethinkContent,
		rethinkTime:    rethinkTime,
	}
}
