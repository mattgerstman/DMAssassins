package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
)

var db *sql.DB

const (
	usersPath = "/users/"
)

func serveGetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Path[len(usersPath):]
		if err != nil {
			http.Error(w, "Couldn't parse ID from users path.", http.StatusBadRequest)
			return
		}

		var user_id string

		err = db.QueryRow("SELECT user_id FROM dm_users WHERE email = $1", email).Scan(&user_id)
		switch {
			case err == sql.ErrNoRows:
				http.NotFound(w, r)
				return
			case err != nil:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		w.Write([]byte(user_id))
	}
}

func serveUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			serveGetUser(db)(w, r)
		// case "POST":
		// 	servePostUser(db)(w, r)
		// case "DELETE":
		// 	serveDeleteUser(db)(w, r)
		default:
		}
	}
}


func connect() {
	db, err = sql.Open("postgres", "postgres://localhost?dbname=dmassassins&sslmode=disable")
}

func StartServer() {
	connect()
	defer db.Close()

	http.HandleFunc(usersPath, serveUsers(db))
	http.ListenAndServe(":8000", nil)
}
