package query

import "time"

type RespReportAllGroupList struct {
	Total int64
	List  []RespReportAllGroupItem
}

type RespReportAllGroupItem struct {
	GroupID   string
	GroupName string
	List      []RespReportSingleTypeItem
}

type RespReportSingleTypeItem struct {
	Content             string
	ReportTime          time.Time
	RethinkShortContent string
}

type RespReportGroupList []RespReportGroupItem

type RespReportGroupItem struct {
	GroupID string
	Name    string
	Count   int
}

type AllReport struct {
	ID         string
	Content    string
	ReportTime time.Time
	GroupID    string
	GroupName  string
}
