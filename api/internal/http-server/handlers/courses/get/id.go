package get

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	resp "api/internal/lib/api/response"
	"api/internal/storage"

	"api/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func GetByID(log *slog.Logger, getter CourseGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Error("invalid course id", slog.String("id", idStr))
			render.JSON(w, r, resp.Error("invalid course id"))
			return
		}

		course, err := getter.GetCourseByID(r.Context(), id)
		if errors.Is(err, storage.ErrCourseNotFound) {
			log.Info("course not found", slog.Int("id", id))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("course not found"))
			return
		}
		if err != nil {
			log.Error("failed to get course by id", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("course retrieved", slog.Int("id", course.Id))

		returnCourse(w, r, course)

	}
}

func returnCourse(w http.ResponseWriter, r *http.Request, course storage.Course) {
	render.JSON(w, r, CourseResponse{
		Response: resp.OK(),
		Course:   course,
	})
}
