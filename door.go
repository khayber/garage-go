package main

import (
    "fmt"
    "log"
    "time"
    "github.com/stianeikeland/go-rpio"
)

// time we detected door open
var open_time time.Time = time.Time{}

var control_pin rpio.Pin
var sensor_pin rpio.Pin

func setup(control_pin_num, sensor_pin_num int) {
    rpio.Open()
    control_pin = rpio.Pin(control_pin_num)
    control_pin.Output()

    sensor_pin = rpio.Pin(sensor_pin_num)
    sensor_pin.Input()
    sensor_pin.PullUp()
}

func cleanup() {
    rpio.Close()
}

func toggle_door() {
    log.Printf("toggle_door")
    control_pin.Low()
    time.Sleep(1000 * time.Millisecond)
    control_pin.High()
}

func is_open() bool {
    return sensor_pin.Read() == rpio.Low
}

func status() bool {
    if is_open() {
        if DEBUG {
            fmt.Printf("Door has been Open for %v\n", time.Now().Round(time.Second).Sub(open_time))
        }
        if open_time.IsZero() {
            open_time = time.Now().Round(time.Second)
        }
        return true
    } else {
        open_time = time.Time{}
        return false
    }
}

func monitor(autoclose bool, closetime float64) {
    for true {
        if status() && (time.Since(open_time).Minutes() > closetime) {
            log.Printf("Door has been open too long, auto-closing.")
            if autoclose {
                toggle_door()
            }
        }
        time.Sleep(20000 * time.Millisecond)
    }
}

func check_door() string {
    if status() {
        return fmt.Sprintf("Door has been Open for %v", time.Now().Round(time.Second).Sub(open_time))
    } else {
        return "Door is Closed"
    }
}

func open_door() string {
    if is_open() {
        return "Door is already Open"
    } else {
        toggle_door()
        return "Opening..."
    }
}

func close_door() string {
    if !is_open() {
        return "Door is already Closed"
    } else {
        toggle_door()
        return "Closing..."
    }
}