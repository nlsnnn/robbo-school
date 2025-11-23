package update

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	resp "api/internal/lib/api/response"
	"api/internal/lib/logger/sl"
	"api/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type Response struct {
	resp.Response
	storage.Course
}

type CourseUpdater interface {
	GetCourseByID(ctx context.Context, id int) (storage.Course, error)
	UpdateCourse(ctx context.Context, id int, course storage.Course) (storage.Course, error)
}

func New(log *slog.Logger, courseUpdater CourseUpdater) http.HandlerFunc {
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

		course, err := courseUpdater.GetCourseByID(r.Context(), id)
		if errors.Is(err, storage.ErrCourseNotFound) {
			log.Info("course not found", slog.Int("id", id))
			render.JSON(w, r, resp.Error("course not found"))
			return
		}
		if err != nil {
			log.Error("failed to get course by id", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		var req Request

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		log.Info("request body decoded")

		course, err = courseUpdater.UpdateCourse(r.Context(), id, storage.Course{
			Title:       req.Title,
			Description: req.Description,
		})
		if err != nil {
			log.Error("failed to update course", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update course"))
			return
		}

		log.Info("course updated", slog.Int("course_id", course.Id))

		responseOK(w, r, course)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, course storage.Course) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Course:   course,
	})
}
