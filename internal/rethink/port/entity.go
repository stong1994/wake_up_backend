package port

type EntityCreateReport struct {
	GroupID string `json:"group_id"`
}

type EntityCreateReportGroup struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
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
