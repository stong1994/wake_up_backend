package http

import (
	"encoding/json"
	"net/http"
	"wake_up_backend/internal/common/auth"
	"wake_up_backend/internal/common/server"
	"wake_up_backend/internal/common/server/httperr"
	"wake_up_backend/internal/rethink/app/command"
)

func (h HttpServer) CompleteRethink(w http.ResponseWriter, r *http.Request) {
	var data CompleteRethink
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		httperr.BadRequest(err.Error(), err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.CompleteRethink.Handle(r.Context(), command.CompleteRethink{
		RethinkID:      data.ID,
		UserID:         user.ID,
		ReportContent:  data.ReportContent,
		RethinkContent: data.RethinkContent,
	})
	if err != nil {
		httperr.InternalError(err.Error(), err, w, r)
		return
	}
	server.RenderResponse(w, r, nil)
}
