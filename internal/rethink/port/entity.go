package port

type EntityCreateReport struct {
	ReportID string `json:"report_id"`
	GroupID  string `json:"group_id"`
	Content  string `json:"content"`
	UserID   string `json:"user_id"`
}

type EntityCreateReportGroup struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
}
