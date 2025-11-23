package storage

import (
	"errors"
	"time"
)

var (
	ErrCourseNotFound = errors.New("course not found")
)

type Course struct {
	Id          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
