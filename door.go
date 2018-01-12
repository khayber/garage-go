package main

import (
    "fmt"
    "log"
    "time"
    "github.com/stianeikeland/go-rpio"
)

var FAKE = false

type State int
const (
    unknown State = iota
    closed
    closing
    open
    opening
    holding
)

type Door struct {
    state State
    // time we detected door open
    open_time time.Time

    control_pin rpio.Pin
    sensor_pin rpio.Pin
}

func NewDoor(control_pin_num, sensor_pin_num int) (*Door, error) {
    err := rpio.Open()
    door := &Door {
        open_time: time.Time{},
        state: unknown,
        control_pin: rpio.Pin(control_pin_num),
        sensor_pin: rpio.Pin(sensor_pin_num),
    }
    if err == nil {
        door.control_pin.Output()
        door.sensor_pin.Input()
        door.sensor_pin.PullUp()
    } else {
        FAKE = true
        log.Printf("FAKE!!!")
    }

    return door, nil
}

func (door *Door) cleanup() {
    if FAKE {return}
    rpio.Close()
}

func (door *Door) toggle() {
    log.Printf("toggle_door")
    if FAKE {return}
    door.control_pin.Low()
    time.Sleep(1000 * time.Millisecond)
    door.control_pin.High()
}

func (door *Door) is_open() bool {
    if FAKE {
        switch door.state {
        case open:
            return true
        case opening:
            door.state = open
            return false
        case closed:
            return false
        case closing:
            door.state = closed
            return false
        case holding:
            return true
        case unknown:
            return false
        }
    }
    return door.sensor_pin.Read() == rpio.Low
}

func (door *Door) status() bool {
    if door.is_open() {
        if door.open_time.IsZero() {
            door.open_time = time.Now().Round(time.Second)
        }
        return true
    } else {
        door.open_time = time.Time{} //reset to 0
        return false
    }
}

func (door *Door) monitor(closetime float64) {
    for true {
        if door.status() && (time.Since(door.open_time).Minutes() > closetime && door.state != holding) {
            log.Printf("Door has been open too long, auto-closing.")
            door.toggle()
        }
        time.Sleep(20000 * time.Millisecond)
    }
}

func (door *Door) check() (string, bool) {
    if door.status() {
        duration := time.Now().Round(time.Second).Sub(door.open_time)
        if door.state == holding {
            return fmt.Sprintf("Door has been Holding for %v\n", duration), true
        } else {
            return fmt.Sprintf("Door has been Open for %v\n", duration), true
        }
    } else {
        return "Door is Closed", false
    }
}

func (door *Door) open() chan string {
    c := make(chan string)
    go func(c chan string) {
        if door.is_open() {
            c <- "Open"
        } else {
            door.state = opening
            c <- "Opening"
            door.toggle()
            for !door.is_open() {
                time.Sleep(1000 * time.Millisecond)
            }
            door.state = open
            c <- "Open"
        }
        close(c)
    } (c)
    return c
}

func (door *Door) close() chan string {
    c := make(chan string)
    go func(c chan string) {
        if !door.is_open() {
            c <- "Closed"
        } else {
            door.state = closing
            c <- "Closing..."
            door.toggle()
            for door.is_open() {
                time.Sleep(1000 * time.Millisecond)
            }
            door.state = closed
            c <- "Closed"
        }
        close(c)
    } (c)
    return c
}

func (door *Door) hold() (string, bool) {
    if !door.is_open() {
        return "Closed", true
    } else {
        door.state = holding
        return "Holding", false
    }
}
