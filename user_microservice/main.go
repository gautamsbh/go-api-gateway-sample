// 3 Microservices
// 1: API Gateway or Proxy Microservice
// API gateway takes the client request and check whether request endpoint is secured or not and routes the request accordingly
// 2: User Service
// User Service takes the username in header and gets the user profile from DB
// 3: Auth Service
// Auth service do authentication check from username in request header
// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// error messages
const (
	UnAuthorized = "you are not authorized for access"
	BadRequest   = "Username not found in header or empty"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string, host string) {
	// DB Initialization
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	// Gorilla Mux Initialization
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/profile", a.profile).Methods("GET")
	a.Router.HandleFunc("/service", a.service).Methods("GET")
}

func (a *App) profile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	u := user{Username: username}
	if err := u.getUserByUsername(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusUnauthorized, UnAuthorized)
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) service(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "user-microservice",
	})
}

// respond with Error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func (a *App) Run(addr string) {
	log.Printf("User microservice starting on port user_microservice:%v", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// DB struct and methods
// User struct
// DB query to fetch user profile
type user struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// get user by username
func (u *user) getUserByUsername(db *sql.DB) error {
	return db.QueryRow("SELECT id, first_name, last_name, username FROM users WHERE username=$1", u.Username).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Username)
}

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_HOST"),
	)
	a.Run(":" + (os.Getenv("SERVICE_PORT")))
}
