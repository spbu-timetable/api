package model

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID           string    `json:"id" validate:"required"`
	Firstname    string    `json:"firstname" validate:"required"`
	Lastname     string    `json:"lastname" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,gte=6"`
	CreationDate time.Time `json:"creation_date" validate:"required"`
	RefreshToken string    `json:"token"`
	AccessToken  string    `json:"access_token"`
}

// Validate function validates User structure and returns error if some of the fields are not valid
func (u *User) Validate() error {
	if err := validator.New().Struct(u); err != nil {
		return err
	}
	return nil
}

// ComparePassword ...
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// GenerateToken ...
func (u *User) GenerateToken(tokenType string) error {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = u.ID
	if tokenType == "access" {
		claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	}

	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		return errors.New("cannot find jwt secret key")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	if tokenType == "refresh" {
		u.RefreshToken = tokenString
	} else if tokenType == "access" {
		u.AccessToken = tokenString
	}

	return nil
}

// UpdateAccessToken ...
func (u *User) UpdateAccessToken(refreshToken string) error {

	if u.RefreshToken != refreshToken {
		return errors.New("refresh token is not valid")
	}

	err := u.GenerateToken("access")
	if err != nil {
		return err
	}

	return nil
}
