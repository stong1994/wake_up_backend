package http

import (
	"github.com/stong1994/kit_golang/sweb"
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/rethink/app/query"
)

func (h HttpServer) FindAllReportWithGroup(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.app.Queries.ReportAllTypeList.Handle(r.Context(), query.ReportAllTypeList{
		PageNum:  pageNo,
		PageSize: pageSize,
		UserID:   user.ID,
	})
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, convert2ReportWithGroup(data))
}

func convert2ReportWithGroup(data query.RespReportAllGroupList) RespReportAllGroupList {
	var groupList []RespReportAllGroupItem
	for _, v := range data.List {
		var rethinkList []RespReportSingleTypeItem
		for _, vv := range v.List {
			rethinkList = append(rethinkList, RespReportSingleTypeItem{
				RethinkID:      vv.RethinkID,
				ReportContent:  vv.ReportContent,
				ReportTime:     vv.ReportTime.Unix(),
				RethinkContent: vv.RethinkShortContent,
			})
		}
		groupList = append(groupList, RespReportAllGroupItem{
			GroupID:   v.GroupID,
			GroupName: v.GroupName,
			List:      rethinkList,
		})
	}
	return RespReportAllGroupList{
		Total: data.Total,
		List:  groupList,
	}
}
