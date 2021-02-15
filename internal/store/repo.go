package store

import (
	"github.com/spbu-timetable/api/internal/model"
)

type UserRepo interface {
	Create(*model.User) error
	Find(value string, field string) (*model.User, error)
	Update(value string, field string, id string) error
}
