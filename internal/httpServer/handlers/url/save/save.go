package save

import (
	"log/slog"
	"net/http"

	resp "github.com/GlebusDev/urlShortener/internal/lib/api/response"
	"github.com/GlebusDev/urlShortener/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			var errText = "failed to decode request body"
			log.Error(errText, sl.Err(err))
 
			render.JSON(w, r, resp.Error(errText))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			var errText = "invalid request"
			log.Error(errText, sl.Err(err))

			render.JSON(w, r, resp.Error(errText))

			return
		}
	}
}
