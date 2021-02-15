package store

type Store interface {
	User() UserRepo
}
