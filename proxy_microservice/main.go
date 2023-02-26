// Proxy Services behaves as API gateway
// Accepts /profile and /service endpoint
// call auth service for authentication, if authenticated then sends the request to user service to get profile
// for /service which is unsecured endpoint, will directly call user service and return response

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

var authServiceUrl string
var userServiceUrl string

const (
	UnProcessableEntity  = "error from downstream service"
	UserNameRequestError = "username not found in request header"
)

type appError struct {
	str string
}

func (e appError) Error() string {
	return fmt.Sprintf("%s", e.str)
}

func makeGetRequest(url string, headers map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("some error occurred: %#v", err)
		return nil, err
	}
	// set all headers in request
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: time.Second * 100}
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("some error occurred: %#v", err)
		return resp, err
	}
	// Print the HTTP Status Code and Status Name
	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
	log.Printf("response %#v", resp)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	} else {
		log.Printf("request handle error, http status code %#v", resp.StatusCode)
		return resp, &appError{str: UnProcessableEntity}
	}
}

func (a *App) Initialize(userServiceHost string, authServiceHost string) {
	// initialize microservice url
	userServiceUrl = userServiceHost
	authServiceUrl = authServiceHost
	// Gorilla Mux Initialization
	a.Router = mux.NewRouter()
	a.Router.Handle("/profile", a.authMiddleware(http.HandlerFunc(a.profile))).Methods("GET")
	a.Router.HandleFunc("/service", a.service).Methods("GET")
}

// middleware check the username in header
// if username not found, throws
func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("username")
		if username == "" {
			respondWithError(w, http.StatusUnprocessableEntity, UserNameRequestError)
			return
		}
		resp, err := makeGetRequest(authServiceUrl+"/auth", map[string]string{"username": r.Header.Get("username")})
		if err != nil {
			log.Printf("error in request %#v", err)
			respondWithError(w, resp.StatusCode, err.Error())
			return
		}
		defer resp.Body.Close()
		next.ServeHTTP(w, r)
	})
}

// user profile request
func (a *App) profile(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing profile API proxy service")
	resp, err := makeGetRequest(userServiceUrl+"/profile", map[string]string{"username": r.Header.Get("username")})
	if err != nil {
		log.Printf("error in request %#v", err)
		respondWithError(w, resp.StatusCode, err.Error())
		return
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	respondWithJSON(w, http.StatusOK, result)
}

// service request
func (a *App) service(w http.ResponseWriter, r *http.Request) {
	resp, err := makeGetRequest(userServiceUrl+"/service", map[string]string{})
	if err != nil {
		log.Printf("error in request %#v", err)
		respondWithError(w, resp.StatusCode, err.Error())
		return
	}
	defer resp.Body.Close()
	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	respondWithJSON(w, http.StatusOK, result)
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

// start server
func (a *App) Run(addr string) {
	log.Printf("Proxy microservice starting on socker proxy_microservice:%v", addr)
	log.Fatalln(http.ListenAndServe(addr, a.Router))
}

// entry point
func main() {
	a := App{}
	a.Initialize(
		os.Getenv("USER_MICROSERVICE_HTTP_SCHEME")+"://"+os.Getenv("USER_MICROSERVICE_HOST"),
		os.Getenv("AUTH_MICROSERVICE_HTTP_SCHEME")+"://"+os.Getenv("AUTH_MICROSERVICE_HOST"),
	)
	a.Run(":" + os.Getenv("SERVICE_PORT"))
}
