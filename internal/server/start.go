package server

import (
	"database/sql"
	"net/http"
	"os"

	// importing mysql driver for "database/sql" package
	_ "github.com/go-sql-driver/mysql"
	"github.com/spbu-timetable/api/internal/store/sqlstore"
)

// Start creates new db connection and starts the server on port declared in "config.env" file
func Start() error {
	db, err := newDB()
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)

	server := newServer(store)
	port := ":" + os.Getenv("PORT")

	return http.ListenAndServe(port, server)
}

func newDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("CLEARDB_DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
