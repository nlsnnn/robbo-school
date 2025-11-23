package postgresql

import (
	"api/internal/config"
	"fmt"

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
