package sqlstore

import (
	"database/sql"

	"github.com/spbu-timetable/api/internal/store"
)

// SQLStore ...
type SQLStore struct {
	db       *sql.DB
	userRepo *UserRepo
}

// New ...
func New(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

// User ...
func (s *SQLStore) User() store.UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepo{store: s}

	return s.userRepo
}
