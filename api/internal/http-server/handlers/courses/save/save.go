package save

import (
	"context"
	"log/slog"
	"net/http"

	resp "api/internal/lib/api/response"
	"api/internal/lib/logger/sl"
	"api/internal/storage"

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

type CourseSaver interface {
	SaveCourse(ctx context.Context, course storage.Course) (storage.Course, error)
}

func New(log *slog.Logger, courseSaver CourseSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
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

		course, err := courseSaver.SaveCourse(r.Context(), storage.Course{
			Title:       req.Title,
			Description: req.Description,
		})
		if err != nil {
			log.Error("failed to save course", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save course"))
			return
		}

		log.Info("course saved", slog.Int("course_id", course.Id))

		responseOK(w, r, course)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, course storage.Course) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Course:   course,
	})
}
