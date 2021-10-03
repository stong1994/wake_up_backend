package port

import (
	"github.com/kataras/iris/v12"
	"wake_up_backend/internal/rethink/app"
	"wake_up_backend/internal/rethink/app/query"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) FindAllReportWithGroup(ctx iris.Context)  {
	pageSize, _ := ctx.URLParamInt("page_size")
	pageNo, _ := ctx.URLParamInt("page_no")
	userID := ctx.URLParam("user_id") // TODO
	data, err := h.app.Queries.ReportAllTypeList.Handle(ctx.Request().Context(), query.ReportAllTypeList{
		PageNum:  pageNo,
		PageSize: pageSize,
		UserID:   userID,
	})
	// todo
	if err != nil {
		ctx.WriteString(err.Error())
		return
	}
	ctx.JSON(map[string]interface{}{"data": data})
}
