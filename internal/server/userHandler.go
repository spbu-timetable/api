package server

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/spbu-timetable/api/internal/model"
	"golang.org/x/crypto/bcrypt"
)
// register godoc
// @Summary user registration
func (s *server) register() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		req := &model.User{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		uuid, err := exec.Command("uuidgen").Output()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u := &model.User{
			ID:           strings.TrimSuffix(string(uuid), "\n"),
			Firstname:    req.Firstname,
			Lastname:     req.Lastname,
			Email:        req.Email,
			Password:     string(hashed),
			CreationDate: time.Now(),
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}

}

func (s *server) login() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(req.Email, "email")
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errorIncorrectEmailOrPassword)
			return
		}

		if err := u.GenerateToken("refresh"); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := u.GenerateToken("access"); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.store.User().Update(u.RefreshToken, "token", u.ID)

		res := &response{u.RefreshToken, u.AccessToken}

		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) updateAccessToken() http.HandlerFunc {
	type request struct {
		ID           string `json:"id"`
		RefreshToken string `json:"token"`
	}

	type response struct {
		AccessToken string
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(req.ID, "id")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := u.UpdateAccessToken(req.RefreshToken); err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		res := response{u.AccessToken}

		s.respond(w, r, http.StatusOK, res)
	}
}

func (s *server) updateUser() http.HandlerFunc {

	type request struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
		Value       string `json:"value"`
		Field       string `json:"field"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().Update(req.Value, req.Field, req.ID); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) getUser() http.HandlerFunc {
	type request struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
	}

	type response struct {
		ID        string `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(req.ID, "id")
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		res := response{u.ID, u.Firstname, u.Lastname}

		s.respond(w, r, http.StatusOK, res)
	}

}
