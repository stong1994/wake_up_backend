package main

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/rethink/port"
	"wake_up_backend/internal/rethink/service"
)

func main() {
	ctx := context.Background()

	app, cleanup := service.NewApplication(ctx)
	defer cleanup()

	server.RunHTTPServer(8080, func(router chi.Router) http.Handler {
		return port.HandlerFromMux(port.NewHttpServer(app), router)
	})
}
