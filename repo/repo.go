package repo

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("index not found")

type PostgresDB struct {
	//db *sql.DB
	db *gorm.DB
}

func NewPostgresDB(db *gorm.DB) *PostgresDB {
	return &PostgresDB{db: db}
}
