package main

import (
    "fmt"
    "log"
    "time"
    "github.com/stianeikeland/go-rpio"
)


const (
    CONTROL_PIN_NUM = 4  //physical pin 7
    SENSOR_PIN_NUM  = 17 //physical pin 11
)

// time we detected door open
var open_time time.Time = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

var control_pin rpio.Pin
var sensor_pin rpio.Pin

func setup() {
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

func is_open() bool {
    return sensor_pin.Read() == rpio.Low
}

func status() bool {
    if is_open() {
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
        if status() && time.Since(open_time).Hours() > 1 {
            log.Printf("Door has been open for > 1 hour.  Closing door...")
            toggle_door()
        }
        time.Sleep(5000 * time.Millisecond)
    }
}

func check_door() string {
    if status() {
        return fmt.Sprintf("Door has been Open for %f minutes", time.Since(open_time).Minutes())
    } else {
        return fmt.Sprintf("Door is Closed")
    }
}

func open_door() string {
    if is_open() {
        return fmt.Sprintf("Door is already Open")
    } else {
        toggle_door()
        return fmt.Sprintf("Opening...")
    }
}

func close_door() string {
    if !is_open() {
        return fmt.Sprintf("Door is already Closed")
    } else {
        toggle_door()
        return fmt.Sprintf("Closing...")
    }
}
