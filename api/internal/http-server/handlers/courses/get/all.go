package get

import (
	"log/slog"
	"net/http"

	resp "api/internal/lib/api/response"

	"api/internal/lib/logger/sl"
	"api/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func GetAll(log *slog.Logger, getter CourseGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		courses, err := getter.GetCourses(r.Context())
		if err != nil {
			log.Error("failed to get courses", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("courses retrieved", slog.Int("count", len(courses)))

		returnCourses(w, r, courses)
	}
}

func returnCourses(w http.ResponseWriter, r *http.Request, courses []storage.Course) {
	render.JSON(w, r, CoursesResponse{
		Response: resp.OK(),
		Courses:  courses,
	})
}
