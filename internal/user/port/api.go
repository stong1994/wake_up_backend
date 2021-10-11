package port

import (
	"github.com/go-chi/chi"
	"net/http"
)

type ServerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
}

// HandlerFromMux creates http.Handler with routing
func HandlerFromMux(si ServerInterface, router chi.Router) http.Handler {
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	router.Post("/user/login", si.Login)
	return router
}
