package storage

import "errors"

var (
	ErrCourseNotFound = errors.New("course not found")
)

type Course struct {
	Id          int
	Title       string
	Description string
	Created_at  string
}
