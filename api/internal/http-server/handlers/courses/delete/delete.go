package delete

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	resp "api/internal/lib/api/response"
	"api/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
}

type CourseDeleter interface {
	DeleteCourse(ctx context.Context, id int) error
}

func New(log *slog.Logger, courseDeleter CourseDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Error("invalid course id", slog.String("id", idStr))
			render.JSON(w, r, resp.Error("invalid course id"))
			return
		}

		err = courseDeleter.DeleteCourse(r.Context(), id)
		if err != nil {
			log.Error("failed to delete course by id", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("course deleted", slog.Int("id", id))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
