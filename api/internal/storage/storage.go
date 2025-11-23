package storage

import (
	"errors"
	"time"
)

var (
	ErrCourseNotFound = errors.New("course not found")
)

type Course struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}
