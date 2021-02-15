package server

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func (s *server) authenticate(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			s.error(w, r, http.StatusUnauthorized, nil)
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)

		tkn, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil || !tkn.Valid {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (s *server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Type", "text/css")
		w.Header().Add("Content-Type", "text/html")

		next.ServeHTTP(w, r)
	})
}
