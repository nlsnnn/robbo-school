package get

import (
	resp "api/internal/lib/api/response"
	"api/internal/storage"
	"context"
)

type CoursesResponse struct {
	resp.Response
	Courses []storage.Course `json:"courses"`
}

type CourseResponse struct {
	resp.Response
	Course storage.Course `json:"course"`
}

type CourseGetter interface {
	GetCourseByID(ctx context.Context, id int) (storage.Course, error)
	GetCourses(ctx context.Context) ([]storage.Course, error)
}
