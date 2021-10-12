package port

import (
	"encoding/json"
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/user/app"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) Login(w http.ResponseWriter, r *http.Request) {
	var data LoginInfo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}

	user, err := h.app.Queries.GetUser.Handle(r.Context(), data.Account, data.Password)
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	token, err := auth.GenToken(auth.NewTokenInfo(user.ID, user.DisplayName))
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, map[string]interface{}{"token": token})
}
