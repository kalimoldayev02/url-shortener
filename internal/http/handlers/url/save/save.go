package save

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	resp "github.com/kalimoldayev02/url/pkg/api/response"
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

type UrlServer interface {
	Save(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlServer UrlServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.url.save"
		var req Request

		log = log.With(
			slog.String("operation", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

	}
}
