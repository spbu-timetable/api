package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spbu-timetable/api/internal/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	errorIncorrectEmailOrPassword = errors.New("Incorrect email or password")
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureCors() {
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{http.MethodOptions, http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	s.router.Use(handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins))
}

func (s *server) configureRouter() {

	s.configureCors()
	s.router.Use(s.setContentType)

	userRouter := s.router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", s.register()).Methods(http.MethodPost, http.MethodOptions)
	userRouter.HandleFunc("/login", s.login()).Methods(http.MethodPost, http.MethodOptions)
	userRouter.HandleFunc("/token", s.updateAccessToken()).Methods("POST")
	userRouter.HandleFunc("/update", s.authenticate(s.updateUser())).Methods("PUT")
	userRouter.HandleFunc("/get", s.authenticate(s.getUser())).Methods("POST")

}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
