package port

import (
	"encoding/json"
	"github.com/stong1994/kit_golang/sweb"
	"net/http"
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
	err := h.app.Commands.AddReport.Handle(r.Context(), command.AddReport{
		ReportID: data.ReportID,
		GroupID:  data.GroupID,
		Content:  data.Content,
		UserID:   data.UserID,
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
	userID, _ := sweb.URLParamString(r, "user_id") // TODO context
	data, err := h.app.Queries.ReportAllTypeList.Handle(r.Context(), query.ReportAllTypeList{
		PageNum:  pageNo,
		PageSize: pageSize,
		UserID:   userID,
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
	err := h.app.Commands.AddReportGroup.Handle(r.Context(), command.AddReportGroup{
		ID:     data.GroupID,
		Name:   data.Name,
		UserID: data.UserID,
	})
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, nil)
}

func (h HttpServer) FindReportGroups(w http.ResponseWriter, r *http.Request) {
	if err := sweb.ParseForm(r); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}
	userID, _ := sweb.URLParamString(r, "user_id") // TODO context
	data, err := h.app.Queries.ReportGroupList.Handle(r.Context(), userID)
	// todo
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
