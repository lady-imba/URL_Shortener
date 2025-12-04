package save

import (
	"URL_SHORTENER/internal/lib/api/response"
	"URL_SHORTENER/internal/lib/logger/sl"
	"URL_SHORTENER/internal/lib/random"
	"URL_SHORTENER/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

const aliasLength = 10

type Request struct {
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request 

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return 
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return 
		}

		alias := req.Alias
		if alias == ""{
			alias = random.NewRandomString(aliasLength)
		}

		err = urlSaver.SaveURL(req.URL, req.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exist", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exist"))
			return 
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add url"))
			return 
		}

		log.Info("url added")
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias: alias,
		})
	}
}