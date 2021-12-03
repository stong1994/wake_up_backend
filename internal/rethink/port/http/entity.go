package http

type EntityCreateReport struct {
	GroupID string `json:"group_id"`
}

type EntityCreateReportGroup struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type ReportGroupList []ReportGroupItem

type ReportGroupItem struct {
	GroupID string `json:"group_id"`
	Name    string `json:"name"`
	Count   int    `json:"count"`
}

type AllReport struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	ReportTime int64  `json:"report_time"`

	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type CompleteRethink struct {
	ID             string `json:"id"`
	ReportContent  string `json:"report_content"`
	RethinkContent string `json:"rethink_content"`
}

type RespReportAllGroupList struct {
	Total int64                    `json:"total"`
	List  []RespReportAllGroupItem `json:"list"`
}

type RespReportAllGroupItem struct {
	GroupID   string                     `json:"group_id"`
	GroupName string                     `json:"group_name"`
	List      []RespReportSingleTypeItem `json:"list"`
}

type RespReportSingleTypeItem struct {
	RethinkID      string `json:"rethink_id"`
	ReportContent  string `json:"report_content"`
	ReportTime     int64  `json:"report_time"`
	RethinkContent string `json:"rethink_content"`
}
