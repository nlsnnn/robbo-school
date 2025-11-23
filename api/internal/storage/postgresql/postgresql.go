package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"api/internal/config"
	"api/internal/storage"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const driverName = "pgx"

type Storage struct {
	db *sqlx.DB
}

func New(dbConfig config.Db) (*Storage, error) {
	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%v/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	db, err := sqlx.Connect(driverName, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveCourse(ctx context.Context, course storage.Course) (storage.Course, error) {
	query := "INSERT INTO courses (title, description) VALUES (:title, :description) RETURNING id, created_at"

	rows, err := s.db.NamedQueryContext(ctx, query, course)
	if err != nil {
		return storage.Course{}, fmt.Errorf("sql insert error: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&course.Id, &course.CreatedAt); err != nil {
			return storage.Course{}, fmt.Errorf("scan error: %w", err)
		}
		return course, nil
	}

	return storage.Course{}, fmt.Errorf("no rows returned")
}

func (s *Storage) GetCourses(ctx context.Context) ([]storage.Course, error) {
	courses := []storage.Course{}
	err := s.db.SelectContext(ctx, &courses, "SELECT * FROM courses")
	if err != nil {
		return []storage.Course{}, fmt.Errorf("sql select error: %w", err)
	}

	return courses, nil
}

func (s *Storage) GetCourseByID(ctx context.Context, id int) (storage.Course, error) {
	course := storage.Course{}
	err := s.db.GetContext(ctx, &course, "SELECT * FROM courses WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.Course{}, storage.ErrCourseNotFound
		}
		return storage.Course{}, fmt.Errorf("sql get error: %w", err)
	}

	return course, nil
}

func (s *Storage) UpdateCourse(ctx context.Context, id int, course storage.Course) (storage.Course, error) {
	query := "UPDATE courses SET title = :title, description = :description WHERE id = :id RETURNING id, created_at"

	course.Id = id
	rows, err := s.db.NamedQueryContext(ctx, query, course)
	if err != nil {
		return storage.Course{}, fmt.Errorf("sql update error: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&course.Id, &course.CreatedAt); err != nil {
			return storage.Course{}, fmt.Errorf("scan error: %w", err)
		}
		return course, nil
	}

	return storage.Course{}, storage.ErrCourseNotFound
}

func (s *Storage) DeleteCourse(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM courses where id = $1", id)
	if err != nil {
		return fmt.Errorf("sql delete error: %w", err)
	}

	return nil
}
