package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/stong1994/kit_golang/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"wake_up_backend/internal/common/auth"
)

func RunHTTPServer(port int, createHandler func(router chi.Router) http.Handler) {
	RunHTTPServerOnAddr(":"+strconv.Itoa(port), createHandler)
}

func RunHTTPServerOnAddr(addr string, createHandler func(router chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	// we are mounting all APIs under /api path
	rootRouter.Mount("/api", createHandler(apiRouter))

	fmt.Println("Starting HTTP server, listing:", addr)

	http.ListenAndServe(addr, rootRouter)
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(slog.NewStructuredLogger())
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router)
	addAuthMiddleware(router)

	router.Use(
		// Configuring your server to return the X-ReportContent-Type-Options HTTP response header set to nosniff will
		// instruct browsers that support MIME sniffing to use the server-provided ReportContent-Type and not interpret
		// the content as a different content type.
		middleware.SetHeader("X-ReportContent-Type-Options", "nosniff"),
		// deny: The page cannot be displayed in a frame, regardless of the site attempting to do so.
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func addAuthMiddleware(router *chi.Mux) {
	if disableAuth, _ := strconv.ParseBool(os.Getenv("DISABLE_AUTH")); disableAuth {
		//router.Use(auth.HttpMockMiddleware)
		return
	}
	router.Use(auth.HttpMiddleware{}.Middleware)
}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "ReportContent-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}

type respData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RenderResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	resp := respData{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
	bts, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("ReportContent-Type", "application/json")
	if _, err = w.Write(bts); err != nil {
		panic(err)
	}
}
