package http

import (
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/rethink/app/query"
)

func (h HttpServer) FindReportGroups(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	data, err := h.app.Queries.ReportGroupList.Handle(r.Context(), user.ID)
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, groupListModel2Resp(data))
}

func groupListModel2Resp(list query.RespReportGroupList) ReportGroupList {
	resp := make([]ReportGroupItem, len(list))
	for i, v := range list {
		resp[i] = ReportGroupItem{
			GroupID: v.GroupID,
			Name:    v.Name,
			Count:   v.Count,
		}
	}
	return resp
}
