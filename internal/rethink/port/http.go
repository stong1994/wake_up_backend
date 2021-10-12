package port

import (
	"encoding/json"
	"github.com/stong1994/kit_golang/sweb"
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/rethink/app"
	"wake_up_backend/internal/rethink/app/command"
	"wake_up_backend/internal/rethink/app/query"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) AddReport(w http.ResponseWriter, r *http.Request) {
	var data EntityCreateReport
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.AddReport.Handle(r.Context(), command.AddReport{
		GroupID: data.GroupID,
		UserID:  user.ID,
	})
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, nil)
}

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
	// todo
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, data)
}

func (h HttpServer) AddReportGroup(w http.ResponseWriter, r *http.Request) {
	var data EntityCreateReportGroup
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.AddReportGroup.Handle(r.Context(), command.AddReportGroup{
		ID:     data.GroupID,
		Name:   data.Name,
		UserID: user.ID,
	})
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, nil)
}

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
