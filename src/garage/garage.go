package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "strings"
    "net/http"
    "encoding/base64"
    "github.com/gorilla/mux"
    "github.com/stianeikeland/go-rpio"
)


const (
    CONTROL_PIN_NUM = 4  //physical pin 7
    SENSOR_PIN_NUM  = 17 //physical pin 11
)


var USER = "god"
var PASS = "damn"


// time we detected door open
var open_time time.Time = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

var control_pin rpio.Pin
var sensor_pin rpio.Pin

func setup() {

    u := os.Getenv("USER")
    if len(u) > 0 {USER = u}
    p := os.Getenv("PASS")
    if len(p) > 0 {PASS = p}

    rpio.Open()
    control_pin = rpio.Pin(CONTROL_PIN_NUM)
    control_pin.Output()

    sensor_pin = rpio.Pin(SENSOR_PIN_NUM)
    sensor_pin.Input()
    sensor_pin.PullUp()
}


func cleanup() {
    rpio.Close()
}


func toggle_door() {
    control_pin.Low()
    time.Sleep(500 * time.Millisecond)
    control_pin.High()
}


func check_door() bool {
    return sensor_pin.Read() == rpio.Low
}


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


func check() bool {
    if check_door() {
        if open_time.IsZero() {
            open_time = time.Now()
        }
        return true
    } else {
        open_time = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
        return false
    }
}


func monitor() {
    for true {
        if check() && time.Since(open_time).Hours() > 1 {
            log.Printf("Door has been open for > 1 hour.  Closing door...")
            toggle_door()
        }
        time.Sleep(5000 * time.Millisecond)
    }
}


func Door(w http.ResponseWriter, r *http.Request) {
    if authenticate(w, r) {
        if check() {
            fmt.Fprintf(w, "Door has been Open for %f minutes", time.Since(open_time).Minutes())
        } else {
            fmt.Fprintf(w, "Door is Closed")
        }
    }
    logger(r)
}


func Open(w http.ResponseWriter, r *http.Request) {
    if authenticate(w, r) {
        if check_door() {
            fmt.Fprintf(w, "Door is already Open")
        } else {
            fmt.Fprintf(w, "Opening...")
            toggle_door()
        }
    }
    logger(r)
}


func Close(w http.ResponseWriter, r *http.Request) {
    if authenticate(w, r) {
        if !check_door() {
            fmt.Fprintf(w, "Door is already Closed")
        } else {
            fmt.Fprintf(w, "Closing...")
            toggle_door()
        }
    }
    logger(r)
}


func main() {
    setup()

    go monitor()

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/door", Door).Methods("GET")
    router.HandleFunc("/door/close", Close).Methods("POST")
    router.HandleFunc("/door/open", Open).Methods("POST")
    //log.Fatal(http.ListenAndServe(":8080", router))
    log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", router))

    cleanup()
}

