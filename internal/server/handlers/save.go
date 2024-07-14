package save

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Device   string `json:"device"`
	IP string `json:"ip,omitempty"`
}

type LaunchSaver interface {
	SaveGameLaunch(device string, ip string) (int64, error)
}

func New(log *slog.Logger, launchSaver LaunchSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			//
		}

		id, err :=  launchSaver.SaveGameLaunch(req.Device, req.IP)
		log.Info("The launch is saved", slog.Int64("id", id))
	}
}