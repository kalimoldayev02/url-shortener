package save

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/kalimoldayev02/url/internal/repository/storage"
	resp "github.com/kalimoldayev02/url/pkg/api/response"
	"github.com/kalimoldayev02/url/pkg/lib/logger/sl"
	"github.com/kalimoldayev02/url/pkg/random"
	"golang.org/x/exp/slog"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias, omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias, omitempty"`
}

type URLSaver interface {
	SaveUrl(urlToSave, alias string) (int64, error)
}

// TODO: move to config
const aliasLength = 4

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.url.save"
		var req Request

		log = log.With(
			slog.String("operation", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("reuest", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.Url, req.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.Url))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Alias:    alias,
		})
	}
}
