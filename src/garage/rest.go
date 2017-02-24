package main

import (
    "fmt"
    "log"
    "strings"
    "net/http"
    "encoding/base64"
    "github.com/gorilla/mux"
)


var USER = ""
var PASS = ""

func authenticate(w http.ResponseWriter, r *http.Request) bool {
    auth_header := r.Header["Authorization"]
    if len(auth_header) > 0 {
        auth := strings.SplitN(auth_header[0], " ", 2)
        if len(auth) == 2 && auth[0] == "Basic" {
            payload, _ := base64.StdEncoding.DecodeString(auth[1])
            pair := strings.SplitN(string(payload), ":", 2)
            if len(pair) == 2 && pair[0] == USER && pair[1] == PASS {
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

func Door(w http.ResponseWriter, r *http.Request) {
    if authenticate(w, r) {
        msg := check_door()
        fmt.Fprintf(w, msg)
    }
    logger(r)
}

func Open(w http.ResponseWriter, r *http.Request) {
    if authenticate(w, r) {
        msg := open_door()
        fmt.Fprintf(w, msg)
    }
    logger(r)
}

func Close(w http.ResponseWriter, r *http.Request) {
    if authenticate(w, r) {
        msg := close_door()
        fmt.Fprintf(w, msg)
    }
    logger(r)
}

func rest(user, pass string, port int, use_ssl bool) {
    USER = user
    PASS = pass
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/door", Door).Methods("GET")
    router.HandleFunc("/door/close", Close).Methods("POST")
    router.HandleFunc("/door/open", Open).Methods("POST")
    if use_ssl {
        log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", port), "server.crt", "server.key", router))
    } else {
        log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
    }
}

