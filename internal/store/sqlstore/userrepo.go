package sqlstore

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/spbu-timetable/api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo ...
type UserRepo struct {
	store *SQLStore
}

// Create inserts new user into the database and returns error if something went wrong
func (ur *UserRepo) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.GenerateToken("refresh"); err != nil {
		return err
	}

	if err := u.GenerateToken("access"); err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO users (id, firstname, lastname, email, password, creation_date, token) 
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')`,
		u.ID, u.Firstname, u.Lastname, u.Email, u.Password, u.CreationDate, u.RefreshToken)
	if _, err := ur.store.db.Exec(query); err != nil {
		return err
	}

	return nil
}

// Find ...
func (ur *UserRepo) Find(value string, field string) (*model.User, error) {

	if field != "id" && field != "email" {
		return nil, errors.New("undefined field")
	}

	u := model.User{}

	query := fmt.Sprintf("SELECT * FROM users WHERE `%s`='%s'", field, value)

	token := sql.NullString{}

	if err := ur.store.db.QueryRow(query).Scan(
		&u.ID,
		&u.Firstname,
		&u.Lastname,
		&u.Email,
		&u.Password,
		&u.CreationDate,
		&token); err != nil {
		return nil, err
	}
	u.RefreshToken = token.String

	return &u, nil
}

// Update ...
func (ur *UserRepo) Update(value string, field string, id string) error {

	if field == "password" {
		byteValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		value = string(byteValue)
	}

	query := fmt.Sprintf("UPDATE users SET `%s`='%s' WHERE id='%s'", field, value, id)
	if _, err := ur.store.db.Exec(query); err != nil {
		return err
	}

	return nil
}
