package port

import (
	"github.com/go-chi/chi"
	"net/http"
)

type ServerInterface interface {
	FindAllReportWithGroup(w http.ResponseWriter, r *http.Request)
	AddReport(w http.ResponseWriter, r *http.Request)
	AddReportGroup(w http.ResponseWriter, r *http.Request)
	FindReportGroups(w http.ResponseWriter, r *http.Request)
	FindAllReport(w http.ResponseWriter, r *http.Request)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
//func Handler(si ServerInterface) http.Handler {
//	return HandlerFromMux(si)
//}

// HandlerFromMux creates http.Handler with routing
func HandlerFromMux(si ServerInterface, router chi.Router) http.Handler {
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	router.Post("/report", si.AddReport)
	router.Get("/report/list", si.FindAllReportWithGroup)
	router.Post("/report/group", si.AddReportGroup)
	router.Get("/report/group/list", si.FindReportGroups)
	router.Get("/report/all", si.FindAllReport)
	return router
}
