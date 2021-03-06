package main

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"os"
	"wake_up_backend/internal/common/logs"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/user/port"
	"wake_up_backend/internal/user/service"
)

func main() {
	// 模拟环境变量，用于禁止token鉴权
	_ = os.Setenv("DISABLE_AUTH", "true")

	ctx := context.Background()

	logs.Init()

	app, cleanup := service.NewApplication(ctx)
	defer cleanup()

	server.RunHTTPServer(8081, func(router chi.Router) http.Handler {
		return port.HandlerFromMux(port.NewHttpServer(app), router)
	})
}
