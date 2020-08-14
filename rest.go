package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

// RestService defines the service
type RestService struct {
	user    string
	pass    string
	port    int
	useSsl bool
	door    *Door
	router  *mux.Router
}

// NewRestService is a constructor
func NewRestService(door *Door, user, pass string, port int, useSsl bool) (*RestService, error) {
	rest := &RestService{
		user:    user,
		pass:    pass,
		useSsl: useSsl,
		port:    port,
		door:    door,
		router:  mux.NewRouter().StrictSlash(true),
	}
	rest.router.HandleFunc("/door", rest.Status).Methods("GET")
	rest.router.HandleFunc("/door/close", rest.Close).Methods("POST")
	rest.router.HandleFunc("/door/open", rest.Open).Methods("POST")
	rest.router.HandleFunc("/door/hold", rest.Hold).Methods("POST")
	return rest, nil
}

// Listen for new connections
func (rest *RestService) Listen() {
	if rest.useSsl {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", rest.port), "server.crt", "server.key", rest.router))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", rest.port), rest.router))
	}
}

func (rest *RestService) authenticate(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header["Authorization"]
	if len(authHeader) > 0 {
		auth := strings.SplitN(authHeader[0], " ", 2)
		if len(auth) == 2 && auth[0] == "Basic" {
			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) == 2 && pair[0] == rest.user && pair[1] == rest.pass {
				return true
			}
		}
	}
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
	return false
}

func logger(r *http.Request) {
	log.Printf("%s    %s", r.Method, r.RequestURI)
}

// Status message is returned here
func (rest *RestService) Status(w http.ResponseWriter, r *http.Request) {
	if rest.authenticate(w, r) {
		msg, _ := rest.door.check()
		fmt.Fprintf(w, msg)
	}
	logger(r)
}

// Open the door
func (rest *RestService) Open(w http.ResponseWriter, r *http.Request) {
	if rest.authenticate(w, r) {
		for msg := range rest.door.open() {
			fmt.Fprintf(w, msg)
		}
	}
	logger(r)
}

// Close the door
func (rest *RestService) Close(w http.ResponseWriter, r *http.Request) {
	if rest.authenticate(w, r) {
		for msg := range rest.door.close() {
			fmt.Fprintf(w, msg)
		}
	}
	logger(r)
}

// Hold the door / prevent auto-close
func (rest *RestService) Hold(w http.ResponseWriter, r *http.Request) {
	if rest.authenticate(w, r) {
		msg, _ := rest.door.hold()
		fmt.Fprintf(w, msg)
	}
	logger(r)
}
