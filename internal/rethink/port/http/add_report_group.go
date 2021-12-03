package http

import (
	"encoding/json"
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/rethink/app/command"
)

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
		Name:   data.Name,
		UserID: user.ID,
	})
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, nil)
}
