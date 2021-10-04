package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

	fmt.Println("Starting HTTP server")

	http.ListenAndServe(addr, rootRouter)
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router)
	//addAuthMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

//func addAuthMiddleware(router *chi.Mux) {
//	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
//		router.Use(auth.HttpMockMiddleware)
//		return
//	}
//
//	var opts []option.ClientOption
//	if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
//		opts = append(opts, option.WithCredentialsFile(file))
//	}
//
//	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
//	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
//	if err != nil {
//		logrus.Fatalf("error initializing app: %v\n", err)
//	}
//
//	authClient, err := firebaseApp.Auth(context.Background())
//	if err != nil {
//		logrus.WithError(err).Fatal("Unable to create firebase Auth client")
//	}
//
//	router.Use(auth.FirebaseHttpMiddleware{authClient}.Middleware)
//}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
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

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(bts); err != nil {
		panic(err)
	}
}
