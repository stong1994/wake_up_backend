package http

import (
	"github.com/stong1994/kit_golang/sweb"
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/rethink/app/query"
)

func (h HttpServer) FindAllReport(w http.ResponseWriter, r *http.Request) {
	if err := sweb.ParseForm(r); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	pageSize, _ := sweb.URLParamInt(r, "page_size")
	pageNo, _ := sweb.URLParamInt(r, "page_no")
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	data, err := h.app.Queries.AllReport.Handle(r.Context(), user.ID, pageNo, pageSize)
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, allReportModel2Resp(data))
}

func allReportModel2Resp(list []query.AllReport) []AllReport {
	resp := make([]AllReport, len(list))
	for i, v := range list {
		resp[i] = AllReport{
			ID:         v.ID,
			Content:    v.Content,
			ReportTime: v.ReportTime.Unix(),

			GroupID:   v.GroupID,
			GroupName: v.GroupName,
		}
	}
	return resp
}
