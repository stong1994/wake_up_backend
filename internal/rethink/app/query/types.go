package query

import "time"

type RespReportAllGroupList struct {
	Total int64
	List []RespReportAllGroupItem
}

type RespReportAllGroupItem struct {
	GroupID string
	GroupName string
	List []RespReportSingleTypeItem
}

type RespReportSingleTypeItem struct {
	Content string
	ReportTime time.Time
	RethinkShortContent string
	RethinkContentID string
}
