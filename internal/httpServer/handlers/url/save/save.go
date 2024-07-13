package save

import (
	"log/slog"
	"net/http"
	"errors"

	resp "github.com/GlebusDev/urlShortener/internal/lib/api/response"
	"github.com/GlebusDev/urlShortener/internal/lib/logger/sl"
	"github.com/GlebusDev/urlShortener/internal/lib/random"
	"github.com/GlebusDev/urlShortener/internal/storage"
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

// TODO: move length to config
const aliasLength = 6

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
			validatorErrors := err.(validator.ValidationErrors)	
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationErrors(validatorErrors))

			return
		}

		alias := req.Alias
		var isAutoGen = false
		if alias == "" {
			isAutoGen = true
			alias = random.NewRandomString(aliasLength)
		}
		
		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			if isAutoGen {
				for i := 0; i < 10 && errors.Is(err, storage.ErrURLExists); i++ {
					alias = random.NewRandomString(aliasLength)
					err = urlSaver.SaveURL(req.URL, alias)
				}
				if errors.Is(err, storage.ErrURLExists) {
					errText := "failed to generate random string"
					log.Error(errText, sl.Err(err))
					render.JSON(w, r, resp.Error(errText))
					return
				}
			} else {
			errText := "url already exists"
			log.Info(errText, slog.String("url", req.URL))

			render.JSON(w, r, resp.Error(errText))
			return
			}
		}

		if err != nil {
			errText := "failed to add url"
			log.Error(errText, sl.Err(err))

			render.JSON(w, r, resp.Error(errText))
			return
		}

		log.Info("url added", slog.String("alias", alias))

		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias: alias,
	})
}